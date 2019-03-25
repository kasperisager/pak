package build

import (
	"fmt"
	"hash"
	"hash/fnv"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/kasperisager/pak/pkg/asset"
	"github.com/kasperisager/pak/pkg/asset/css"
	"github.com/kasperisager/pak/pkg/cli"
)

func Build(cmd *cli.Command) {
	flag := cmd.Flag()

	var (
		out    = *flag.String("o", "dist", "The directory to write files to")
		root   = *flag.String("root", "", "The root directory of entry files")
		vendor = *flag.String("vendor", "vendor", "The vendor directory of external files")
	)

	cmd.Usage("[flags] [entry files]")

	cmd.HandleFunc(func(filenames []string) {
		var err error

		if root == "" {
			root, err = computeRoot(filenames)

			if err != nil {
				cmd.Fatal(err)
			}
		}

		urls := make([]*url.URL, len(filenames))

		for i, filename := range filenames {
			url, err := resource(filename, root)

			if err != nil {
				cmd.Fatal(err)
			}

			urls[i] = url
		}

		graph, err := read(urls, root)

		if err != nil {
			cmd.Fatal(err)
		}

		if err = write(graph, out, vendor); err != nil {
			cmd.Fatal(err)
		}
	})
}

func computeRoot(filenames []string) (string, error) {
	var common []string

	for i, filename := range filenames {
		parts := strings.Split(filepath.Dir(filename), string(os.PathSeparator))

		if parts[0] == ".." {
			return "", fmt.Errorf("%s: file outside working directory", filename)
		}

		if i == 0 {
			common = parts
		} else {
			if len(common) < len(parts) {
				parts = parts[:len(common)]
			} else {
				common = common[:len(parts)]
			}

			for i, part := range parts {
				if common[i] != part {
					common = common[:i]
					break
				}
			}
		}
	}

	return filepath.Join(common...), nil
}

func resource(filename string, base string) (*url.URL, error) {
	filename, err := filepath.Rel(base, filename)

	if err != nil {
		return nil, err
	}

	return &url.URL{Path: "/" + filepath.ToSlash(filename)}, nil
}

func filename(url *url.URL, base string) string {
	var filename string

	if url.IsAbs() {
		filename = path.Join(base, url.Host, url.Path)
	} else {
		filename = path.Join(base, url.Path)
	}

	return filepath.FromSlash(filename)
}

func read(urls []*url.URL, root string) (*asset.Graph, error) {
	graph := asset.NewGraph()

	entries := make([]asset.Asset, len(urls))

	for i, url := range urls {
		asset, err := resolve(url, root, graph)

		if err != nil {
			return nil, err
		}

		entries[i] = asset
	}

	if err := compress(graph, entries); err != nil {
		return nil, err
	}

	return graph, nil
}

func resolve(url *url.URL, root string, graph *asset.Graph) (asset.Asset, error) {
	if asset, ok := graph.Lookup(url); ok {
		return asset, nil
	}

	var asset asset.Asset

	switch filepath.Ext(url.Path) {
	case ".css":
		bytes, err := fetch(url, root)

		if err != nil {
			return nil, err
		}

		asset, err = css.Asset(url, bytes)

		if err != nil {
			return nil, err
		}

	default:
		return nil, fmt.Errorf("%s: unsupported file type", url)
	}

	graph.Add(asset)

	for _, url := range asset.References() {
		reference, err := resolve(url, root, graph)

		if err != nil {
			return nil, err
		}

		graph.Reference(asset, reference)
	}

	return asset, nil
}

func compress(graph *asset.Graph, entries []asset.Asset) error {
	partitions, err := partition(graph, entries)

	if err != nil {
		return err
	}

	visited := make(map[asset.Asset]bool)

	for _, entry := range entries {
		_, err := merge(graph, partitions, entry, visited)

		if err != nil {
			return err
		}
	}

	return nil
}

func merge(
	graph *asset.Graph,
	partitions map[asset.Asset]string,
	asset asset.Asset,
	visited map[asset.Asset]bool,
) (asset.Asset, error) {
	if visited[asset] {
		return asset, nil
	}

	visited[asset] = true

	references, _ := graph.Outgoing(asset)

	for _, reference := range references {
		if visited[reference] {
			continue
		}

		reference, err := merge(graph, partitions, reference, visited)

		if err != nil {
			return nil, err
		}

		if partitions[asset] == partitions[reference] {
			if asset.Merge(reference) {
				graph.Merge(asset, reference)
			}
		}
	}

	return asset, nil
}

func partition(graph *asset.Graph, entries []asset.Asset) (map[asset.Asset]string, error) {
	hashes := make(map[asset.Asset]hash.Hash)

	for _, entry := range entries {
		data, err := entry.URL().MarshalBinary()

		if err != nil {
			return nil, err
		}

		err = mark(graph, hashes, entry, data, make(map[asset.Asset]bool))

		if err != nil {
			return nil, err
		}
	}

	partitions := make(map[asset.Asset]string)

	for asset, hash := range hashes {
		partitions[asset] = string(hash.Sum(nil))
	}

	return partitions, nil
}

func mark(
	graph *asset.Graph,
	hashes map[asset.Asset]hash.Hash,
	asset asset.Asset,
	data []byte,
	visited map[asset.Asset]bool,
) error {
	if visited[asset] {
		return nil
	}

	visited[asset] = true

	hash, ok := hashes[asset]

	if !ok {
		hash = fnv.New64()
		hashes[asset] = hash
	}

	_, err := hash.Write(data)

	if err != nil {
		return err
	}

	references, _ := graph.Outgoing(asset)

	for _, reference := range references {
		err := mark(graph, hashes, reference, data, visited)

		if err != nil {
			return err
		}
	}

	return nil
}

func fetch(url *url.URL, root string) ([]byte, error) {
	var (
		bytes []byte
		err   error
	)

	if url.IsAbs() {
		response, err := http.Get(url.String())

		if err != nil {
			return nil, err
		}

		defer response.Body.Close()

		bytes, err = ioutil.ReadAll(response.Body)
	} else {
		bytes, err = ioutil.ReadFile(filename(url, root))
	}

	return bytes, err
}

func write(graph *asset.Graph, out, vendor string) error {
	for _, asset := range graph.Assets() {
		var target string

		url := asset.URL()

		if url.IsAbs() {
			target = filename(url, filepath.Join(out, vendor))
		} else {
			target = filename(url, out)
		}

		if err := os.MkdirAll(filepath.Dir(target), 0755); err != nil {
			return err
		}

		if err := ioutil.WriteFile(target, asset.Data(), 0644); err != nil {
			return err
		}
	}

	return nil
}
