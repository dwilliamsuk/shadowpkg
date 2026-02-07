package main

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"sync"

	"github.com/dwilliamsuk/shadowpkg/internal/manifest"
)

var (
	envPtr    = flag.String("environment", "prod", "The environment to download packages from")
	outputPtr = flag.String("output", ".", "The output directory to download packages to")
)

func downloadPackage(outputFolder string, c *manifest.Component) error {
	// Download and save an installer package

	// Download installer package bytes to b
	b, err := c.Download()
	if err != nil {
		return err
	}

	// Write installer package bytes to file
	outputFp := filepath.Join(outputFolder, c.Name, c.Version, fmt.Sprintf("%s.msi", c.Name))

	err = os.WriteFile(outputFp, b, 0644)
	if err != nil {
		return err
	}

	// Calculate SHA256 sum of bytes, and store with file
	sha256Sum := sha256.Sum256(b)
	sha256SumString := fmt.Sprintf("%s  %s", hex.EncodeToString(sha256Sum[:]), filepath.Base(outputFp))

	err = os.WriteFile(fmt.Sprintf("%s.sha256", outputFp), []byte(sha256SumString), 0644)
	if err != nil {
		return err
	}

	return nil
}

func saveManifest(outputFp string, m *manifest.Manifest) error {
	// Save manifest struct as JSON

	// Marshal struct to JSON
	manifestJson, err := json.MarshalIndent(m, "", "  ")
	if err != nil {
		return err
	}

	// Save bytes to fp
	err = os.WriteFile(outputFp, manifestJson, 0644)
	if err != nil {
		return err
	}

	return nil
}

func init() {
	flag.Parse()
}

func main() {
	// Download all packages within env defined by args

	// Get arg values
	env := *envPtr
	outputFolder := *outputPtr

	// Get manifest struct for specified env
	manifest, err := manifest.New(env)
	if err != nil {
		panic(err)
	}

	// Pre-create folders for files
	for _, c := range manifest.Components {
		outputFolder := filepath.Join(outputFolder, c.Name, c.Version)
		err = os.MkdirAll(outputFolder, 0755)
		if err != nil {
			panic(err)
		}
	}

	// Download all packages within manifest concurrently
	var wg sync.WaitGroup
	wg.Add(len(manifest.Components))

	fmt.Printf("Downloading [%d] packages...\n", len(manifest.Components))

	for _, c := range manifest.Components {
		go func() {
			defer wg.Done()

			err := downloadPackage(outputFolder, &c)
			if err != nil {
				panic(err)
			}
		}()
	}
	wg.Wait()

	// Save manifest as json with packages
	outputFp := filepath.Join(outputFolder, "manifest.json")
	err = saveManifest(outputFp, manifest)
	if err != nil {
		panic(err)
	}

	absPath, _ := filepath.Abs(outputFolder)

	if absPath != "" {
		fmt.Printf("Done! Saved packages to [%s]\n", absPath)
	} else {
		fmt.Printf("Done! Saved packages to [%s]\n", filepath.Clean(outputFolder))
	}
}
