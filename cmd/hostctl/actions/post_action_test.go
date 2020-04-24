package actions

import (
	"bytes"
	"io/ioutil"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func Test_postActionCmd(t *testing.T) {
	cmd := NewRootCmd()

	tmp := makeTempHostsFile(t, "postActionCmd")
	defer os.Remove(tmp.Name())

	t.Run("Wait and disable", func(t *testing.T) {
		b := bytes.NewBufferString("")
		args := []string{"enable", "profile1", "--host-file", tmp.Name(), "--wait", "10ms"}

		cmd.SetOut(b)
		cmd.SetArgs(args)

		err := cmd.Execute()
		assert.NoError(t, err)

		out, _ := ioutil.ReadAll(b)
		assert.Contains(t, string(out), "Waiting for 10ms or ctrl+c to disable from profile 'profile1'")
	})

	t.Run("Wait and disable on SIGTERM", func(t *testing.T) {
		b := bytes.NewBufferString("")
		args := []string{"enable", "profile1", "--host-file", tmp.Name(), "--wait", "0"}

		cmd.SetOut(b)
		cmd.SetArgs(args)

		proc, err := os.FindProcess(os.Getpid())
		assert.NoError(t, err)

		go func() {
			time.Sleep(10 * time.Millisecond)
			err = proc.Signal(os.Interrupt)
			assert.NoError(t, err)
		}()

		err = cmd.Execute()
		assert.NoError(t, err)

		out, _ := ioutil.ReadAll(b)
		assert.Contains(t, string(out), "Waiting until ctrl+c to disable from profile 'profile1'")
	})
}

func Test_waitSignalOrDuration(t *testing.T) {
	t.Run("timeout", func(t *testing.T) {
		d := 10 * time.Millisecond
		done := waitSignalOrDuration(d)
		v, ok := <-done
		assert.Equal(t, v, struct{}{})
		assert.Equal(t, ok, true)
	})

	t.Run("timeout abs", func(t *testing.T) {
		d := -10 * time.Millisecond
		done := waitSignalOrDuration(d)
		v, ok := <-done
		assert.Equal(t, v, struct{}{})
		assert.Equal(t, ok, true)
	})

	t.Run("SIGINT", func(t *testing.T) {
		proc, err := os.FindProcess(os.Getpid())
		assert.NoError(t, err)

		d := 1 * time.Hour
		done := waitSignalOrDuration(d)

		time.Sleep(10 * time.Millisecond)
		err = proc.Signal(os.Interrupt)
		assert.NoError(t, err)

		v, ok := <-done
		assert.Equal(t, v, struct{}{})
		assert.Equal(t, ok, true)
	})

	t.Run("Wait forever", func(t *testing.T) {
		proc, err := os.FindProcess(os.Getpid())
		assert.NoError(t, err)

		done := waitSignalOrDuration(0)

		time.Sleep(10 * time.Millisecond)
		err = proc.Signal(os.Interrupt)
		assert.NoError(t, err)

		v, ok := <-done
		assert.Equal(t, v, struct{}{})
		assert.Equal(t, ok, true)
	})
}
