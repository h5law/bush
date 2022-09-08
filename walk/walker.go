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

		// Get levels
		levels := viper.GetUint("levels")
		if uint(indent) > levels {
			return filepath.SkipDir
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

		// Convert filesize
		size := ConvertBytes(info.Size())

		name := d.Name()
		if path == root {
			name = path
		}
		fmt.Printf("[ %s] %s\n", size, name)
		return nil
	})

	if err != nil {
		return errors.New(fmt.Sprintf("\"%s\" [error: %s]", root, err))
	}

	return nil
}

func ConvertBytes(bytes int64) string {
	if bytes >= 1073741824 {
		f := fmt.Sprintf("%.1f", float64(bytes/1073741824))
		f = strings.Trim(f, ".0")
		return fmt.Sprintf("%sGB", f)
	} else if bytes >= 1048576 {
		f := fmt.Sprintf("%.1f", float64(bytes/1048576))
		f = strings.Trim(f, ".0")
		return fmt.Sprintf("%sMB", f)
	} else if bytes >= 1024 {
		f := fmt.Sprintf("%.1f", float64(bytes/1024))
		f = strings.Trim(f, ".0")
		return fmt.Sprintf("%sKB", f)
	} else {
		return fmt.Sprintf("%d", bytes)
	}
}
