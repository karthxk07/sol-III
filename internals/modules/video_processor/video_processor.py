import subprocess
import json
import os

def get_video_dimensions(video_path):
    """Get video width and height using ffprobe."""
    cmd = [
        'ffprobe', '-v', 'error',
        '-select_streams', 'v:0',
        '-show_entries', 'stream=width,height',
        '-of', 'json',
        video_path
    ]
    result = subprocess.run(cmd, capture_output=True, text=True)
    data = json.loads(result.stdout)
    width = data['streams'][0]['width']
    height = data['streams'][0]['height']
    return width, height

def calculate_crop_params(width, height, target_ratio=(9, 16)):
    """Calculate crop parameters for 9:16 ratio."""
    current_ratio = width / height
    target_ratio_val = target_ratio[0] / target_ratio[1]
    
    if current_ratio > target_ratio_val:
        # Video is wider, crop width
        new_width = int(height * target_ratio_val)
        new_height = height
        x_offset = (width - new_width) // 2
        y_offset = 0
    else:
        # Video is taller, crop height
        new_width = width
        new_height = int(width / target_ratio_val)
        x_offset = 0
        y_offset = (height - new_height) // 2
    
    return new_width, new_height, x_offset, y_offset

def extract_audio(video_path, audio_path):
    """Extract audio from video."""
    cmd = [
        'ffmpeg', '-i', video_path,
        '-vn', '-acodec', 'pcm_s16le',
        '-ar', '16000', '-ac', '1',
        audio_path, '-y'
    ]
    subprocess.run(cmd, check=True)

def transcribe_audio_vosk(audio_path):
    """Transcribe audio using Vosk (lightweight, offline)."""
    try:
        from vosk import Model, KaldiRecognizer
        import wave
        
        # Download model first: 
        # https://alphacephei.com/vosk/models
        # Extract to ./model directory
        model_path = "model"
        
        if not os.path.exists(model_path):
            print("Error: Vosk model not found!")
            print("Download from: https://alphacephei.com/vosk/models")
            print("Extract to './model' directory")
            return None
        
        model = Model(model_path)
        wf = wave.open(audio_path, "rb")
        rec = KaldiRecognizer(model, wf.getframerate())
        rec.SetWords(True)
        
        results = []
        while True:
            data = wf.readframes(4000)
            if len(data) == 0:
                break
            if rec.AcceptWaveform(data):
                result = json.loads(rec.Result())
                if 'result' in result:
                    results.append(result)
        
        # Final result
        final_result = json.loads(rec.FinalResult())
        if 'result' in final_result:
            results.append(final_result)
        
        return results
    except ImportError:
        print("Vosk not installed. Install with: pip install vosk")
        return None
    except Exception as e:
        print(f"Transcription error: {e}")
        return None

def format_time(seconds):
    """Convert seconds to SRT time format."""
    hours = int(seconds // 3600)
    minutes = int((seconds % 3600) // 60)
    secs = int(seconds % 60)
    millis = int((seconds % 1) * 1000)
    return f"{hours:02d}:{minutes:02d}:{secs:02d},{millis:03d}"

def create_srt_from_vosk(results, srt_path, max_words_per_subtitle=10):
    """Create SRT file from Vosk results."""
    with open(srt_path, 'w', encoding='utf-8') as f:
        subtitle_index = 1
        
        for result in results:
            if 'result' not in result:
                continue
                
            words = result['result']
            if not words:
                continue
            
            # Group words into subtitles
            i = 0
            while i < len(words):
                # Take up to max_words_per_subtitle words
                chunk = words[i:i + max_words_per_subtitle]
                
                start_time = chunk[0]['start']
                end_time = chunk[-1]['end']
                text = ' '.join([w['word'] for w in chunk])
                
                f.write(f"{subtitle_index}\n")
                f.write(f"{format_time(start_time)} --> {format_time(end_time)}\n")
                f.write(f"{text}\n\n")
                
                subtitle_index += 1
                i += max_words_per_subtitle

def crop_and_add_subtitles(input_video, output_video, srt_path):
    """Crop video to 9:16 and add subtitles."""
    # Get video dimensions
    width, height = get_video_dimensions(input_video)
    print(f"Original dimensions: {width}x{height}")
    
    # Calculate crop parameters
    new_width, new_height, x_offset, y_offset = calculate_crop_params(width, height)
    print(f"Cropping to: {new_width}x{new_height} at offset ({x_offset}, {y_offset})")
    
    # Build ffmpeg command
    crop_filter = f"crop={new_width}:{new_height}:{x_offset}:{y_offset}"
    subtitle_filter = f"subtitles={srt_path}:force_style='Fontname=Helvetica,Fontsize=15,PrimaryColour=&H00FFFFFF,OutlineColour=&H00000000,Outline=2,Shadow=1'"
    
    cmd = [
        'ffmpeg', '-i', input_video,
        '-vf', f"{crop_filter},{subtitle_filter}",
        '-c:a', 'copy',
        output_video, '-y'
    ]
    
    subprocess.run(cmd, check=True)

def crop_only(input_video, output_video):
    """Crop video to 9:16 without subtitles."""
    width, height = get_video_dimensions(input_video)
    new_width, new_height, x_offset, y_offset = calculate_crop_params(width, height)
    crop_filter = f"crop={new_width}:{new_height}:{x_offset}:{y_offset}"
    
    cmd = [
        'ffmpeg', '-i', input_video,
        '-vf', crop_filter,
        '-c:a', 'copy',
        output_video, '-y'
    ]
    subprocess.run(cmd, check=True)

def main():
    input_video = "input_video.mp4"
    output_video = "output_video.mp4"
    audio_path = "temp_audio.wav"
    srt_path = "subtitles.srt"
    
    print("Step 1: Extracting audio...")
    extract_audio(input_video, audio_path)
    
    print("Step 2: Transcribing audio with Vosk...")
    results = transcribe_audio_vosk(audio_path)
    
    if results:
        print("Step 3: Creating SRT file...")
        create_srt_from_vosk(results, srt_path)
        print(f"Subtitles saved to {srt_path}")
        
        print("Step 4: Cropping video and adding subtitles...")
        crop_and_add_subtitles(input_video, output_video, srt_path)
        
        print(f"Done! Output saved to {output_video}")
        
        # Cleanup
        os.remove(audio_path)
        os.remove(srt_path)
    else:
        print("Transcription failed. Cropping without subtitles...")
        crop_only(input_video, output_video)
        print(f"Done! Output saved to {output_video}")

if __name__ == "__main__":
    main()
