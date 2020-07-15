package actions

import (
	"bytes"
	"context"
	"io"
	"io/ioutil"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/trusch/stackctl/pkg/compose"
)

func TestCreatePod(t *testing.T) {
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
    ports:
      - 127.0.0.1:8080:80/tcp
`)))
	require.NoError(t, err)
	require.NoError(t, f.Close())
	defer os.Remove(f.Name())
	project, err := compose.Load(f.Name())
	require.NoError(t, err)
	project.Name = "test"

	err = CreatePod(context.Background(), project)
	require.NoError(t, err)
}
