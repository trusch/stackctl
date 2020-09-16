package snapshot

import (
	"archive/tar"
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/compose-spec/compose-go/types"
	"github.com/dolmen-go/contextio"
	"github.com/sirupsen/logrus"
	"github.com/ulikunitz/xz"
	"golang.org/x/crypto/openpgp"
)

type PackOptions struct {
	Output   string   // path to output file
	Password string   // eventually a passphrase
	Sources  []string // list of files/directories to include in the snapshot
}

type UnpackOptions struct {
	OutputRoot   string
	Password     string
	SnapshotFile string
}

func Create(ctx context.Context, project *types.Project, volumes []string, password, output string) error {
	args := []string{"run", "-v", ".:/.working-directory"}
	for _, volName := range volumes {
		if _, ok := project.Volumes[volName]; ok {
			args = append(args, "-v", fmt.Sprintf("%s:/%s", volName, volName))
			continue
		}
		if vol, ok := project.Configs[volName]; ok {
			args = append(args, "-v", fmt.Sprintf("%s:/%s", vol.File, volName))
			continue
		}
		if vol, ok := project.Secrets[volName]; ok {
			args = append(args, "-v", fmt.Sprintf("%s:/%s", vol.File, volName))
			continue
		}
		return fmt.Errorf("volume %s not found", volName)
	}
	args = append(args, "containers.trusch.io/stackctl", "stackctl", "snapshot", "pack", "--output="+filepath.Join("/.working-directory", output))
	if password != "" {
		args = append(args, "--password", password)
	}
	for _, volName := range volumes {
		args = append(args, "/"+volName)
	}
	c := exec.Command("podman", args...)
	c.Stdout = os.Stdout
	c.Stdin = os.Stdin
	c.Stderr = os.Stderr
	return c.Run()
}

func Load(ctx context.Context, project *types.Project, volumes []string, password, snapshotFile string) error {
	args := []string{"run", "-v", snapshotFile + ":/.snapshot"}
	for _, volName := range volumes {
		if _, ok := project.Volumes[volName]; ok {
			args = append(args, "-v", fmt.Sprintf("%s:/%s", volName, volName))
			continue
		}
		if vol, ok := project.Configs[volName]; ok {
			args = append(args, "-v", fmt.Sprintf("%s:/%s", vol.File, volName))
			continue
		}
		if vol, ok := project.Secrets[volName]; ok {
			args = append(args, "-v", fmt.Sprintf("%s:/%s", vol.File, volName))
			continue
		}
		return fmt.Errorf("volume %s not found", volName)
	}
	args = append(args, "containers.trusch.io/stackctl", "stackctl", "snapshot", "unpack", "--input=/.snapshot", "--output=/")
	if password != "" {
		args = append(args, "--password", password)
	}
	c := exec.Command("podman", args...)
	c.Stdout = os.Stdout
	c.Stdin = os.Stdin
	c.Stderr = os.Stderr
	return c.Run()
}

func Pack(ctx context.Context, opts PackOptions) (err error) {
	var writer io.WriteCloser
	f, err := os.Create(opts.Output)
	if err != nil {
		return err
	}
	defer f.Close()

	xzWriter, err := xz.NewWriter(f)
	if err != nil {
		return err
	}
	defer xzWriter.Close()

	writer = xzWriter

	if opts.Password != "" {
		encrypter, err := openpgp.SymmetricallyEncrypt(writer, []byte(opts.Password), nil, nil)
		if err != nil {
			return err
		}
		defer encrypter.Close()
		writer = encrypter
	}

	tw := tar.NewWriter(writer)
	defer tw.Close()
	for _, s := range opts.Sources {
		st, err := os.Stat(s)
		switch {
		case err != nil:
			return err
		case st.IsDir():
			err := filepath.Walk(st.Name(), func(path string, info os.FileInfo, err error) error {
				if err != nil {
					return err
				}
				return addFileToTarWriter(ctx, path, info, tw)
			})
			if err != nil {
				return err
			}
		case st.Mode().IsRegular():
			return addFileToTarWriter(ctx, s, st, tw)
		default:
			logrus.Warnf("failed to process %s", s)
		}
	}
	return nil
}

func Unpack(ctx context.Context, opts UnpackOptions) error {
	var reader io.Reader
	f, err := os.Open(opts.SnapshotFile)
	if err != nil {
		return err
	}
	defer f.Close()

	xzReader, err := xz.NewReader(f)
	if err != nil {
		return err
	}

	reader = xzReader

	if opts.Password != "" {
		alreadyPrompted := false
		md, err := openpgp.ReadMessage(reader, nil, func(keys []openpgp.Key, symmetric bool) ([]byte, error) {
			if alreadyPrompted {
				return nil, errors.New("Could not decrypt data using supplied passphrase")
			} else {
				alreadyPrompted = true
			}
			return []byte(opts.Password), nil
		}, nil)
		if err != nil {
			return err
		}
		reader = md.UnverifiedBody
	}

	tr := tar.NewReader(reader)
	ctxTarReader := contextio.NewReader(ctx, tr)
	for {
		header, err := tr.Next()
		if err == io.EOF {
			break // End of archive
		}
		if err != nil {
			return err
		}
		target := filepath.Join(opts.OutputRoot, header.Name)
		// check the type
		switch header.Typeflag {
		// if its a dir and it doesn't exist create it (with 0755 permission)
		case tar.TypeDir:
			if _, err := os.Stat(target); err != nil {
				if err := os.MkdirAll(target, 0755); err != nil {
					return err
				}
			}
		// if it's a file create it (with same permission)
		case tar.TypeReg:
			fileToWrite, err := os.OpenFile(target, os.O_CREATE|os.O_RDWR, os.FileMode(header.Mode))
			if err != nil {
				return err
			}
			// copy over contents
			if _, err := io.Copy(contextio.NewWriter(ctx, fileToWrite), ctxTarReader); err != nil {
				return err
			}
			// manually close here after each file operation; defering would cause each file close
			// to wait until all operations have completed.
			fileToWrite.Close()
		}
	}
	return nil
}

func addFileToTarWriter(ctx context.Context, filePath string, stat os.FileInfo, tarWriter *tar.Writer) error {
	logrus.Debugf("add '%s' to tar", filePath)
	if err := ctx.Err(); err != nil {
		return err
	}

	header, err := tar.FileInfoHeader(stat, "")
	if err != nil {
		return err
	}
	header.Name = filePath

	err = tarWriter.WriteHeader(header)
	if err != nil {
		return errors.New(fmt.Sprintf("Could not write header for file '%s', got error '%s'", filePath, err.Error()))
	}

	if stat.Mode().IsRegular() {
		file, err := os.Open(filePath)
		if err != nil {
			return errors.New(fmt.Sprintf("Could not open file '%s', got error '%s'", filePath, err.Error()))
		}

		_, err = io.Copy(contextio.NewWriter(ctx, tarWriter), contextio.NewReader(ctx, file))
		if err != nil {
			return errors.New(fmt.Sprintf("Could not copy the file '%s' data to the tarball, got error '%s'", filePath, err.Error()))
		}
		file.Close()
	}

	return nil
}
