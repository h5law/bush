package walk

import (
	"errors"
	"fmt"
	"io/fs"
	"path/filepath"
	"strings"

	"github.com/fatih/color"
	"github.com/spf13/viper"
)

func Walk(root string) error {
	err := filepath.WalkDir(root, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			fmt.Printf("\"%s\" [error checking dir: %v]", path, err)
			return filepath.SkipDir
		}

		// Skip anything matching ignore pattern
		ignoreStr := viper.GetString("ignore")
		ignore := strings.Split(ignoreStr, ",")
		for _, v := range ignore {
			if v != "" && strings.Contains(path, v) {
				return nil
			}
		}

		// Check if dir
		isDir := d.IsDir()
		if isDir {
			color.Set(color.Bold, color.FgBlue)
			defer color.Unset()
		}

		// Get indent
		shortPath := strings.TrimPrefix(path, root)
		shortPath = strings.TrimPrefix(shortPath, "/")
		indent := strings.Count(shortPath, "/") + 1
		if path == root {
			indent = 0
		}

		// Get fs.FileInfo from d
		info, err := d.Info()
		if err != nil {
			fmt.Printf("\"%s\" [error getting path info]", path)
			return filepath.SkipDir
		}

		name := d.Name()
		if path == root {
			name = path
		}
		fmt.Printf("%s[%v] %s\n", strings.Repeat("\t", indent), info.Size(), name)
		return nil
	})

	if err != nil {
		return errors.New(fmt.Sprintf("\"%s\" [error: %s]", root, err))
	}

	return nil
}
