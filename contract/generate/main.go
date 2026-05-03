// This program compiles contract/tredd.sol using a pinned, automatically downloaded version of solc (via go-solc).
//
// Run it from the module root with:
//
//	go run ./contract/generate
//
// or via `task contract_go`.
package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/bobg/errors"
	solc "github.com/lmittmann/go-solc"
)

const pinnedVersion = solc.Version0_7_2

func main() {
	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %s\n", err)
		os.Exit(1)
	}
}

func run() error {
	modRoot, err := findModRoot()
	if err != nil {
		return errors.Wrap(err, "finding module root")
	}

	solcBin := solcBinaryPath(modRoot)

	// Trigger the download/verification of the pinned solc binary.
	// go-solc stores it at <modRoot>/.solc/bin/solc_v<version>.
	// We call Compile on the actual source; even if compilation produces
	// warnings we don't care, we only need the side-effect of the download.
	// On error we check whether it's a download failure (fatal) or a
	// compilation issue (we'll let the direct invocation report it properly).
	compiler := solc.New(pinnedVersion)
	if _, err := compiler.Compile("contract", "Tredd"); err != nil {
		// If the binary doesn't exist after this, it was a download failure.
		if _, statErr := os.Stat(solcBin); statErr != nil {
			return errors.Wrap(err, "Downloading solc")
		}
		// Otherwise the download succeeded; the error was a compilation issue
		// that the direct invocation below will surface more clearly.
	}

	// Compile tredd.sol → ERC20.abi, Tredd.abi, Tredd.bin
	cmd := exec.Command(solcBin, "--abi", "--bin", "--overwrite", "-o", ".", "tredd.sol")
	cmd.Dir = "contract"
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

// solcBinaryPath returns the path where go-solc caches the binary.
// This mirrors the naming convention in go-solc's download.go:
//
//	<modRoot>/.solc/bin/solc_v<version>
func solcBinaryPath(modRoot string) string {
	name := fmt.Sprintf("solc_v%s", pinnedVersion)
	if runtime.GOOS == "windows" {
		name += ".exe"
	}
	return filepath.Join(modRoot, ".solc", "bin", name)
}

// findModRoot returns the directory containing go.mod.
func findModRoot() (string, error) {
	out, err := exec.Command("go", "env", "GOMOD").Output()
	if err != nil {
		return "", errors.Wrap(err, "running go env GOMOD")
	}
	path := strings.TrimSpace(string(out))
	if path == "" || path == os.DevNull {
		return "", fmt.Errorf("no go.mod found; run from within the module")
	}
	return filepath.Dir(path), nil
}
