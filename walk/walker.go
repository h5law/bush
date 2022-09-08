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

func Walk(root string, dc, fc *int) error {
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

		// Get indent
		shortPath := strings.TrimPrefix(path, root)
		shortPath = strings.TrimPrefix(shortPath, "/")
		indent := strings.Count(shortPath, "/") + 1
		if path == root {
			indent = 0
		}

		/* Print tree chars
		├ \u251C
		─ \u2500
		└ \u2514
		*/
		if indent == 1 {
			fmt.Printf("%s%s ", "\u251C", strings.Repeat("\u2500", 2))
		}
		if indent > 1 {
			fmt.Printf("%s%s ", "\u251C", strings.Repeat("\u2500", indent*4))
		}

		// Get fs.FileInfo from d
		info, err := d.Info()
		if err != nil {
			fmt.Printf("\"%s\" [error getting path info]", path)
			return filepath.SkipDir
		}

		// Check if dir
		isDir := d.IsDir()
		if isDir {
			color.Set(color.Bold, color.FgBlue)
			defer color.Unset()
		}

		// Increment counts
		if isDir {
			*dc += 1
		} else {
			*fc += 1
		}

		name := d.Name()
		if path == root {
			name = path
		}
		fmt.Printf("[ %v] %s\n", info.Size(), name)
		return nil
	})

	if err != nil {
		return errors.New(fmt.Sprintf("\"%s\" [error: %s]", root, err))
	}

	return nil
}
