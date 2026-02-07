package manifest

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

const (
	baseUrl      = "https://oblivion.shadow.tech"
	envPath      = "environments"
	manifestPath = "manifests"
	packagesPath = "packages"
)

type Manifest struct {
	Metadata   Metadata    `json:"metadata"`
	Components []Component `json:"vm_components"`
}

type Metadata struct {
	Environment           string `json:"environment"`
	SaveAsCurrentManifest bool   `json:"save_as_current_manifest"`
	ProcessUsingDiff      bool   `json:"process_using_diff"`
}

func New(env string) (*Manifest, error) {
	// Create a new manifest struct for env from Shadow CDN manifest

	fileId := func() (string, error) {
		// Get manifest artifact id

		// Get id of manifest from environments URL
		resp, err := http.Get(fmt.Sprintf("%s/%s/%s", baseUrl, envPath, env))
		if err != nil {
			return "", err
		}
		defer resp.Body.Close()

		b, err := io.ReadAll(resp.Body)
		if err != nil {
			return "", err
		}

		// Remove whitespace / newlines from web server
		b = bytes.TrimSpace(b)

		return string(b), nil
	}

	manifestBytes := func() ([]byte, error) {
		// Get manifest JSON bytes

		// Get manifest file ID
		manifestFileId, err := fileId()
		if err != nil {
			return nil, fmt.Errorf("Unable to retrieve manifest artifact id! Err: %s", err)
		}

		// Download manifest JSON
		resp, err := http.Get(fmt.Sprintf("%s/%s/%s", baseUrl, manifestPath, manifestFileId))
		if err != nil {
			return nil, err
		}
		defer resp.Body.Close()

		// Read all bytes to b
		b, err := io.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}

		return b, nil
	}

	mb, err := manifestBytes()
	if err != nil {
		return nil, fmt.Errorf("Unable to retrieve manifest JSON! Err: %s", err)
	}

	// Unmarshal JSON
	m := &Manifest{}
	err = json.Unmarshal(mb, &m)

	return m, nil
}
