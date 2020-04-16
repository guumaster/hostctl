package host

import (
	"fmt"
	"io"
	"path"
	"time"
)

// Backup creates a copy of your hosts file to a new location with the date as extension
func (f *File) Backup(dst string) (string, error) {
	_, _ = f.src.Seek(0, io.SeekStart)
	bkpFilename := fmt.Sprintf("%s.%s", f.src.Name(), time.Now().UTC().Format("20060102"))
	bkpFilename = path.Join(dst, path.Base(bkpFilename))

	b, err := f.fs.Create(bkpFilename)
	if err != nil {
		return "", err
	}
	defer b.Close()

	_, err = io.Copy(b, f.src)

	return bkpFilename, err
}
