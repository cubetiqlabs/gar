// Package archive provides compression and extraction functionality
package archive

import (
	"archive/tar"
	"compress/gzip"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/cubetiqlabs/gar/internal/models"
)

func compressTarGz(inputPath string, info os.FileInfo, writer io.Writer, opts *models.ArchiveOptions) error {
	// Setup gzip
	var gzLevel int
	switch opts.CompressionLevel {
	case models.LevelFastest:
		gzLevel = gzip.BestSpeed
	case models.LevelBest:
		gzLevel = gzip.BestCompression
	default:
		gzLevel = gzip.DefaultCompression
	}

	gzWriter, err := gzip.NewWriterLevel(writer, gzLevel)
	if err != nil {
		return err
	}
	defer gzWriter.Close()

	tarWriter := tar.NewWriter(gzWriter)
	defer tarWriter.Close()

	if info.IsDir() {
		return filepath.Walk(inputPath, func(path string, fi os.FileInfo, err error) error {
			if err != nil {
				return err
			}

			header, err := tar.FileInfoHeader(fi, fi.Name())
			if err != nil {
				return err
			}

			relPath, err := filepath.Rel(inputPath, path)
			if err != nil {
				return err
			}
			header.Name = filepath.ToSlash(relPath)

			if err := tarWriter.WriteHeader(header); err != nil {
				return err
			}

			if !fi.IsDir() {
				file, err := os.Open(path)
				if err != nil {
					return err
				}
				defer file.Close()

				if opts.Verbose {
					fmt.Printf("  Adding: %s\n", relPath)
				}

				_, err = io.Copy(tarWriter, file)
				return err
			}

			return nil
		})
	}

	// Single file
	file, err := os.Open(inputPath)
	if err != nil {
		return err
	}
	defer file.Close()

	header, err := tar.FileInfoHeader(info, info.Name())
	if err != nil {
		return err
	}
	header.Name = filepath.Base(inputPath)

	if err := tarWriter.WriteHeader(header); err != nil {
		return err
	}

	_, err = io.Copy(tarWriter, file)
	return err
}

func extractTarGz(reader io.Reader, outputPath string, opts *models.ArchiveOptions) error {
	gzReader, err := gzip.NewReader(reader)
	if err != nil {
		return err
	}
	defer gzReader.Close()

	tarReader := tar.NewReader(gzReader)

	for {
		header, err := tarReader.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}

		destPath := filepath.Join(outputPath, header.Name)

		// Security check: prevent path traversal
		if !strings.HasPrefix(filepath.Clean(destPath), filepath.Clean(outputPath)) {
			return fmt.Errorf("illegal file path: %s", header.Name)
		}

		if opts.Verbose {
			fmt.Printf("  Extracting: %s\n", header.Name)
		}

		switch header.Typeflag {
		case tar.TypeDir:
			if err := os.MkdirAll(destPath, 0755); err != nil {
				return err
			}
		case tar.TypeReg:
			if err := os.MkdirAll(filepath.Dir(destPath), 0755); err != nil {
				return err
			}

			outFile, err := os.Create(destPath)
			if err != nil {
				return err
			}

			if _, err := io.Copy(outFile, tarReader); err != nil {
				outFile.Close()
				return err
			}
			outFile.Close()

			if err := os.Chmod(destPath, os.FileMode(header.Mode)); err != nil {
				return err
			}
		}
	}

	return nil
}

func listTarGz(inputPath string) error {
	file, err := os.Open(inputPath)
	if err != nil {
		return err
	}
	defer file.Close()

	gzReader, err := gzip.NewReader(file)
	if err != nil {
		return err
	}
	defer gzReader.Close()

	tarReader := tar.NewReader(gzReader)

	fmt.Println("Archive contents:")
	for {
		header, err := tarReader.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}

		fmt.Printf("  %s (%d bytes)\n", header.Name, header.Size)
	}

	return nil
}
