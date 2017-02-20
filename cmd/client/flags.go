package main

import (
	"flag"
	"fmt"
	"github.com/fatih/color"
	"github.com/rot256/fugl"
	"os"
	"time"
)

const tagFlags = "flag"

type Flags struct {
	Proof       string        // path to proof
	PublicKey   string        // path to pgp public key
	PrivateKey  string        // path to pgp private key
	Author      string        // creator of canary
	Description string        // file containing canary description
	Expire      time.Duration // expiration delta
	Proxy       string        // proxy to use (e.g 127.0.0.1:9050)
	Address     string        // address of submission point
	Manifest    string        // manifest, for creating canaries
	Operation   string        // operation to apply
	Debug       bool          // used during development
	Json        bool          // enable json output
	Help        bool          // print help
}

const (
	FlagNamePublicKey  = "public-key"
	FlagNamePrivateKey = "private-key"
	FlagNameProxy      = "proxy"
	FlagNameAddress    = "address"
	FlagNameOperation  = "operation"
	FlagNameDebug      = "debug"
	FlagNameHelp       = "help"
	FlagNameJson       = "json"
	FlagNameManifest   = "manifest"
	FlagNameProof      = "proof"
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

func exitError(code int, format string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, format, args...)
	fmt.Fprintf(os.Stderr, "\n")
	os.Exit(code)
}

func (opt FlagOpt) Check() {
	okay := true
	for name := range opt.required {
		okay = okay && opt.required[name]
	}
	if !okay {
		for name := range opt.required {
			if !opt.required[name] {
				line := fmt.Sprintf("required argument: '%s' (missing)\n", name)
				color.Red(line)
			} else {
				line := fmt.Sprintf("required argument: '%s' (supplied)\n", name)
				color.Green(line)
			}
		}
		for name := range opt.optional {
			if !opt.optional[name] {
				line := fmt.Sprintf("optional argument: '%s' (missing)\n", name)
				color.Blue(line)
			} else {
				line := fmt.Sprintf("optional argument: '%s' (supplied)\n", name)
				color.Green(line)
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
	flag.StringVar(&flags.Manifest, FlagNameManifest, "./manifest.toml", "canary manifest, for creating new canaries")
	flag.StringVar(&flags.Proof, FlagNameProof, "./temp"+fugl.ProofFileExtension, "path to proof")
	flag.StringVar(&flags.PrivateKey, FlagNamePrivateKey, "", "path to a PGP private key")
	flag.StringVar(&flags.PublicKey, FlagNamePublicKey, "", "path to a PGP public key")
	flag.StringVar(&flags.Proxy, FlagNameProxy, "", "socks5 proxy")
	flag.StringVar(&flags.Address, FlagNameAddress, "", "address of canary server")
	flag.StringVar(&flags.Operation, FlagNameOperation, "", "operation, supported: pull, push, verify")
	flag.BoolVar(&flags.Debug, FlagNameDebug, false, "enable debugging")
	flag.BoolVar(&flags.Help, FlagNameHelp, false, "print this help page")
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

    push   : uploads a new canary to a server
    pull   : downloads the latest canary from the remote
    verify : verifies a locally stored canary
    create : creates a new canary locally

  Using --operation=[action]
  You may specify any one of these to see what arguments they require.

2. List of flags:
`
	fmt.Fprintf(os.Stderr, msg)
	flag.PrintDefaults()
}
