package service

import (
	"path"

	"github.com/spf13/afero"
)

var FS = afero.NewOsFs()

// Creates a new or overwrites an existing file with the specified content
// Note: if a path to the file is given, all sub folders on the path are also created
func WriteToFile(data []byte, filename string) error {
	if err := FS.MkdirAll(path.Dir(filename), 0755); err != nil {
		return err
	}

	file, err := FS.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	if _, err = file.Write(data); err != nil {
		return err
	}
	return nil
}
