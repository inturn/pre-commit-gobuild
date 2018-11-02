package helpers

import (
	"os"
	"path/filepath"
	"regexp"
)

// DirsWith finds all the directories in the root (including the root) containing files matching regex
func DirsWith(root, mask string) []string {
	var dirs []string
	filepath.Walk(root, func(path string, f os.FileInfo, _ error) error {
		if !f.IsDir() {
			r, err := regexp.MatchString(mask, f.Name())
			if err == nil && r {
				d := filepath.Dir(path)
				if !Contains(dirs, d) {
					dirs = append(dirs, filepath.Dir(path))
				}
			}
		}
		return nil
	})
	return dirs
}
