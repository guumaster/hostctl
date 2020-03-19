package host

import (
	"testing"
	"time"
)

func TestGetBackupFilename(t *testing.T) {
	t.Run("Test backup filename", func(t *testing.T) {
		nowTime := time.Now()
		filename := getBackupFilename("/etc/hosts", "/tmp", nowTime)
		if len(filename) == 0 {
			t.Fatalf("filename in empty")
		}
		checkValue := "/tmp/hosts." + nowTime.UTC().Format("20060102")
		if filename != checkValue {
			t.Fatalf("unexpected filename; got %s; want %s", filename, checkValue)
		}
	})
	t.Run("Test backup filename with last time", func(t *testing.T) {
		nowTime := time.Now().Add(-24 * time.Hour)
		filename := getBackupFilename("/tmp/path/fake_file", "/tmp/fake_path", nowTime)
		checkValue := "/tmp/fake_path/fake_file." + nowTime.UTC().Format("20060102")
		if filename != checkValue {
			t.Fatalf("unexpected filename; got %s; want %s", filename, checkValue)
		}
	})
}
