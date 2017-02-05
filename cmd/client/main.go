package main

import (
	"fmt"
	"os"
)

const (
	SCRYPT_N  = 2 << 20
	SCRYPT_R  = 8
	SCRYPT_P  = 1
	SALT_SIZE = 16
)

/* Initally this was implemented in python,
 * but relied on an unmaintained gnupg wrapper.
 *
 * The idea is to make canary generation easy and allow for easy automation.
 * Furthermore the program attempts to guard against common pitfall and metadata leakage
 */

/*

func CheckErr(err error) {
	if err != nil {
		panic(err)
	}
}

func YesNo(prompt string) bool {
	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Print(prompt)
		text, _ := reader.ReadString('\n')
		if text == "y\n" {
			return true
		}
		if text == "n\n" {
			return false
		}
	}
}

func EnterToContinue() {
	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Print("Press enter to continue ")
		text, _ := reader.ReadString('\n')
		if text == "\n" {
			return
		}
	}
}

func OperationCheck() {
	// Load public key
	f, err := os.Open(FlagPublicKeyPath)
	if err != nil {
		fmt.Println("Failed to open public key:", FlagPublicKeyPath, err)
		os.Exit(EXIT_INVALID_ADDRESS)
		return
	}
	block, err := armor.Decode(f)
	if err != nil {
		log.Println("Failed to load public key:", FlagPublicKeyPath, err)
		os.Exit(EXIT_INVALID_PUBLIC_KEY)
		return
	} else if block.Type != openpgp.PrivateKeyType {
		log.Println("Not a PGP public key:", FlagPublicKeyPath, err)
		os.Exit(EXIT_INVALID_PUBLIC_KEY)
		return
	}

	// Load entity
	entity, err := openpgp.ReadEntity(packet.NewReader(block.Body))
	if err != nil {
		log.Println("Not a PGP public key:", FlagPublicKeyPath, err)
		os.Exit(EXIT_INVALID_PUBLIC_KEY)
		return
	}

	// Load existing canary from disk
	files, err := ioutil.ReadDir(FlagCanaryDirectory)
	if err != nil {
		log.Println("Unable to list canaries in directory:", FlagCanaryDirectory)
		os.Exit(EXIT_NO_SUCH_FILE)
	}
	var canary *Canary
	if len(files) > 0 {
		log.Println("Loading prior canary")
		for i := len(files) - 1; i >= 0; i-- {
			// Check name format
			if files[i].IsDir() {
				continue
			}
			match, err := regexp.MatchString("canary-\\d{4}-\\d{2}-\\d{2}-.*", files[i].Name())
			if err != nil {
				panic(err)
			}
			if !match {
				continue
			}

			// Read contents
			canaryFile := path.Join(FlagCanaryDirectory, files[len(files)-1].Name())
			proof, err := ioutil.ReadFile(canaryFile)
			if err != nil {
				log.Println("Failed to read canary", err)
				os.Exit(EXIT_FILE_READ_ERROR)
			}
			log.Println("DEBUG:", string(proof))

			// Parse file contents
			canary, err = CanaryParse(entity, proof)
			if err != nil {
				log.Println("Failed to parse local canary", err)
				os.Exit(EXIT_INVALID_LOCAL_CANARY)
			}
			break
		}
	}

	// Download new canary
	canary, hash, err := CanaryFetchLatest(entity, FlagAddress)
	log.Println("Latest canary hash:", canary, hash, err)
}

*/

func main() {
	// parse flags
	flags := parseFlags()

	// handle different operations
	switch flags.Operation {
	case "create":
		operationCreate(flags)
	case "add":
		operationAdd(flags)
	case "push":
		operationPush(flags)
	case "pull":
		operationPull(flags)
	case "":
		printHelp()
	default:
		fmt.Fprintf(os.Stderr, "Invalid operation: %s\n", flags.Operation)
		os.Exit(EXIT_INVALID_OPERATION)
	}
}
