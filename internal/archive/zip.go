// Package archive provides compression and extraction functionality
package archive

import (
	"archive/zip"
	"compress/gzip"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/cubetiqlabs/gar/internal/models"
)

func compressZip(inputPath string, info os.FileInfo, writer io.Writer, opts *models.ArchiveOptions) error {
	zipWriter := zip.NewWriter(writer)
	defer zipWriter.Close()

	// Set compression level
	switch opts.CompressionLevel {
	case models.LevelFastest:
		zipWriter.RegisterCompressor(zip.Deflate, func(out io.Writer) (io.WriteCloser, error) {
			return gzip.NewWriterLevel(out, gzip.BestSpeed)
		})
	case models.LevelBest:
		zipWriter.RegisterCompressor(zip.Deflate, func(out io.Writer) (io.WriteCloser, error) {
			return gzip.NewWriterLevel(out, gzip.BestCompression)
		})
	}

	if info.IsDir() {
		return filepath.Walk(inputPath, func(path string, fi os.FileInfo, err error) error {
			if err != nil {
				return err
			}

			header, err := zip.FileInfoHeader(fi)
			if err != nil {
				return err
			}

			relPath, err := filepath.Rel(inputPath, path)
			if err != nil {
				return err
			}
			header.Name = filepath.ToSlash(relPath)

			if fi.IsDir() {
				header.Name += "/"
			} else {
				header.Method = zip.Deflate
			}

			w, err := zipWriter.CreateHeader(header)
			if err != nil {
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

				_, err = io.Copy(w, file)
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

	header, err := zip.FileInfoHeader(info)
	if err != nil {
		return err
	}
	header.Name = filepath.Base(inputPath)
	header.Method = zip.Deflate

	w, err := zipWriter.CreateHeader(header)
	if err != nil {
		return err
	}

	_, err = io.Copy(w, file)
	return err
}

func extractZip(inputPath, outputPath string, opts *models.ArchiveOptions) error {
	zipReader, err := zip.OpenReader(inputPath)
	if err != nil {
		return err
	}
	defer zipReader.Close()

	// Use worker pool for parallel extraction
	var wg sync.WaitGroup
	sem := make(chan struct{}, opts.Workers)
	errChan := make(chan error, 1)

	for _, file := range zipReader.File {
		wg.Add(1)
		sem <- struct{}{}

		go func(f *zip.File) {
			defer wg.Done()
			defer func() { <-sem }()

			if err := extractZipFile(f, outputPath, opts); err != nil {
				select {
				case errChan <- err:
				default:
				}
				fmt.Fprintf(os.Stderr, "Error extracting %s: %v\n", f.Name, err)
			}
		}(file)
	}

	wg.Wait()

	// Check if there was any error
	select {
	case err := <-errChan:
		return err
	default:
		return nil
	}
}

func extractZipFile(f *zip.File, outputPath string, opts *models.ArchiveOptions) error {
	destPath := filepath.Join(outputPath, f.Name)

	// Security check: prevent path traversal
	// Convert both paths to absolute to handle relative paths like "." correctly
	absDestPath, err := filepath.Abs(destPath)
	if err != nil {
		return fmt.Errorf("invalid destination path: %s", f.Name)
	}
	absOutputPath, err := filepath.Abs(outputPath)
	if err != nil {
		return fmt.Errorf("invalid output path: %s", outputPath)
	}
	if !strings.HasPrefix(absDestPath, absOutputPath+string(filepath.Separator)) && absDestPath != absOutputPath {
		return fmt.Errorf("illegal file path: %s", f.Name)
	}

	if f.FileInfo().IsDir() {
		return os.MkdirAll(destPath, f.Mode())
	}

	if opts.Verbose {
		fmt.Printf("  Extracting: %s\n", f.Name)
	}

	// Create parent directories
	if err := os.MkdirAll(filepath.Dir(destPath), 0755); err != nil {
		return err
	}

	// Extract file
	rc, err := f.Open()
	if err != nil {
		return err
	}
	defer rc.Close()

	outFile, err := os.OpenFile(destPath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
	if err != nil {
		return err
	}
	defer outFile.Close()

	_, err = io.Copy(outFile, rc)
	return err
}

func listZip(inputPath string) error {
	zipReader, err := zip.OpenReader(inputPath)
	if err != nil {
		return err
	}
	defer zipReader.Close()

	fmt.Println("Archive contents:")
	for _, f := range zipReader.File {
		fmt.Printf("  %s (%d bytes)\n", f.Name, f.UncompressedSize64)
	}

	return nil
}
