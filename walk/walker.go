package walk

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"sort"
	"strings"

	"github.com/fatih/color"
	"github.com/spf13/viper"
)

const Line = "\u2502"      //│
const LineRight = "\u251C" //├
const Dash = "\u2500"      //─
const Corner = "\u2514"    //└

var finalDir bool

type ByDirFirst []fs.DirEntry

func (b ByDirFirst) Len() int {
	return len(b)
}

func (b ByDirFirst) Swap(i, j int) {
	b[i], b[j] = b[j], b[i]
}

func (b ByDirFirst) Less(i, j int) bool {
	if b[i].IsDir() && b[j].IsDir() {
		// Both dirs go by name
		return b[i].Name() < b[j].Name()
	} else if b[i].IsDir() && !b[j].IsDir() {
		// Put dirs first
		return true
	} else if !b[i].IsDir() && b[j].IsDir() {
		return false
	} else {
		// Sort files by name
		return b[i].Name() < b[j].Name()
	}
}

func Walk(root string, dc, fc *int) error {
	// Open file at path
	f, err := os.Open(root)
	if err != nil {
		msg := fmt.Sprintf("\"%s\" [error opening dir]", root)
		return errors.New(msg)
	}

	// Increment for root
	*dc += 1

	// Print root path info
	color.Set(color.Bold, color.FgBlue)
	info, err := f.Stat()
	if err != nil {
		msg := fmt.Sprintf("\"%s\" [error checking path]", f.Name())
		return errors.New(msg)
	}
	fmt.Printf("[ %s] %s\n", ConvertBytes(info.Size()), f.Name())
	color.Unset()

	if err := loopDirs(root, root, 0, dc, fc); err != nil {
		fmt.Println(err)
	}

	return nil
}

func loopDirs(base, root string, depth int, dc, fc *int) error {
	// Get slice containing directory contents
	dirs, err := readDir(root)
	if err != nil {
		msg := fmt.Sprintf("\"%s\" [error getting directory contents]", root)
		return errors.New(msg)
	}

	lastDir := false
	for i, d := range dirs {
		if i == len(dirs)-1 {
			lastDir = true
		}
		root = strings.TrimSuffix(root, "/")
		path := root + "/" + d.Name()

		if ignore := checkIgnore(path); ignore {
			continue
		}

		// Get file info
		info, err := d.Info()
		if err != nil {
			msg := fmt.Sprintf("\"%s\" [error checking path]", d.Name())
			return errors.New(msg)
		}

		// Get levels
		levels := viper.GetUint("levels")
		if uint(depth) > levels && levels > 0 {
			continue
		}

		// TODO Figure this out
		/* Print tree chars
		│ Line
		├ LineRight
		─ Dash
		└ Corner
		*/
		if depth == 0 && i != len(dirs)-1 {
			fmt.Printf("%s%s ", LineRight, strings.Repeat(Dash, (depth+1)*4))
		} else if depth == 0 && i == len(dirs)-1 {
			finalDir = true
			fmt.Printf("%s%s ", Corner, strings.Repeat(Dash, (depth+1)*4))
		}
		if depth > 0 && !finalDir {
			fmt.Printf("%s", strings.Repeat(Line+"     ", depth))
		} else if depth > 0 && !lastDir && !finalDir {
			fmt.Printf("%s", strings.Repeat(Line+"     ", depth))
		} else if depth > 0 {
			fmt.Printf("      %s", strings.Repeat(Line+"     ", depth-1))
		}
		if depth > 0 && i != len(dirs)-1 {
			fmt.Printf("%s%s ",
				LineRight,
				strings.Repeat(Dash, 4),
			)
		} else if depth > 0 && i == len(dirs)-1 {
			fmt.Printf("%s%s ",
				Corner,
				strings.Repeat(Dash, 4),
			)
		}

		// Increment counters and set color
		if d.IsDir() {
			*dc += 1
			color.Set(color.Bold, color.FgBlue)
		} else {
			*fc += 1
		}

		// Print file info
		fmt.Printf("[ %s] %s\n", ConvertBytes(info.Size()), info.Name())
		color.Unset()

		// Loop through next dir
		if d.IsDir() {
			err := loopDirs(base, root+"/"+d.Name(), depth+1, dc, fc)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

// return true to ignore false to keep
func checkIgnore(path string) bool {
	// Skip anything matching ignore pattern
	ignoreStr := viper.GetString("ignore")
	ignore := strings.Split(ignoreStr, ",")
	for _, v := range ignore {
		if v != "" && strings.Contains(path, v) {
			return true
		}
	}

	return false
}

func readDir(path string) ([]fs.DirEntry, error) {
	// Open file at path
	f, err := os.Open(path)
	if err != nil {
		msg := fmt.Sprintf("\"%s\" [error opening dir]", path)
		return nil, errors.New(msg)
	}

	// Get directory contents
	dirs, err := f.ReadDir(-1)
	f.Close()
	if err != nil {
		msg := fmt.Sprintf("\"%s\" [error getting directory contents]", path)
		return nil, errors.New(msg)
	}

	// Sort directory contents
	dirsFirst := viper.GetBool("dirs-first")
	if dirsFirst {
		sort.Sort(ByDirFirst(dirs))
	} else {
		sort.Slice(dirs, func(i, j int) bool { return dirs[i].Name() < dirs[j].Name() })
	}
	if err != nil {
		return nil, err
	}

	return dirs, nil
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
