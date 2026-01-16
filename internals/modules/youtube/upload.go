package youtube

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"

	"github.com/karthxk07/sol-III/internals/modules/google"
)

// write a function to upload the video

func UploadYouTubeVideo(snippet Snippet, file multipart.File, size int64) error {

	var token google.Token
	if err := google.Prepare(&token); err != nil {
		return errors.Join(errors.New("error while preparing oauth"), err)
	}

	location, err := createYouTubeResumableSession(&token, snippet, size)
	if err != nil {
		return err
	}

	switch err := executeResumablePUT(location, file, size); err {

	case errUploadSuccess:
		return nil

	case errUploadResumable:
		return errUploadResumable // caller should retry with backoff

	case errUploadExpired:
		return errUploadExpired // caller should re-initiate session

	default:
		return err // permanent failure
	}
}

func createYouTubeResumableSession(token *google.Token, snippet Snippet, videoSize int64) (string, error) {
	url := "https://www.googleapis.com/upload/youtube/v3/videos?uploadType=resumable&part=snippet,status,contentDetails"

	body := map[string]any{
		"snippet": snippet,
		"status": map[string]any{
			"privacyStatus": "public",
			"embeddable":    true,
			"license":       "youtube",
		},
	}

	b, err := json.Marshal(body)
	if err != nil {
		return "", err
	}

	req, err := http.NewRequest("POST", url, bytes.NewReader(b))
	if err != nil {
		return "", err
	}

	req.Header.Set("Authorization", "Bearer "+token.AccessToken)
	req.Header.Set("Content-Type", "application/json; charset=UTF-8")
	req.Header.Set("X-Upload-Content-Type", "video/*")
	req.Header.Set("X-Upload-Content-Length", strconv.FormatInt(videoSize, 10))

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		d, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("youtube init failed: %s", d)
	}

	location := resp.Header.Get("Location")
	if location == "" {
		return "", fmt.Errorf("no resumable location returned")
	}

	return location, nil
}

var (
	errUploadSuccess       = errors.New("upload successful")
	errUploadResumable     = errors.New("upload interrupted, resumable")
	errUploadExpired       = errors.New("resumable session expired")
	errUploadPermanentFail = errors.New("upload failed permanently")
)

func executeResumablePUT(location string, file multipart.File, size int64) error {

	req, err := http.NewRequest("PUT", location, file)
	if err != nil {
		return err
	}

	req.Header.Set("Content-Length", strconv.FormatInt(size, 10))
	req.Header.Set("Content-Type", "video/*")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		// connection lost → resumable
		return errUploadResumable
	}
	defer resp.Body.Close()

	switch resp.StatusCode {

	case http.StatusCreated, 200: // 201
		return errUploadSuccess

	case 500, 502, 503, 504:
		return errUploadResumable

	case http.StatusNotFound: // 404 expired session
		return errUploadExpired

	default:
		b, _ := io.ReadAll(resp.Body)
		return errors.Join(errUploadPermanentFail, errors.New(string(b)))
	}
}

func ConvertToShort(input multipart.File) (multipart.File, int64, error) {
	// Get the project root directory (assuming this file is at ./internals/modules/youtube/upload.go)
	projectRoot, err := filepath.Abs("./")
	if err != nil {
		return nil, 0, fmt.Errorf("failed to get project root: %w", err)
	}

	videoProcessorDir := filepath.Join(projectRoot, "internals", "modules", "video_processor")
	inputVideoPath := filepath.Join(videoProcessorDir, "input_video.mp4")
	outputVideoPath := filepath.Join(videoProcessorDir, "output_video.mp4")
	pythonScriptPath := filepath.Join(videoProcessorDir, "video_processor.py")

	// Save the uploaded file as input_video.mp4
	inputFile, err := os.Create(inputVideoPath)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to create input file: %w", err)
	}

	_, err = io.Copy(inputFile, input)
	inputFile.Close()
	if err != nil {
		os.Remove(inputVideoPath)
		return nil, 0, fmt.Errorf("failed to copy input video: %w", err)
	}

	// Clean up output file if it exists from previous run
	os.Remove(outputVideoPath)

	// Run the Python script
	pythonPath := filepath.Join(videoProcessorDir, ".venv", "bin", "python")
	cmd := exec.Command(pythonPath, pythonScriptPath)
	cmd.Dir = videoProcessorDir // Set working directory to video_processor

	// Capture stdout and stderr for debugging
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		// Clean up input file
		os.Remove(inputVideoPath)
		return nil, 0, fmt.Errorf("video processing failed: %w, stderr: %s", err, stderr.String())
	}

	// Open the processed output video
	outputFile, err := os.Open(outputVideoPath)
	if err != nil {
		os.Remove(inputVideoPath)
		return nil, 0, fmt.Errorf("failed to open output video: %w", err)
	}

	// Get file size
	stat, err := outputFile.Stat()
	if err != nil {
		outputFile.Close()
		os.Remove(inputVideoPath)
		return nil, 0, fmt.Errorf("failed to stat output video: %w", err)
	}

	// Clean up input file (keep output file open for return)
	os.Remove(inputVideoPath)

	return outputFile, stat.Size(), nil
}
