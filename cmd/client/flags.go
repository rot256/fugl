package main

import (
	"flag"
	"fmt"
	"github.com/fatih/color"
	"os"
	"time"
)

const tagFlags = "flag"

type Flags struct {
	Store      string        // directory for storing canaries (when validating)
	PublicKey  string        // path to pgp public key
	PrivateKey string        // path to pgp private key
	Message    string        // path to canary message (notification)
	Input      string        // input file of operation
	Output     string        // output file of operation
	Expire     time.Duration // expiration delta
	Proxy      string        // proxy to use (e.g 127.0.0.1:9050)
	Address    string        // address of submission point
	Operation  string        // operation to apply
	Debug      bool          // used during development
	Json       bool          // enable json output
	Help       bool          // print help
}

const (
	FlagNameInput      = "input"
	FlagNameOutput     = "output"
	FlagNameStore      = "store"
	FlagNamePublicKey  = "public-key"
	FlagNamePrivateKey = "private-key"
	FlagNameMessage    = "message"
	FlagNameProxy      = "proxy"
	FlagNameAddress    = "address"
	FlagNameOperation  = "operation"
	FlagNameDebug      = "debug"
	FlagNameHelp       = "help"
	FlagNameExpire     = "expire"
	FlagNameJson       = "json"
)

func init() {
	flag.Usage = printHelp
}

type FlagOpt struct {
	required map[string]bool
	optional map[string]bool
}

func (opt *FlagOpt) Required(name string, isEnabled bool) {
	if opt.required == nil {
		opt.required = make(map[string]bool)
	}
	opt.required[name] = isEnabled
}

func (opt *FlagOpt) Optional(name string, isEnabled bool) {
	if opt.optional == nil {
		opt.optional = make(map[string]bool)
	}
	opt.optional[name] = isEnabled
}

func (opt FlagOpt) Check() {
	okay := true
	for name := range opt.required {
		okay = okay && opt.required[name]
	}
	if !okay {
		for name := range opt.required {
			line := fmt.Sprintf("required argument: '%s'\n", name)
			if !opt.required[name] {
				color.Red(line)
			}
		}
		for name := range opt.optional {
			line := fmt.Sprintf("optional argument: '%s'\n", name)
			if !opt.optional[name] {
				color.Blue(line)
			}
		}
		os.Exit(EXIT_INVALID_ARGUMENTS)
	}
}

func requiredArgument(name string) string {
	return fmt.Sprintf("requires argument: '%s'\n", name)
}

func optionalArgument(name string) string {
	return fmt.Sprintf("optional argument: '%s'\n", name)
}

func parseFlags() Flags {
	var flags Flags
	flag.StringVar(&flags.Input, FlagNameInput, "", "path to input file")
	flag.StringVar(&flags.Output, FlagNameOutput, "", "path to output file")
	flag.StringVar(&flags.Store, FlagNameStore, "", "path to store/state")
	flag.StringVar(&flags.PrivateKey, FlagNamePrivateKey, "", "path to a PGP private key")
	flag.StringVar(&flags.PublicKey, FlagNamePublicKey, "", "path to a PGP public key")
	flag.StringVar(&flags.Message, FlagNameMessage, "", "path to a custom notification")
	flag.StringVar(&flags.Proxy, FlagNameProxy, "", "socks5 proxy")
	flag.StringVar(&flags.Address, FlagNameAddress, "", "address of submission point")
	flag.StringVar(&flags.Operation, FlagNameOperation, "", "operation, supported: pull, push, verify")
	flag.BoolVar(&flags.Json, FlagNameJson, false, "enable json output (stdout only)")
	flag.BoolVar(&flags.Debug, FlagNameDebug, false, "enable debugging")
	flag.BoolVar(&flags.Help, FlagNameHelp, false, "print help")
	flag.DurationVar(&flags.Expire, FlagNameExpire, 0, "expiration delta")
	flag.Parse()
	if flags.Debug {
		fmt.Println("flags:", flags)
	}
	if flags.Help {
		printHelp()
		os.Exit(EXIT_SUCCESS)
	}
	if flags.Json {
		fmt.Fprintf(os.Stderr, "JSON output not yet supported")
		os.Exit(EXIT_INVALID_ARGUMENTS)
	}
	return flags
}

func printHelp() {
	msg := `Help:
1. Getting started
  This is a client for the fugl canary system.
  To use this client you must specify 1 of three operations:

    push   : uploads a new canary to a remote http server
    pull   : downloads the latest canary from the remote
    verify : verifies a locally stored canary against the state of store
    create : creates a new canary locally

  You may specify any one of these to see what arguments they require.
`

	fmt.Println(msg)
}
