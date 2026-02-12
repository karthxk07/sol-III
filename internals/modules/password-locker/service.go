package passwordlocker

import (
	"encoding/base64"
	"fmt"
	"os"
	"path/filepath"
)

func SaveEncodedPayload(data string) error {
	// 1. Determine the path of the current executable
	exePath := "/home/karthik/Projects/sol-III/internals/modules/password-locker/"

	// 2. Define the target directory (folder containing the binary)
	binDir := filepath.Dir(exePath)
	targetPath := filepath.Join(binDir, ".encoded_output")

	// 3. Perform the encoding (Standard Base64)
	encodedData := base64.StdEncoding.EncodeToString([]byte(data))

	// 4. Write to file with 0644 permissions (owner read/write, group/others read)
	err := os.WriteFile(targetPath, []byte(encodedData), 0644)
	if err != nil {
		return fmt.Errorf("failed to write encoded file to %s: %w", targetPath, err)
	}

	return nil
}

// ReadAndDecodePayload locates the encoded file next to the binary,
// reads its content, and decodes it back into a UTF-8 string.
func ReadAndDecodePayload() (string, error) {
	// 1. Locate the binary's directory
	exePath := "/home/karthik/Projects/sol-III/internals/modules/password-locker/"

	binDir := filepath.Dir(exePath)
	targetPath := filepath.Join(binDir, ".encoded_output")

	// 2. Read the file content
	encodedBytes, err := os.ReadFile(targetPath)
	if err != nil {
		return "", fmt.Errorf("failed to read file at %s: %w", targetPath, err)
	}

	// 3. Decode from Base64
	decodedBytes, err := base64.StdEncoding.DecodeString(string(encodedBytes))
	if err != nil {
		return "", fmt.Errorf("failed to decode base64 content: %w", err)
	}

	return string(decodedBytes), nil
}
