import whisper
from flask import Flask,request
import threading

app = Flask(__name__)

model = whisper.load_model("base")

lock = threading.Lock()

text = ""

def perform_transcription(path: str) -> str:
    ''' Read the file and perform the transcription'''
    audio = whisper.load_audio("sample.wav")
    audio = whisper.pad_or_trim(audio)

    mel = whisper.log_mel_spectrogram(audio).to(model.device)

    _, probs = model.detect_language(mel)
    print(f"Detected language: {max(probs, key=probs.get)}")

    options = whisper.DecodingOptions(fp16 = False)
    return whisper.decode(model, mel, options)

@app.route('/',methods=['GET', 'POST','OPTIONS'])
def recieve_file():

    print("Got request")

    # Enter the critical section
    lock.acquire()

    # Save request into the file
    file = open("sample.wav", "wb")
    file.write(request.data)
    file.close()

    # Perform speech-to-text
    result = perform_transcription("sample.wav")

    # Append the generated text to the final text
    text = text + " " + result.text
        
    # Save the text in a txt file
    text_file = open("Output.txt", "w")
    text_file.write(text)
    text_file.close()

    lock.release()

    # Return the text to the frontend
    return '<h1>{}</h1>'.format(text)


if __name__ == "__main__":
    app.run(debug=True)