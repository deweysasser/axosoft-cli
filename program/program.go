package program

import (
	"fmt"
	"io"
	"os"
	"runtime"

	"github.com/alecthomas/kong"
	"github.com/mattn/go-colorable"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

var exitFunc = os.Exit

// Options is the structure of program options
type Options struct {
	Version bool `help:"Show program version"`

	Debug        bool   `group:"Info" help:"Show debugging information"`
	OutputFormat string `group:"Info" enum:"auto,jsonl,terminal" default:"auto" help:"How to show program output (auto|terminal|jsonl)"`
	Quiet        bool   `group:"Info" help:"Be less verbose than usual"`

	BaseURL  string `env:"AXOSOFT_URL" help:"Axosoft base URL (e.g. https://example.axosoft.com)"`
	ClientID string `env:"AXOSOFT_CLIENT_ID" help:"Axosoft OAuth client ID"`

	Login LoginCmd `cmd:"" help:"Authenticate with Axosoft via OAuth browser flow"`
	Serve ServeCmd `cmd:"" help:"Start MCP server via stdio transport"`

	Projects ProjectsCmd `cmd:"" help:"List all projects"`
	Releases ReleasesCmd `cmd:"" help:"List releases"`
	Items    ItemsCmd    `cmd:"" help:"List work items"`
	Item     ItemCmd     `cmd:"" help:"Get a single work item"`
}

// Parse calls the CLI parsing routines
func (program *Options) Parse(args []string) (*kong.Context, error) {
	parser, err := kong.New(program,
		kong.ShortUsageOnError(),
		kong.Description("Axosoft MCP server and CLI tool"),
	)

	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return parser.Parse(args)
}

// AfterApply runs after the options are parsed but before anything runs
func (program *Options) AfterApply() error {
	program.initLogging()
	return nil
}

func (program *Options) initLogging() {
	if program.Version {
		fmt.Println(Version)
		exitFunc(0)
	}

	switch {
	case program.Debug:
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	case program.Quiet:
		zerolog.SetGlobalLevel(zerolog.WarnLevel)
	default:
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
	}

	// Always log to stderr — stdout is reserved for MCP protocol in serve mode
	var out io.Writer = os.Stderr

	if os.Getenv("TERM") == "" && runtime.GOOS == "windows" {
		out = colorable.NewColorableStderr()
	}

	if program.OutputFormat == "terminal" ||
		(program.OutputFormat == "auto" && isTerminal(os.Stderr)) {
		log.Logger = log.Output(zerolog.ConsoleWriter{Out: out})
	} else {
		log.Logger = log.Output(out)
	}

	log.Logger.Debug().
		Str("version", Version).
		Str("program", os.Args[0]).
		Msg("Starting")
}

// isTerminal returns true if the file given points to a character device (i.e. a terminal)
func isTerminal(file *os.File) bool {
	if fileInfo, err := file.Stat(); err != nil {
		log.Err(err).Msg("Error running stat")
		return false
	} else {
		return (fileInfo.Mode() & os.ModeCharDevice) != 0
	}
}
