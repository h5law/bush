package walk

import (
	"fmt"
	"io/fs"
	"path/filepath"
	"strings"

	"github.com/fatih/color"
)

func Walker(path string, d fs.DirEntry, err error) error {
	if err != nil {
		fmt.Printf("\"%s\" [error checking dir: %v]", path, err)
		return filepath.SkipDir
	}

	// Check if dir
	isDir := d.IsDir()
	if isDir {
		color.Set(color.Bold, color.FgBlue)
		defer color.Unset()
	}

	// Get indent
	indent := strings.Count(path, "/")

	// Get fs.FileInfo from d
	info, err := d.Info()
	if err != nil {
		fmt.Printf("\"%s\" [error getting path info]", path)
		return filepath.SkipDir
	}

	fmt.Printf("%s[%v] %s\n", strings.Repeat("\t", indent), info.Size(), d.Name())
	return nil
}
