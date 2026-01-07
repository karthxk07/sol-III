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

	inTmp, _ := os.CreateTemp("", "yt-in-*.mp4")
	defer os.Remove(inTmp.Name())
	io.Copy(inTmp, input)
	inTmp.Close()

	outTmp, _ := os.CreateTemp("", "yt-out-*.mp4")
	outTmp.Close()

	// crop to vertical 9:16 and re-encode
	cmd := exec.Command("ffmpeg",
		"-y",
		"-i", inTmp.Name(),
		"-vf", "crop=ih*9/16:ih",
		"-c:v", "libx264",
		"-preset", "veryfast",
		"-crf", "23",
		"-movflags", "+faststart",
		outTmp.Name(),
	)

	if err := cmd.Run(); err != nil {
		return nil, 0, err
	}

	f, err := os.Open(outTmp.Name())
	if err != nil {
		return nil, 0, err
	}

	stat, _ := f.Stat()
	return f, stat.Size(), nil
}
