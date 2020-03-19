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
	"github.com/kasperisager/pak/pkg/cli"

	"github.com/kasperisager/pak/pkg/asset/blob"
	"github.com/kasperisager/pak/pkg/asset/css"
	"github.com/kasperisager/pak/pkg/asset/html"
	"github.com/kasperisager/pak/pkg/asset/importmap"
	"github.com/kasperisager/pak/pkg/asset/js"
	"github.com/kasperisager/pak/pkg/asset/webmanifest"
)

func Command(cmd *cli.Command) {
	flag := cmd.Flag()

	var (
		out  = flag.String("o", "dist", "The directory to write files to")
		root = flag.String("root", "", "The root directory of entry files")
		_    = flag.String("vendor", "vendor", "The vendor directory of external files")
	)

	cmd.Usage("[flags] [entry files]")

	cmd.HandleFunc(func(filenames []string) {
		root, out := *root, *out

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

		if err = write(graph, out); err != nil {
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
		resolved, err := resolve(url, root, graph, nil)

		if err != nil {
			return nil, err
		}

		graph.Add(resolved)

		err = collect(resolved, url, root, graph)

		if err != nil {
			return nil, err
		}

		entries[i] = resolved
	}

	if err := compress(graph, entries); err != nil {
		return nil, err
	}

	return graph, nil
}

func parse(url *url.URL, data []byte, mediaType string, flags asset.Flags) (asset.Asset, error) {
	switch mediaType {
	case css.MediaType:
		return css.From(url, data)

	case html.MediaType:
		return html.From(url, data)

	case js.MediaType:
		return js.From(url, data, flags)

	case importmap.MediaType:
		return importmap.From(url, data)

	case webmanifest.MediaType:
		return webmanifest.From(url, data, flags)

	default:
		return blob.From(url, data), nil
	}
}

func resolve(
	url *url.URL,
	root string,
	graph *asset.Graph,
	flags asset.Flags,
) (asset.Asset, error) {
	if asset, ok := graph.Lookup(asset.ByURL(url)); ok {
		return asset, nil
	}

	mediaType, data, err := fetch(url, root, flags)

	if err != nil {
		return nil, err
	}

	resolved, err := parse(url, data, mediaType, flags)

	if err != nil {
		return nil, err
	}

	return resolved, nil
}

func collect(
	asset asset.Asset,
	url *url.URL,
	root string,
	graph *asset.Graph,
) error {
	for _, reference := range asset.References() {
		url := asset.URL().ResolveReference(reference.URL())
		flags := reference.Flags()

		referenced, err := resolve(url, root, graph, flags)

		if err != nil {
			return err
		}

		graph.Add(referenced)
		graph.Relate(asset, referenced, reference)

		err = collect(referenced, url, root, graph)

		if err != nil {
			return err
		}
	}

	for _, embed := range asset.Embeds() {
		data := embed.Data()
		mediaType := embed.MediaType()
		flags := embed.Flags()

		embedded, err := parse(url, data, mediaType, flags)

		if err != nil {
			return err
		}

		graph.Add(embedded)
		graph.Relate(asset, embedded, embed)

		err = collect(embedded, url, root, graph)

		if err != nil {
			return err
		}
	}

	return nil
}

func compress(graph *asset.Graph, entries []asset.Asset) error {
	partitions, err := partition(graph, entries)

	if err != nil {
		return err
	}

	visited := make(map[asset.Asset]bool)

	for _, entry := range entries {
		err := merge(graph, partitions, entry, visited)

		if err != nil {
			return err
		}
	}

	return nil
}

func merge(
	graph *asset.Graph,
	partitions map[asset.Asset]string,
	target asset.Asset,
	visited map[asset.Asset]bool,
) error {
	if visited[target] {
		return nil
	}

	visited[target] = true

	edges, _ := graph.Outgoing(target)

	for _, related := range edges {
		err := merge(graph, partitions, related, visited)

		if err != nil {
			return err
		}

		if partitions[target] == partitions[related] {
			graph.Merge(target, related)
		}
	}

	return nil
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

	edges, _ := graph.Outgoing(asset)

	for _, related := range edges {
		err := mark(graph, hashes, related, data, visited)

		if err != nil {
			return err
		}
	}

	return nil
}

func fetch(url *url.URL, root string, flags asset.Flags) (mediaType string, data []byte, err error) {
	if url.IsAbs() {
		response, err := http.Get(url.String())

		if err != nil {
			return "", nil, err
		}

		defer response.Body.Close()

		mediaType = response.Header.Get("content-type")

		data, err = ioutil.ReadAll(response.Body)
	} else {
		data, err = ioutil.ReadFile(filename(url, root))

		if flags.Has("mediaType") {
			mediaType = flags.Get("mediaType").(string)
		} else {
			mediaType = asset.MediaTypeByURL(url)

			if mediaType == "" {
				mediaType = http.DetectContentType(data)
			}
		}
	}

	return mediaType, data, err
}

func write(graph *asset.Graph, out string) error {
	for _, asset := range graph.Assets() {
		var target string

		url := asset.URL()

		if url.IsAbs() {
			continue
		}

		target = filename(url, out)

		if err := os.MkdirAll(filepath.Dir(target), 0755); err != nil {
			return err
		}

		if err := ioutil.WriteFile(target, asset.Data(), 0644); err != nil {
			return err
		}
	}

	return nil
}
