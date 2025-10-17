// Package cli provides command-line interface functionality
package cli

import (
	"flag"
	"fmt"
	"runtime"
	"strings"

	"github.com/cubetiqlabs/gar/internal/models"
)

// Parser handles command-line argument parsing
type Parser struct {
	flagSet *flag.FlagSet
}

// NewParser creates a new CLI parser
func NewParser() *Parser {
	return &Parser{
		flagSet: flag.NewFlagSet("gar", flag.ContinueOnError),
	}
}

// Parse parses command-line arguments and returns CLIArgs
func (p *Parser) Parse(args []string) (*models.CLIArgs, error) {
	// Pre-process arguments to expand combined flags like -cvf to -c -v -f
	processedArgs := p.processUnixStyleFlags(args)

	// Command line flags
	var (
		// Long-form flags (backward compatibility)
		action      = p.flagSet.String("action", "", "Action: compress, extract, list")
		input       = p.flagSet.String("input", "", "Input file or directory")
		output      = p.flagSet.String("output", "", "Output file or directory")
		format      = p.flagSet.String("format", "zip", "Archive format: zip, tar.gz")
		password    = p.flagSet.String("password", "", "Password for encryption")
		compression = p.flagSet.String("compression", "normal", "Compression level: fastest, normal, best")
		workers     = p.flagSet.Int("workers", runtime.NumCPU(), "Number of worker threads")
		verbose     = p.flagSet.Bool("verbose", false, "Verbose output")
		version     = p.flagSet.Bool("version", false, "Show version")
		help        = p.flagSet.Bool("help", false, "Show help message")
		h           = p.flagSet.Bool("h", false, "Show help message (short)")

		// Unix-style single char flags
		c = p.flagSet.Bool("c", false, "(Unix-style) Compress")
		x = p.flagSet.Bool("x", false, "(Unix-style) Extract")
		t = p.flagSet.Bool("t", false, "(Unix-style) Test/List archive")
		v = p.flagSet.Bool("v", false, "(Unix-style) Verbose")
		_ = p.flagSet.Bool("f", false, "(Unix-style) File (archive path)")
		z = p.flagSet.Bool("z", false, "(Unix-style) Force gzip/TAR.GZ")
		j = p.flagSet.Bool("j", false, "(Unix-style) Force bzip2")
		Z = p.flagSet.Bool("Z", false, "(Unix-style) Force 7zip")
	)

	// Parse the pre-processed flags
	if err := p.flagSet.Parse(processedArgs); err != nil {
		return nil, err
	}

	// Build result
	result := &models.CLIArgs{
		Workers: *workers,
	}

	// Handle help flags
	if *help || *h {
		result.Help = true
		return result, nil
	}

	if *version {
		result.Version = true
		return result, nil
	}

	// Get remaining arguments
	posArgs := p.flagSet.Args()

	// Build options from Unix-style flags if they were used
	unixVerbose := *v
	unixFormat := *format
	if *z {
		unixFormat = "tar.gz"
	} else if *j {
		unixFormat = "bzip2"
	} else if *Z {
		unixFormat = "7zip"
	}

	var unixAction string
	if *c {
		unixAction = "compress"
	} else if *x {
		unixAction = "extract"
	} else if *t {
		unixAction = "list"
	}

	// Parse positional arguments for Unix-style
	var unixInput, unixOutput string

	if *c && len(posArgs) >= 1 {
		// Compress: first arg is output archive, second is input path
		unixOutput = posArgs[0]
		if len(posArgs) > 1 {
			unixInput = posArgs[1]
		}
	} else if (*x || *t) && len(posArgs) >= 1 {
		// Extract or List: first arg is input archive, second is output path
		unixInput = posArgs[0]
		if len(posArgs) > 1 {
			unixOutput = posArgs[1]
		}
	} else if len(posArgs) > 0 {
		// No Unix flags, treat as traditional
		unixInput = posArgs[0]
		if len(posArgs) > 1 {
			unixOutput = posArgs[1]
		}
	}

	// If we have action from flags, use the parsed values; otherwise use traditional flags
	if unixAction != "" {
		result.Action = unixAction
		result.Input = unixInput
		result.Output = unixOutput
	} else {
		result.Action = *action
		result.Input = *input
		result.Output = *output
	}

	if unixVerbose {
		result.Verbose = unixVerbose
	} else {
		result.Verbose = *verbose
	}

	result.Format = unixFormat
	result.Password = *password
	result.Compression = *compression

	return result, nil
}

// ProcessUnixStyleFlags converts Unix-style combined flags (like -cvf) into separate flags
// This allows us to support tar-like commands like: gar -cvf archive.zip folder
func (p *Parser) processUnixStyleFlags(args []string) []string {
	var result []string

	for i := 0; i < len(args); i++ {
		arg := args[i]

		// Check if this is a combined short flag (starts with -, has multiple chars, not --)
		if strings.HasPrefix(arg, "-") && !strings.HasPrefix(arg, "--") && len(arg) > 2 {
			// This could be a combined flag like -cvf or -xvf
			flags := arg[1:] // Remove leading dash

			// Check if it contains only valid flag characters
			allValidFlags := true
			for _, ch := range flags {
				if !strings.ContainsRune("cvxtfjzZ", ch) {
					allValidFlags = false
					break
				}
			}

			if allValidFlags && strings.ContainsAny(flags, "cxt") {
				// It's a Unix-style combined flag
				// Expand it: -cvf becomes -c -v -f
				for _, ch := range flags {
					result = append(result, "-"+string(ch))
				}
				continue
			}
		}

		result = append(result, arg)
	}

	return result
}

// PrintUsage prints the usage information
func (p *Parser) PrintUsage(version string) {
	fmt.Println("GoArchive (gar) - High-Performance Cross-Platform Archive Manager")
	fmt.Println()
	fmt.Printf("Version: %s\n", version)
	fmt.Println()
	fmt.Println("Usage (Unix-style):")
	fmt.Println("  gar -cvf archive.zip folder              Compress folder with verbose")
	fmt.Println("  gar -xvf archive.zip [output_path]       Extract archive with verbose")
	fmt.Println("  gar -tvf archive.zip                     List archive contents")
	fmt.Println()
	fmt.Println("Usage (Long-form flags):")
	fmt.Println("  gar -action=compress -input=<path> -output=<file> [options]")
	fmt.Println("  gar -action=extract -input=<file> -output=<path> [options]")
	fmt.Println("  gar -action=list -input=<file> [options]")
	fmt.Println()
	fmt.Println("Unix-style Options:")
	fmt.Println("  c              Compress")
	fmt.Println("  x              Extract")
	fmt.Println("  t              Test/List archive contents")
	fmt.Println("  v              Verbose output")
	fmt.Println("  f              File (archive path) - must follow other options")
	fmt.Println("  z              Force gzip compression (TAR.GZ format)")
	fmt.Println("  j              Force bzip2 compression")
	fmt.Println("  Z              Force 7zip compression")
	fmt.Println()
	fmt.Println("Long-form Options:")
	p.flagSet.PrintDefaults()
	fmt.Println()
	fmt.Println("Examples:")
	fmt.Println("  gar -cvf my-archive.zip my-folder")
	fmt.Println("  gar -xvf my-archive.zip")
	fmt.Println("  gar -xvf my-archive.zip /tmp/extract")
	fmt.Println("  gar -tvf my-archive.zip")
	fmt.Println("  gar -action=compress -input=folder -output=archive.zip -verbose")
}
