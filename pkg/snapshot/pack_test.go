package snapshot

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestPackUnpack(t *testing.T) {
	ctx := context.Background()

	err := Pack(ctx, PackOptions{
		Output:   "/tmp/test.snapshot",
		Password: "foobar",
		Sources:  []string{"./test"},
	})
	require.NoError(t, err)

	err = Unpack(ctx, UnpackOptions{
		OutputRoot:   "/tmp/restored",
		Password:     "foobar",
		SnapshotFile: "/tmp/test.snapshot",
	})
	require.NoError(t, err)
}
