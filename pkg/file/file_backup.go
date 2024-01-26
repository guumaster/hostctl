package file

import (
	"fmt"
	"io"
	"os"
	"path"
	"time"
)

// Backup creates a copy of your hosts file to a new location with the date as
// extension. It recursively creates the target directory if it does not exist.
func (f *File) Backup(dst string) (string, error) {
	_, _ = f.src.Seek(0, io.SeekStart)
	bkpFilename := fmt.Sprintf("%s.%s", f.src.Name(), time.Now().UTC().Format("20060102"))
	bkpFilename = path.Join(dst, path.Base(bkpFilename))

	// check if directory exists, else make it
	if _, err := f.fs.Stat(dst); os.IsNotExist(err) {
		err := f.fs.MkdirAll(dst, os.ModePerm)
		if err != nil {
			return "", err
		}
	}

	b, err := f.fs.Create(bkpFilename)
	if err != nil {
		return "", err
	}
	defer b.Close()

	_, err = io.Copy(b, f.src)

	return bkpFilename, err
}
