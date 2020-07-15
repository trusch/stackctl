package compose

import (
	"bytes"
	"io"
	"io/ioutil"
	"os"
	"testing"

	"github.com/davecgh/go-spew/spew"
	"github.com/stretchr/testify/require"
)

func TestBasicLoad(t *testing.T) {
	f, err := ioutil.TempFile(os.TempDir(), "test")
	require.NoError(t, err)
	_, err = io.Copy(f, bytes.NewReader([]byte(`
version: "3.9"
x-project-name: testproject
x-my-global-extension: foobar
services:
  foo:
    image: alpine:latest
    command: ["tail","-f","/dev/null"]
    x-my-extension: foobar
`)))
	require.NoError(t, err)
	require.NoError(t, f.Close())
	defer os.Remove(f.Name())

	project, err := Load(f.Name())
	require.NoError(t, err)
	spew.Dump(project)
}
