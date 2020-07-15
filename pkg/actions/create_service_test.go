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

func TestCreateService(t *testing.T) {
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
      - 8080:80
    volumes:
      - vol-1:/vol-1
      - vol-2:/vol-2
      - ./:/vol-3
    configs:
      - cfg-1
    secrets:
      - secret-1
volumes:
  vol-1:
  vol-2:
configs:
  cfg-1:
    file: ./
secrets:
  secret-1:
    file: ./

`)))

	require.NoError(t, err)

	require.NoError(t, f.Close())
	defer os.Remove(f.Name())
	project, err := compose.Load(f.Name())
	require.NoError(t, err)
	project.Name = "test"
	ctx := context.Background()

	RemovePod(ctx, project)
	err = CreatePod(ctx, project)
	require.NoError(t, err)
	err = CreateVolumes(ctx, project)
	require.NoError(t, err)
	err = CreateService(ctx, project, "foo")
	require.NoError(t, err)
}
