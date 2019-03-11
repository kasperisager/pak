package build

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/kasperisager/pak/pkg/asset"
	"github.com/kasperisager/pak/pkg/asset/css"
	"github.com/kasperisager/pak/pkg/cli"
)

func Build(cmd *cli.Command) {
	flag := cmd.Flag()

	var (
		out  = flag.String("o", "dist", "The directory to write files to")
		root = flag.String("root", "", "The root directory of entry files")
	)

	cmd.Usage("[flags] [entry files]")

	cmd.HandleFunc(func(args []string) {
		base, err := os.Getwd()

		if err != nil {
			cmd.Fatal(err)
		}

		assets, err := resolveAssetGraph(args, base)

		if err != nil {
			cmd.Fatal(err)
		}

		err = writeAssetGraph(assets, *out, *root)

		if err != nil {
			cmd.Fatal(err)
		}
	})
}

func resolveAssetGraph(filenames []string, base string) (map[string]asset.Asset, error) {
	assets := make(map[string]asset.Asset, len(filenames))

	for len(filenames) > 0 {
		var filename string

		filename, filenames = filenames[0], filenames[1:]

		filename, err := filepath.Rel(base, filepath.Join(base, filename))

		if _, ok := assets[filename]; ok {
			continue
		}

		asset, err := resolveAsset(filename)

		if err != nil {
			return nil, err
		}

		assets[filename] = asset

		for _, reference := range asset.References() {
			filenames = append(filenames, filepath.FromSlash(reference.Path))
		}
	}

	return assets, nil
}

func resolveAsset(filename string) (asset asset.Asset, err error) {
	switch filepath.Ext(filename) {
	case ".css":
		r, err := ioutil.ReadFile(filename)

		if err != nil {
			return nil, err
		}

		asset, err = css.Asset(filename, r)

		if err != nil {
			return nil, err
		}

	default:
		return nil, fmt.Errorf("%s: unsupported file type", filename)
	}

	return asset, nil
}

func writeAssetGraph(assets map[string]asset.Asset, out, root string) error {
	if root == "" {
		filenames := make([]string, 0, len(assets))

		for filename, _ := range assets {
			filenames = append(filenames, filename)
		}

		var err error

		root, err = commonDir(filenames)

		if err != nil {
			return err
		}
	}

	for filename, asset := range assets {
		filename, err := filepath.Rel(root, filename)

		if err != nil {
			return err
		}

		output := filepath.Join(out, filename)

		if !strings.HasPrefix(output, out) {
			return fmt.Errorf("%s: file outside root directory", filename)
		}

		if err := os.MkdirAll(filepath.Dir(output), 0755); err != nil {
			return err
		}

		if err := ioutil.WriteFile(output, asset.Data(), 0644); err != nil {
			return err
		}
	}

	return nil
}

func commonDir(filenames []string) (string, error) {
	var common []string

	for i, filename := range filenames {
		parts := strings.Split(filepath.Dir(filename), string(os.PathSeparator))

		if i == 0 {
			common = parts
		} else {
			if len(common) < len(parts) {
				parts = parts[:len(common)]
			} else {
				common = common[:len(parts)]
			}

			for i, part := range parts {
				if part == ".." {
					return "", fmt.Errorf("%s: file outside working directory", filename)
				}

				if common[i] != part {
					common = common[:i]
					break
				}
			}
		}
	}

	return filepath.Join(common...), nil
}
