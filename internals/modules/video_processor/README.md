# Video Crop & Auto-Subtitle Tool

A Python script that automatically crops videos to 9:16 ratio (perfect for Instagram Reels, TikTok, YouTube Shorts) and adds auto-generated subtitles using offline speech recognition.

## Features

- 🎬 Automatic cropping to 9:16 aspect ratio
- 🗣️ Offline speech-to-text transcription (no API costs!)
- 📝 Auto-generated SRT subtitles
- 🎨 Customizable subtitle styling (font, size, color, outline)
- ⚡ Lightweight and fast processing

## Prerequisites

Before you begin, ensure you have the following installed on your system:

### 1. Python 3.7 or higher
Check your Python version:
```bash
python --version
```

### 2. FFmpeg
FFmpeg is required for video processing.

**Ubuntu/Debian:**
```bash
sudo apt update
sudo apt install ffmpeg
```

**macOS:**
```bash
brew install ffmpeg
```

**Windows:**
1. Download from [ffmpeg.org](https://ffmpeg.org/download.html)
2. Extract and add to your system PATH
3. Verify installation: `ffmpeg -version`

## Installation

### Step 1: Clone or Download the Project
```bash
cd your-project-directory
```

### Step 2: Create Virtual Environment
```bash
# Create virtual environment
python -m venv venv

# Activate virtual environment
# On Linux/macOS:
source venv/bin/activate

# On Windows:
venv\Scripts\activate
```

You should see `(venv)` in your terminal prompt after activation.

### Step 3: Install Python Dependencies
```bash
pip install -r requirements.txt
```

### Step 4: Download Vosk Speech Recognition Model

1. **Visit the Vosk Models page:**
   [https://alphacephei.com/vosk/models](https://alphacephei.com/vosk/models)

2. **Choose a model based on your needs:**
   
   **For English (Recommended):**
   - **Small model** (~40MB): `vosk-model-small-en-us-0.15`
     - Fast, lightweight, good for most use cases
     - Download: [Direct Link](https://alphacephei.com/vosk/models/vosk-model-small-en-us-0.15.zip)
   
   - **Large model** (~1.8GB): `vosk-model-en-us-0.22`
     - More accurate, slower processing
     - Download: [Direct Link](https://alphacephei.com/vosk/models/vosk-model-en-us-0.22.zip)

   **Other Languages:**
   - Spanish: `vosk-model-small-es-0.42`
   - French: `vosk-model-small-fr-0.22`
   - German: `vosk-model-small-de-0.15`
   - And many more on the website

3. **Extract the model:**
   ```bash
   # Download and extract (example for small English model)
   wget https://alphacephei.com/vosk/models/vosk-model-small-en-us-0.15.zip
   unzip vosk-model-small-en-us-0.15.zip
   
   # Rename to 'model'
   mv vosk-model-small-en-us-0.15 model
   ```

4. **Verify your directory structure:**
   ```
   your-project-directory/
   ├── script.py
   ├── requirements.txt
   ├── input_video.mp4
   └── model/                    # Vosk model folder
       ├── am/
       ├── conf/
       ├── graph/
       ├── ivector/
       └── ...
   ```

## Usage

### Step 1: Place Your Video
Put your input video in the project directory and name it `input_video.mp4` (or modify the filename in `script.py`).

### Step 2: Run the Script
```bash
python script.py
```

### Step 3: Output
The processed video will be saved as `output_video.mp4` with:
- 9:16 aspect ratio (centered crop)
- Auto-generated subtitles burned into the video
- SRT subtitle file saved as `subtitles.srt`

## Customization

### Change Font Style
Edit the `crop_and_add_subtitles()` call in `script.py`:

```python
crop_and_add_subtitles(
    input_video, 
    output_video, 
    srt_path,
    font_name="Impact",          # Font family
    font_size=32,                # Font size in pixels
    font_color="&H00FFFFFF",     # White text
    outline_color="&H00000000",  # Black outline
    outline=2,                   # Outline thickness
    shadow=1,                    # Shadow depth
    bold=1                       # 1 for bold, 0 for normal
)
```

### Common Font Styles

**TikTok/Instagram Style (Bold Impact):**
```python
font_name="Impact"
font_size=32
bold=1
```

**YouTube Style (Clean Arial):**
```python
font_name="Arial"
font_size=28
bold=0
```

**Color Options:**
- White: `&H00FFFFFF`
- Black: `&H00000000`
- Yellow: `&H0000FFFF`
- Red: `&H000000FF`
- Blue: `&H00FF0000`

### Change Input/Output Files
Modify these variables in the `main()` function:
```python
input_video = "input_video.mp4"      # Your input file
output_video = "output_video.mp4"    # Your output file
```

### Adjust Subtitle Timing
Modify the `max_words_per_subtitle` parameter in `create_srt_from_vosk()`:
```python
create_srt_from_vosk(results, srt_path, max_words_per_subtitle=8)
```
- Lower value = More frequent subtitles
- Higher value = Longer subtitles

## Troubleshooting

### "Vosk model not found" Error
Make sure the model folder is named exactly `model` and is in the same directory as your script.

### "ffmpeg not found" Error
Ensure FFmpeg is installed and added to your system PATH. Test with:
```bash
ffmpeg -version
```

### Poor Transcription Quality
- Try using a larger Vosk model for better accuracy
- Ensure your audio is clear with minimal background noise
- Check that you're using the correct language model

### Subtitles Not Appearing
- Verify that `subtitles.srt` was created successfully
- Check that the SRT file contains text
- Try adjusting the font size or color for better visibility

## Deactivating Virtual Environment

When you're done working:
```bash
deactivate
```

## License

This project uses:
- **Vosk** (Apache 2.0 License)
- **FFmpeg** (LGPL/GPL License)

## Support

For issues or questions:
- Vosk Documentation: [https://alphacephei.com/vosk/](https://alphacephei.com/vosk/)
- FFmpeg Documentation: [https://ffmpeg.org/documentation.html](https://ffmpeg.org/documentation.html)

---

**Happy video editing! 🎥✨**