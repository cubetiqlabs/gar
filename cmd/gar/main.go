package main

import (
	"fmt"
	"os"

	"github.com/cubetiqlabs/gar/internal/archive"
	"github.com/cubetiqlabs/gar/internal/cli"
	"github.com/cubetiqlabs/gar/internal/models"
	"github.com/cubetiqlabs/gar/pkg/version"
)

func main() {
	// Create CLI parser
	parser := cli.NewParser()

	// Parse arguments
	args, err := parser.Parse(os.Args[1:])
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error parsing arguments: %v\n", err)
		os.Exit(1)
	}

	// Get version from pkg/version
	Version := version.Number()

	// Handle help
	if args.Help {
		parser.PrintUsage(Version)
		os.Exit(0)
	}

	// Handle version
	if args.Version {
		fmt.Printf("GoArchive v%s\n", Version)
		return
	}

	// Validate required arguments
	if args.Action == "" || args.Input == "" {
		parser.PrintUsage(Version)
		os.Exit(1)
	}

	// Build archive options from parsed arguments
	opts := &models.ArchiveOptions{
		Format:   archive.ParseFormat(args.Format),
		Password: args.Password,
		Workers:  args.Workers,
		Verbose:  args.Verbose,
	}

	// Parse compression level
	switch args.Compression {
	case "fastest":
		opts.CompressionLevel = models.LevelFastest
	case "best":
		opts.CompressionLevel = models.LevelBest
	default:
		opts.CompressionLevel = models.LevelNormal
	}

	// Execute action
	operator := archive.NewOperator(opts)
	var actionErr error

	switch args.Action {
	case "compress", "c":
		output := args.Output
		if output == "" {
			output = args.Input + archive.GetExtension(opts.Format)
		}
		actionErr = archive.TimeOperation(
			func() error { return operator.Compress(args.Input, output) },
			opts.Verbose,
			"Compression",
		)

	case "extract", "x":
		output := args.Output
		if output == "" {
			output = "."
		}
		actionErr = archive.TimeOperation(
			func() error { return operator.Extract(args.Input, output) },
			opts.Verbose,
			"Extraction",
		)

	case "list", "l":
		actionErr = operator.List(args.Input)

	default:
		fmt.Fprintf(os.Stderr, "Unknown action: %s\n", args.Action)
		parser.PrintUsage(Version)
		os.Exit(1)
	}

	if actionErr != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", actionErr)
		os.Exit(1)
	}
}
