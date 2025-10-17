// Package models contains shared data structures and types
package models

// ArchiveFormat defines the type of archive format
type ArchiveFormat int

const (
	FormatZip ArchiveFormat = iota
	FormatTarGz
)

// CompressionLevel defines the compression intensity
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

// CLIArgs contains parsed command-line arguments
type CLIArgs struct {
	Action      string
	Input       string
	Output      string
	Format      string
	Password    string
	Compression string
	Workers     int
	Verbose     bool
	Version     bool
	Help        bool
}
