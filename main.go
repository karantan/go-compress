package main

import (
	"archive/tar"
	"compress/gzip"
	"gocompress/logger"
	"gocompress/utils"
	"io"
	"os"
	"path"

	"github.com/klauspost/compress/zstd"
	"github.com/spf13/cobra"
)

var (
	log         = logger.New("main", true)
	source      string
	compression string
)

var cmd = &cobra.Command{
	Use:   "gocompress",
	Short: "Archive and compress a file/folder",
	Long: `A very simple archiver with a few compression options. It's not meant for
production usage. Mainly for testing.
`,
}

func init() {
	cmd.Run = runArchiving
	cmd.Flags().StringVarP(&source, "path", "p", "", "path to the folder you want to archive")
	cmd.MarkFlagRequired("path")
	cmd.Flags().StringVarP(&compression, "compression", "c", ".tar.zst", "compression method. Options: '.tar.zst', '.tar.gz' or '.tar'")
}

func runArchiving(cmd *cobra.Command, args []string) {
	if _, err := os.Stat(source); os.IsNotExist(err) {
		log.Errorf("%s does not exist.", source)
		return
	}

	archiveFile := path.Join(utils.RootDir(), "archive"+compression)
	buf, _ := os.Create(archiveFile)
	defer buf.Close()

	var err error

	switch compression {
	case ".tar":
		err = Archive(source, buf)
	case ".tar.gz":
		enc := gzip.NewWriter(buf)
		defer enc.Close()
		err = Archive(source, enc)
	default:
		enc, err := zstd.NewWriter(buf)
		if err != nil {
			log.Error(err)
			return
		}
		defer enc.Close()
		err = Archive(source, enc)
	}

	if err != nil {
		log.Error(err)
	}

}

func main() {
	cmd.Execute()
}

// Archive  function takes a source file/folder and creates an archive in the
// `targetDir` folder. Last argument is the encoder part. It can be e.g.
// `gzip.NewWriter`, `zstd.NewWriter` or even `os.File` and `bytes.Buffer`
// which won't do any encoding/compressing.
// Keep in mind that archiving is not the same as compressing. Tape archive (tar) is
// a file format for storing a sequence of files that can be read and written in a
// streaming manner. Once a tar archive is created we can compress it.
func Archive(sourceDir string, encoder io.Writer) error {
	files, _ := utils.FilePathWalkDir(sourceDir)

	tw := tar.NewWriter(encoder)
	defer tw.Close()

	// Iterate over files and add them to the tar archive
	for _, file := range files {
		err := addToArchive(tw, file)
		if err != nil {
			return err
		}
	}

	return nil
}

func addToArchive(tw *tar.Writer, filename string) error {
	// Open the file which will be written into the archive
	file, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	// Get FileInfo about our file providing file size, mode, etc.
	info, err := file.Stat()
	if err != nil {
		return err
	}

	// Create a tar Header from the FileInfo data
	header, err := tar.FileInfoHeader(info, info.Name())
	if err != nil {
		return err
	}

	// Use full path as name (FileInfoHeader only takes the basename)
	// If we don't do this the directory structure would not be preserved
	// https://golang.org/src/archive/tar/common.go?#L626
	header.Name = filename

	// Write file header to the tar archive
	if err = tw.WriteHeader(header); err != nil {
		return err
	}

	// Copy file content to tar archive
	_, err = io.Copy(tw, file)
	if err != nil {
		return err
	}

	return nil
}
