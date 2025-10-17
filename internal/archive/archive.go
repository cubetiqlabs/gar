// Package archive provides compression and extraction functionality for multiple formats
package archive

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/cubetiqlabs/gar/internal/crypto"
	"github.com/cubetiqlabs/gar/internal/models"
)

// Operator handles archive operations (compress, extract, list)
type Operator struct {
	opts *models.ArchiveOptions
}

// NewOperator creates a new archive operator
func NewOperator(opts *models.ArchiveOptions) *Operator {
	return &Operator{opts: opts}
}

// Compress creates an archive from input path
func (op *Operator) Compress(inputPath, outputPath string) error {
	if op.opts.Verbose {
		fmt.Printf("Compressing %s to %s...\n", inputPath, outputPath)
	}

	// Check if input exists
	info, err := os.Stat(inputPath)
	if err != nil {
		return fmt.Errorf("input path error: %w", err)
	}

	// Create output file
	outFile, err := os.Create(outputPath)
	if err != nil {
		return fmt.Errorf("create output file: %w", err)
	}
	defer outFile.Close()

	var writer io.Writer = outFile

	// Add encryption if password is provided
	if op.opts.Password != "" {
		writer, err = crypto.NewEncryptedWriter(writer, op.opts.Password)
		if err != nil {
			return fmt.Errorf("encryption setup: %w", err)
		}
	}

	switch op.opts.Format {
	case models.FormatZip:
		return compressZip(inputPath, info, writer, op.opts)
	case models.FormatTarGz:
		return compressTarGz(inputPath, info, writer, op.opts)
	default:
		return fmt.Errorf("unsupported format")
	}
}

// Extract extracts an archive to output path
func (op *Operator) Extract(inputPath, outputPath string) error {
	if op.opts.Verbose {
		fmt.Printf("Extracting %s to %s...\n", inputPath, outputPath)
	}

	inFile, err := os.Open(inputPath)
	if err != nil {
		return fmt.Errorf("open archive: %w", err)
	}
	defer inFile.Close()

	var reader io.Reader = inFile

	// Check for encryption
	if op.opts.Password != "" {
		reader, err = crypto.NewEncryptedReader(reader, op.opts.Password)
		if err != nil {
			return fmt.Errorf("decryption setup: %w", err)
		}
	}

	// Detect format from extension
	ext := strings.ToLower(filepath.Ext(inputPath))
	if ext == ".gz" {
		return extractTarGz(reader, outputPath, op.opts)
	}
	return extractZip(inputPath, outputPath, op.opts)
}

// List lists archive contents
func (op *Operator) List(inputPath string) error {
	ext := strings.ToLower(filepath.Ext(inputPath))

	switch ext {
	case ".zip":
		return listZip(inputPath)
	case ".gz":
		return listTarGz(inputPath)
	}

	return fmt.Errorf("unsupported archive format: %s", ext)
}

// ParseFormat converts string to ArchiveFormat
func ParseFormat(format string) models.ArchiveFormat {
	switch strings.ToLower(format) {
	case "tar.gz", "tgz":
		return models.FormatTarGz
	default:
		return models.FormatZip
	}
}

// GetExtension returns the file extension for a given format
func GetExtension(format models.ArchiveFormat) string {
	switch format {
	case models.FormatTarGz:
		return ".tar.gz"
	default:
		return ".zip"
	}
}

// TimeOperation measures the time taken for an operation
func TimeOperation(fn func() error, verbose bool, operationName string) error {
	start := time.Now()
	err := fn()

	if verbose && err == nil {
		fmt.Printf("%s completed in %v\n", operationName, time.Since(start))
	}

	return err
}
