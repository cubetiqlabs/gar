package main

import (
	"archive/tar"
	"archive/zip"
	"compress/gzip"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
	"time"

	"golang.org/x/crypto/pbkdf2"
)

const (
	Version = "1.0.0"
	// Buffer size for file operations
	BufferSize = 32 * 1024 // 32KB
)

type ArchiveFormat int

const (
	FormatZip ArchiveFormat = iota
	FormatTarGz
)

type CompressionLevel int

const (
	LevelFastest CompressionLevel = iota
	LevelNormal
	LevelBest
)

// ArchiveOptions holds configuration for archive operations
type ArchiveOptions struct {
	Format           ArchiveFormat
	CompressionLevel CompressionLevel
	Password         string
	Workers          int
	Verbose          bool
}

func main() {
	// Command line flags
	var (
		action      = flag.String("action", "", "Action: compress, extract, list")
		input       = flag.String("input", "", "Input file or directory")
		output      = flag.String("output", "", "Output file or directory")
		format      = flag.String("format", "zip", "Archive format: zip, tar.gz")
		password    = flag.String("password", "", "Password for encryption")
		compression = flag.String("compression", "normal", "Compression level: fastest, normal, best")
		workers     = flag.Int("workers", runtime.NumCPU(), "Number of worker threads")
		verbose     = flag.Bool("verbose", false, "Verbose output")
		version     = flag.Bool("version", false, "Show version")
	)

	flag.Parse()

	if *version {
		fmt.Printf("GoArchive v%s\n", Version)
		fmt.Printf("Platform: %s/%s\n", runtime.GOOS, runtime.GOARCH)
		return
	}

	if *action == "" || *input == "" {
		printUsage()
		os.Exit(1)
	}

	// Parse options
	opts := ArchiveOptions{
		Format:   parseFormat(*format),
		Password: *password,
		Workers:  *workers,
		Verbose:  *verbose,
	}

	switch *compression {
	case "fastest":
		opts.CompressionLevel = LevelFastest
	case "best":
		opts.CompressionLevel = LevelBest
	default:
		opts.CompressionLevel = LevelNormal
	}

	// Execute action
	start := time.Now()
	var err error

	switch *action {
	case "compress", "c":
		if *output == "" {
			*output = *input + getExtension(opts.Format)
		}
		err = compress(*input, *output, opts)
	case "extract", "x":
		if *output == "" {
			*output = "."
		}
		err = extract(*input, *output, opts)
	case "list", "l":
		err = listArchive(*input, opts)
	default:
		fmt.Printf("Unknown action: %s\n", *action)
		printUsage()
		os.Exit(1)
	}

	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	if opts.Verbose {
		fmt.Printf("Operation completed in %v\n", time.Since(start))
	}
}

func printUsage() {
	fmt.Println("GoArchive (gar) - High-Performance Cross-Platform Archive Manager")
	fmt.Println("\nUsage:")
	fmt.Println("  gar -action=compress -input=<path> -output=<file> [options]")
	fmt.Println("  gar -action=extract -input=<file> -output=<path> [options]")
	fmt.Println("  gar -action=list -input=<file> [options]")
	fmt.Println("\nActions:")
	fmt.Println("  compress, c    Compress files/directories")
	fmt.Println("  extract, x     Extract archive")
	fmt.Println("  list, l        List archive contents")
	fmt.Println("\nOptions:")
	flag.PrintDefaults()
}

func parseFormat(format string) ArchiveFormat {
	switch strings.ToLower(format) {
	case "tar.gz", "tgz":
		return FormatTarGz
	default:
		return FormatZip
	}
}

func getExtension(format ArchiveFormat) string {
	switch format {
	case FormatTarGz:
		return ".tar.gz"
	default:
		return ".zip"
	}
}

// Compress creates an archive from input path
func compress(inputPath, outputPath string, opts ArchiveOptions) error {
	if opts.Verbose {
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
	if opts.Password != "" {
		writer, err = newEncryptedWriter(writer, opts.Password)
		if err != nil {
			return fmt.Errorf("encryption setup: %w", err)
		}
	}

	switch opts.Format {
	case FormatZip:
		return compressZip(inputPath, info, writer, opts)
	case FormatTarGz:
		return compressTarGz(inputPath, info, writer, opts)
	default:
		return fmt.Errorf("unsupported format")
	}
}

func compressZip(inputPath string, info os.FileInfo, writer io.Writer, opts ArchiveOptions) error {
	zipWriter := zip.NewWriter(writer)
	defer zipWriter.Close()

	// Set compression level
	switch opts.CompressionLevel {
	case LevelFastest:
		zipWriter.RegisterCompressor(zip.Deflate, func(out io.Writer) (io.WriteCloser, error) {
			return gzip.NewWriterLevel(out, gzip.BestSpeed)
		})
	case LevelBest:
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
	} else {
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
}

func compressTarGz(inputPath string, info os.FileInfo, writer io.Writer, opts ArchiveOptions) error {
	// Setup gzip
	var gzLevel int
	switch opts.CompressionLevel {
	case LevelFastest:
		gzLevel = gzip.BestSpeed
	case LevelBest:
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
	} else {
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
}

// Extract extracts an archive to output path
func extract(inputPath, outputPath string, opts ArchiveOptions) error {
	if opts.Verbose {
		fmt.Printf("Extracting %s to %s...\n", inputPath, outputPath)
	}

	inFile, err := os.Open(inputPath)
	if err != nil {
		return fmt.Errorf("open archive: %w", err)
	}
	defer inFile.Close()

	var reader io.Reader = inFile

	// Check for encryption
	if opts.Password != "" {
		reader, err = newEncryptedReader(reader, opts.Password)
		if err != nil {
			return fmt.Errorf("decryption setup: %w", err)
		}
	}

	// Detect format from extension
	ext := strings.ToLower(filepath.Ext(inputPath))
	if ext == ".gz" {
		return extractTarGz(reader, outputPath, opts)
	} else {
		return extractZip(inputPath, outputPath, opts)
	}
}

func extractZip(inputPath, outputPath string, opts ArchiveOptions) error {
	zipReader, err := zip.OpenReader(inputPath)
	if err != nil {
		return err
	}
	defer zipReader.Close()

	// Use worker pool for parallel extraction
	var wg sync.WaitGroup
	sem := make(chan struct{}, opts.Workers)

	for _, file := range zipReader.File {
		wg.Add(1)
		sem <- struct{}{}

		go func(f *zip.File) {
			defer wg.Done()
			defer func() { <-sem }()

			if err := extractZipFile(f, outputPath, opts); err != nil {
				fmt.Fprintf(os.Stderr, "Error extracting %s: %v\n", f.Name, err)
			}
		}(file)
	}

	wg.Wait()
	return nil
}

func extractZipFile(f *zip.File, outputPath string, opts ArchiveOptions) error {
	destPath := filepath.Join(outputPath, f.Name)

	// Security check: prevent path traversal
	if !strings.HasPrefix(filepath.Clean(destPath), filepath.Clean(outputPath)) {
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

func extractTarGz(reader io.Reader, outputPath string, opts ArchiveOptions) error {
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

// listArchive lists contents of an archive
func listArchive(inputPath string, opts ArchiveOptions) error {
	ext := strings.ToLower(filepath.Ext(inputPath))

	switch ext {
	case ".zip":
		return listZip(inputPath)
	case ".gz":
		return listTarGz(inputPath)
	}

	return fmt.Errorf("unsupported archive format")
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

// Encryption helpers using AES-256-GCM
func newEncryptedWriter(w io.Writer, password string) (io.Writer, error) {
	// Derive key from password
	salt := make([]byte, 32)
	if _, err := rand.Read(salt); err != nil {
		return nil, err
	}

	key := pbkdf2.Key([]byte(password), salt, 100000, 32, sha256.New)

	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err := rand.Read(nonce); err != nil {
		return nil, err
	}

	// Write salt and nonce first
	if _, err := w.Write(salt); err != nil {
		return nil, err
	}
	if _, err := w.Write(nonce); err != nil {
		return nil, err
	}

	return &encryptedWriter{
		writer: w,
		gcm:    gcm,
		nonce:  nonce,
	}, nil
}

type encryptedWriter struct {
	writer io.Writer
	gcm    cipher.AEAD
	nonce  []byte
}

func (ew *encryptedWriter) Write(p []byte) (n int, err error) {
	encrypted := ew.gcm.Seal(nil, ew.nonce, p, nil)
	return ew.writer.Write(encrypted)
}

func newEncryptedReader(r io.Reader, password string) (io.Reader, error) {
	// Read salt
	salt := make([]byte, 32)
	if _, err := io.ReadFull(r, salt); err != nil {
		return nil, err
	}

	key := pbkdf2.Key([]byte(password), salt, 100000, 32, sha256.New)

	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	// Read nonce
	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(r, nonce); err != nil {
		return nil, err
	}

	return &encryptedReader{
		reader: r,
		gcm:    gcm,
		nonce:  nonce,
	}, nil
}

type encryptedReader struct {
	reader io.Reader
	gcm    cipher.AEAD
	nonce  []byte
}

func (er *encryptedReader) Read(p []byte) (n int, err error) {
	encrypted := make([]byte, len(p)+er.gcm.Overhead())
	n, err = er.reader.Read(encrypted)
	if err != nil && err != io.EOF {
		return 0, err
	}

	decrypted, err := er.gcm.Open(nil, er.nonce, encrypted[:n], nil)
	if err != nil {
		return 0, err
	}

	copy(p, decrypted)
	return len(decrypted), nil
}
