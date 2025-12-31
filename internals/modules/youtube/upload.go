package youtube

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/karthxk07/sol-III/internals/modules/google"
)

// write a function to upload the video
func UploadVideo(data io.Reader) error {

	var token google.Token
	err := google.Prepare(&token)
	if err != nil {
		return errors.Join(errors.New("error while preparing oauth"), err)
	}

	initReq, _ := http.NewRequest("POST",
		"https://www.googleapis.com/upload/youtube/v3/videos?uploadType=resumable&part=snippet,status",
		nil,
	)

	initReq.Header.Set("Authorization", "Bearer "+token.AccessToken)
	initReq.Header.Set("Content-Type", "application/json")

	initReq.Body = io.NopCloser(strings.NewReader(`{
        "snippet": {"title": "Auto upload", "categoryId": "22"},
        "status": {"privacyStatus": "private"}
    }`))

	resp, err := http.DefaultClient.Do(initReq)
	if err != nil {
		return err
	}
	uploadURL := resp.Header.Get("Location")
	resp.Body.Close()

	// Now stream the actual bytes
	uploadReq, _ := http.NewRequest("PUT", uploadURL, data)
	uploadReq.Header.Set("Content-Type", "application/octet-stream")

	uploadResp, err := http.DefaultClient.Do(uploadReq)
	if err != nil {
		return err
	}
	defer uploadResp.Body.Close()

	if uploadResp.StatusCode >= 300 {
		b, _ := io.ReadAll(uploadResp.Body)
		return fmt.Errorf("YT: %s", string(b))
	}
	return nil
}
