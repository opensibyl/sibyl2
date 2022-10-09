package model

import "path/filepath"

func PathStandardize(results []*FileResult, basedir string) error {
	for _, each := range results {
		newPath, err := filepath.Rel(basedir, each.Path)
		if err != nil {
			return err
		}

		each.Path = filepath.ToSlash(newPath)
	}
	return nil
}
