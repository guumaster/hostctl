package actions

import (
	"testing"
)

func Test_Info(t *testing.T) {
	cmd := NewRootCmd()

	r := NewRunner(t, cmd, "info")
	defer r.Clean()

	r.Run("hostctl info").Contains("Your dev tool to manage /etc/hosts like a pro")
}
