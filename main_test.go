package main

import (
	"compress/gzip"
	"gocompress/utils"
	"os"
	"path"
	"testing"

	"github.com/klauspost/compress/zstd"
	"github.com/stretchr/testify/assert"
)

func TestArchive(t *testing.T) {
	assert := assert.New(t)
	fooFolder := path.Join(utils.RootDir(), "fixtures", "foo")

	t.Run("create a tar archive", func(t *testing.T) {
		archiveFile := path.Join(utils.RootDir(), "archive.tar")
		buf, _ := os.Create(archiveFile)
		defer buf.Close()
		defer os.Remove(buf.Name())

		err := Archive(fooFolder, buf)
		assert.NoError(err)

	})
	t.Run("create a zstd tar archive", func(t *testing.T) {
		archiveFile := path.Join(utils.RootDir(), "archive.tar.zst")
		buf, _ := os.Create(archiveFile)
		defer buf.Close()
		defer os.Remove(buf.Name())
		enc, _ := zstd.NewWriter(buf)
		defer enc.Close()

		err := Archive(fooFolder, enc)
		assert.NoError(err)
	})
	t.Run("create a gzip tar archive", func(t *testing.T) {
		archiveFile := path.Join(utils.RootDir(), "archive.tar.gz")
		buf, _ := os.Create(archiveFile)
		defer buf.Close()
		defer os.Remove(buf.Name())
		enc := gzip.NewWriter(buf)
		defer enc.Close()

		err := Archive(fooFolder, enc)
		assert.NoError(err)
		assert.FileExists(archiveFile)
	})

}

func BenchmarkArchiveTar(b *testing.B) {
	fooFolder := path.Join(utils.RootDir(), "fixtures", "foo")
	archiveFile := path.Join(utils.RootDir(), "archive.tar.tar")
	buf, _ := os.Create(archiveFile)
	defer buf.Close()
	defer os.Remove(buf.Name())

	for i := 0; i < b.N; i++ {
		Archive(fooFolder, buf)
	}
}

func BenchmarkArchiveZstd(b *testing.B) {
	fooFolder := path.Join(utils.RootDir(), "fixtures", "foo")
	archiveFile := path.Join(utils.RootDir(), "archive.tar.zst")
	buf, _ := os.Create(archiveFile)
	defer buf.Close()
	defer os.Remove(buf.Name())

	enc, _ := zstd.NewWriter(buf)
	defer enc.Close()

	for i := 0; i < b.N; i++ {
		Archive(fooFolder, enc)
	}
}

func BenchmarkArchiveGzip(b *testing.B) {
	fooFolder := path.Join(utils.RootDir(), "fixtures", "foo")
	archiveFile := path.Join(utils.RootDir(), "archive.tar.gz")
	buf, _ := os.Create(archiveFile)
	defer buf.Close()
	defer os.Remove(buf.Name())

	enc := gzip.NewWriter(buf)
	defer enc.Close()

	for i := 0; i < b.N; i++ {
		Archive(fooFolder, enc)
	}
}
