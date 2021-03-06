// Copyright (c) 2018 Vincent Landgraf

package ghostscript

import (
	"context"
	"crypto/md5"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestThumbnailer13(t *testing.T) {
	gen := func(ctx context.Context, page int) (io.WriteCloser, error) {
		file := filepath.Join("test", fmt.Sprintf("gen-%d.png", page))
		return os.Create(file)
	}

	source, err := os.Open("test/input-1.3.pdf")
	assert.NoError(t, err)
	defer source.Close()

	err = DefaultConfig.NewThumbnailerContext(context.Background(), source, 70, gen)
	assert.NoError(t, err)

	assert.FileExists(t, "test/gen-1.png")
	md5, err := md5File("test/gen-1.png")
	assert.NoError(t, err)
	assert.Equal(t, "76e983bf9a247a84f706d74695de03f6", md5)

	assert.FileExists(t, "test/gen-2.png")
	md5, err = md5File("test/gen-2.png")
	assert.NoError(t, err)
	assert.Equal(t, "45dc970bd4764b8295cbc71816deb585", md5)

	os.Remove("test/gen-1.png")
	os.Remove("test/gen-2.png")
}

type nopWriteCloser struct{}

func (c *nopWriteCloser) Write(p []byte) (n int, err error) { return len(p), nil }
func (c *nopWriteCloser) Close() error                      { return nil }

func TestThumbnailer15(t *testing.T) {
	i := 0
	gen := func(ctx context.Context, page int) (io.WriteCloser, error) {
		i++
		return &nopWriteCloser{}, nil
	}

	source, err := os.Open("test/input-1.5.pdf")
	assert.NoError(t, err)
	defer source.Close()

	err = DefaultConfig.NewThumbnailerContext(context.Background(), source, 70, gen)
	assert.NoError(t, err)
	assert.Equal(t, 81, i)
}

func md5File(path string) (string, error) {
	f, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer f.Close()

	h := md5.New()
	_, err = io.Copy(h, f)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("%x", h.Sum(nil)), nil
}
