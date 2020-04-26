package actions

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"testing"

	"github.com/spf13/cobra"
	as "github.com/stretchr/testify/assert"
)

func NewRunner(t *testing.T, root *cobra.Command, pattern string) Runner {
	t.Helper()

	c := &cmdRunner{t, root, "", nil}
	c.file = c.TempHostfile(pattern)

	return c
}

type Runner interface {
	Run(string) Runner
	RunE(string, error) Runner
	Runf(string, ...interface{}) Runner

	Equal(string) Runner
	Contains(string) Runner
	Containsf(string, ...interface{}) Runner
	Empty() Runner

	TempHostfile(string) *os.File
	Hostfile() string
	Clean()
}

type cmdRunner struct {
	t    *testing.T
	root *cobra.Command
	out  string
	file *os.File
}

func (c *cmdRunner) Equal(expected string) Runner {
	expected = trimLeft(expected)
	as.Equal(c.t, expected, c.out)

	return c
}

func (c *cmdRunner) Hostfile() string {
	return c.file.Name()
}

func trimLeft(s string) string {
	lines := strings.Split(s, "\n")
	for i, l := range lines {
		lines[i] = strings.TrimLeft(strings.ReplaceAll(l, "\t", " "), " ")
	}

	return strings.Join(lines, "\n")
}

func (c *cmdRunner) Contains(expected string) Runner {
	expected = trimLeft(expected)
	as.Contains(c.t, c.out, expected)

	return c
}

func (c *cmdRunner) Containsf(expected string, args ...interface{}) Runner {
	expected = fmt.Sprintf(expected, args...)
	expected = trimLeft(expected)
	as.Contains(c.t, c.out, expected)

	return c
}

func (c *cmdRunner) Empty() Runner {
	as.Empty(c.t, strings.ReplaceAll(c.out, "\n", ""))

	return c
}

func (c *cmdRunner) Runf(format string, args ...interface{}) Runner {
	return c.Run(fmt.Sprintf(format, args...))
}

func (c *cmdRunner) Run(cmd string) Runner {
	assert := as.New(c.t)
	b := bytes.NewBufferString("")

	c.out = ""
	c.root.SetOut(b)

	if !strings.Contains(cmd, "--host-file") {
		cmd += fmt.Sprintf(" --host-file %s", c.file.Name())
	}

	args := strings.Split(cmd, " ")
	args = args[1:]

	c.root.SetArgs(args)

	err := c.root.Execute()
	assert.NoError(err)

	out, err := ioutil.ReadAll(b)
	assert.NoError(err)

	c.out = "\n" + string(out)

	return c
}

func (c *cmdRunner) RunE(cmd string, expectedErr error) Runner {
	assert := as.New(c.t)
	b := bytes.NewBufferString("")

	c.out = ""
	c.root.SetOut(b)

	cmd += fmt.Sprintf(" --host-file %s", c.file.Name())

	args := strings.Split(cmd, " ")
	args = args[1:]

	c.root.SetArgs(args)

	actualErr := c.root.Execute()
	assert.EqualError(actualErr, expectedErr.Error())

	out, err := ioutil.ReadAll(b)
	assert.NoError(err)

	c.out = "\n" + string(out)

	return c
}

func (c *cmdRunner) TempHostfile(pattern string) *os.File {
	file, err := ioutil.TempFile("/tmp", fmt.Sprintf("%s_%s_", c.root.Name(), pattern))
	as.NoError(c.t, err)

	_, _ = file.WriteString(`
127.0.0.1 localhost

# profile.on profile1
127.0.0.1 first.loc
127.0.0.1 second.loc
# end

# profile.off profile2
# 127.0.0.1 first.loc
# 127.0.0.1 second.loc
# end
`)

	return file
}

func (c *cmdRunner) Clean() {
	_ = c.file.Close()
	_ = os.Remove(c.file.Name())
}
