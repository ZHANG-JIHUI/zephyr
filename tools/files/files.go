package files

import (
	"os"
	"path/filepath"
)

func FileList(dir, ext string, recursive bool) ([]string, error) {
	var list []string
	if recursive {
		if err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if !info.IsDir() {
				if ext != "" && !(filepath.Ext(path) == ext) {
					return nil
				}
				file, err := filepath.Abs(path)
				if err != nil {
					return err
				}
				list = append(list, file)
			}
			return nil
		}); err != nil {
			return nil, err
		}
	} else {
		entries, err := os.ReadDir(dir)
		if err != nil {
			return nil, err
		}
		for _, entry := range entries {
			if entry.IsDir() {
				continue
			}
			if ext != "" && !(filepath.Ext(entry.Name()) == ext) {
				continue
			}
			file, err := filepath.Abs(entry.Name())
			if err != nil {
				return nil, err
			}
			list = append(list, file)
		}
	}

	return list, nil
}
