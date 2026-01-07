# YouTube Uploader API

A local Gin-based API that uploads videos to YouTube using resumable uploads.  
Supports automatic YouTube Shorts conversion (9:16 vertical crop + re-encode).

---

## OAuth Setup

You must authorize once before using the upload endpoint.

### Get OAuth Consent URL
GET /google/getOAuthurl

yaml
Copy code

This returns the Google OAuth consent URL.  
Open it in your browser, grant access, and the token will be stored for future uploads.

---

## Upload Endpoint

POST /youtube/upload

yaml
Copy code

Uploads a video to YouTube.  
Optionally converts it into a YouTube Short when `isReel=true`.

---

## Multipart Form Fields

| Field | Type | Required | Description |
|------|-----|----------|-------------|
| title | string | Yes | Video title |
| description | string | Yes | Video description |
| categoryId | number | Yes | YouTube category ID |
| tags | string | No | Comma-separated tags |
| isReel | boolean | No | `true` → convert to Shorts |
| video | file | Yes | MP4 video file |

---

## Example Request

```bash
curl -X POST http://localhost:8080/youtube/upload \
  -F "title=My Reel Test" \
  -F "description=Auto converted to Shorts" \
  -F "categoryId=22" \
  -F "tags=reel,shorts,test" \
  -F "isReel=true" \
  -F "video=@/home/<username>/Downloads/video.mp4"
Shorts Conversion
When isReel=true:

Video is center-cropped to 9:16

Re-encoded to H.264 MP4

Automatically detected by YouTube as a Short (≤60s)

Upload Flow
pgsql
Copy code
multipart request
      ↓
(optional Shorts conversion)
      ↓
create resumable session
      ↓
fault-tolerant resumable upload
      ↓
video published
Error Handling
Case	Action
201 Created	Upload successful
500 / 502 / 503 / 504	Upload can be resumed
404	Session expired — restart upload
Other 4xx / 5xx	Permanent failure

Requirements
Go 1.20+

FFmpeg installed

Google OAuth credentials

YouTube Data API v3 enabled

API Summary
Endpoint	Method	Description
/google/getOAuthurl	GET	Get OAuth consent URL
/youtube/upload	POST	Upload video / Shorts


