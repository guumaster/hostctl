package host

import (
	"io"
	"os"
)

// Restore overwrite content of a hosts file with the content of a backup.
func (f *File) Restore(from string) error {
	fromFile, err := f.fs.Open(from)
	if err != nil {
		return err
	}
	defer fromFile.Close()

	dstFile, err := f.fs.OpenFile(f.src.Name(), os.O_CREATE|os.O_RDWR|os.O_TRUNC, 0644)
	if err != nil {
		return err
	}
	defer dstFile.Close()

	_, err = io.Copy(dstFile, fromFile)

	return err
}
