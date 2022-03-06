package common

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path"
	"path/filepath"
)

// DownloadFile downloads a file to the target folder. If
// a file with same name exists. The file can be overwritten
// by setting the `overwrite` parameter.
//
// Returns the file path if the download was successful,
// an error if the file already exists, or the download
// failed.
func DownloadFile(fileURL string, targetFolder string, overwrite bool) (string, error) {
	// Determine filename
	u, err := url.Parse(fileURL)
	if err != nil {
		return "", err
	}

	fn := path.Base(u.Path)
	fp := filepath.Join(targetFolder, fn)

	// Ensure file does not exist
	if !overwrite {
		if _, err := os.Stat(fp); !os.IsNotExist(err) {
			return "", fmt.Errorf("file '%s' exists", fp)
		}
	}

	// Start file download
	resp, err := http.Get(fileURL)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	// Temporary file name.
	// Partially downloaded file gets saved under
	// a temporary file name, then moved to the final
	// file name when done.
	tfn := "." + fn + ".tmp"
	tfp := filepath.Join(targetFolder, tfn)

	// Close file before rename, because Windows locks
	// the file handle.
	err = func() error {
		// Create new file
		file, err := os.Create(tfp)
		if err != nil {
			return err
		}
		defer file.Close()

		// Stream download into file
		_, err = io.Copy(file, resp.Body)
		if err != nil {
			return err
		}

		return nil
	}()

	if err != nil {
		return "", err
	}

	// Move temporary file into final
	// file location.
	err = os.Rename(tfp, fp)
	if err != nil {
		return "", err
	}

	return "", nil
}
