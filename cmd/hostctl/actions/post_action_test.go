package actions

import (
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func Test_postActionCmd(t *testing.T) {
	cmd := NewRootCmd()

	r := NewRunner(t, cmd, "postAction")
	defer r.Clean()

	t.Run("Wait and disable", func(t *testing.T) {
		r.Run("hostctl enable profile1 --wait 10ms").
			Containsf(`
				[ℹ] Using hosts file: %s

        +----------+--------+-----------+------------+
        | PROFILE  | STATUS |    IP     |   DOMAIN   |
        +----------+--------+-----------+------------+
        | profile1 | on     | 127.0.0.1 | first.loc  |
        | profile1 | on     | 127.0.0.1 | second.loc |
        +----------+--------+-----------+------------+

        [ℹ] Waiting for 10ms or ctrl+c to disable from profile 'profile1'
			`, r.Hostfile()).
			Contains(`
        +----------+--------+-----------+------------+
        | PROFILE  | STATUS |    IP     |   DOMAIN   |
        +----------+--------+-----------+------------+
        | profile1 | off    | 127.0.0.1 | first.loc  |
        | profile1 | off    | 127.0.0.1 | second.loc |
        +----------+--------+-----------+------------+
			`)
	})

	t.Run("Wait and disable on SIGTERM", func(t *testing.T) {
		proc, err := os.FindProcess(os.Getpid())
		assert.NoError(t, err)

		go func() {
			time.Sleep(10 * time.Millisecond)
			err = proc.Signal(os.Interrupt)
			assert.NoError(t, err)
		}()

		r.Run("hostctl enable profile1 --wait 0").
			Containsf(`
        [ℹ] Using hosts file: %s

        +----------+--------+-----------+------------+
        | PROFILE  | STATUS |    IP     |   DOMAIN   |
        +----------+--------+-----------+------------+
        | profile1 | on     | 127.0.0.1 | first.loc  |
        | profile1 | on     | 127.0.0.1 | second.loc |
        +----------+--------+-----------+------------+

				[ℹ] Waiting until ctrl+c to disable from profile 'profile1'
			`, r.Hostfile()).
			Contains(`
        +----------+--------+-----------+------------+
        | PROFILE  | STATUS |    IP     |   DOMAIN   |
        +----------+--------+-----------+------------+
        | profile1 | off    | 127.0.0.1 | first.loc  |
        | profile1 | off    | 127.0.0.1 | second.loc |
        +----------+--------+-----------+------------+
			`)
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
