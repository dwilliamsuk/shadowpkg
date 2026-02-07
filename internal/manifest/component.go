package manifest

import (
	"bytes"
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"io"
	"net/http"
)

type Component struct {
	Name                 string   `json:"name"`
	Version              string   `json:"version"`
	Guid                 string   `json:"guid"`
	UninstallTheseBefore []string `json:"uninstall_these_before,omitempty"`
	InstallTheseBefore   []string `json:"install_these_before,omitempty"`
	RequiredSize         int      `json:"required_size"`
}

func (c *Component) Download() ([]byte, error) {
	// Download a component

	pkgUrl := fmt.Sprintf("%s/%s/%s/%s/%s.msi", baseUrl, packagesPath, c.Name, c.Version, c.Name)
	pkgMD5Url := fmt.Sprintf("%s/%s/%s/%s/%s.md5", baseUrl, packagesPath, c.Name, c.Version, c.Name)

	packageFile := func() ([]byte, []byte, error) {
		// Download package file and calculate MD5 sum

		// Download package file
		resp, err := http.Get(pkgUrl)
		if err != nil {
			return nil, nil, err
		}
		defer resp.Body.Close()

		// Create TeeReader to read bytes and write MD5 hash simultaneously
		m := md5.New()
		tr := io.TeeReader(resp.Body, m)

		// Read all bytes to b
		b, err := io.ReadAll(tr)
		if err != nil {
			return nil, nil, err
		}

		// Calculate MD5 sum
		md5Sum := m.Sum(nil)

		return b, md5Sum[:], nil
	}

	packageFileMD5 := func() ([]byte, error) {
		// Get package MD5 hash

		// Download package MD5 file
		resp, err := http.Get(pkgMD5Url)
		if err != nil {
			return nil, err
		}
		defer resp.Body.Close()

		// Read all bytes to b
		b, err := io.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}

		// Trim whitespace bytes
		b = bytes.TrimSpace(b)

		// Decode hex to bytes
		b, err = hex.DecodeString(string(b))
		if err != nil {
			return nil, err
		}

		return b, nil
	}

	packageBytes, calcMD5, err := packageFile()
	if err != nil {
		return nil, err
	}

	packageMD5, err := packageFileMD5()
	if err != nil {
		return nil, err
	}

	if !bytes.Equal(calcMD5, packageMD5) {
		return nil, fmt.Errorf("MD5 hash mismatch! Expected: [%s], Actual: [%s]", hex.EncodeToString(packageMD5), hex.EncodeToString(calcMD5))
	}

	return packageBytes, nil
}
