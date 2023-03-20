import whisper
from flask import Flask,request
import threading
import subprocess
import numpy as np

import subprocess
import logging

app = Flask(__name__)

model = whisper.load_model("base")
lock = threading.Lock()
text = ""

def bytesStreamToTensor(byteStream) -> np.array:
    ffmpeg = 'ffmpeg'

    cmd = [ffmpeg,
        '-i', 'pipe:',
        '-f', 's16le',
        '-acodec', 'pcm_s16le',
        '-ac','1',
        '-ar','16000',
        'pipe:']

    pipe = subprocess.Popen(cmd, stdout=subprocess.PIPE,stdin=subprocess.PIPE)
    out = pipe.communicate(input=request.data)[0]    
    pipe.wait()

    return np.frombuffer(out, np.int16).flatten().astype(np.float32) / 32768.0

def audioTensorToSpectrogram(tensor: np.array) -> np.array:
    audio = whisper.pad_or_trim(tensor)
    return whisper.log_mel_spectrogram(audio).to(model.device)
    
def getLanguage(spectogram: np.array) -> np.array:
    _, probs = model.detect_language(spectogram)
    return max(probs, key=probs.get)
    

def performTranscription() -> str:
    ''' Read the byte stream and perform the transcription'''
    audioTensor = bytesStreamToTensor(request.data)
    spectogram = audioTensorToSpectrogram(audioTensor)

    print("Detected language: {}".format(getLanguage(spectogram=spectogram)))

    decodingOptions = whisper.DecodingOptions(fp16 = False)
    return whisper.decode(model, spectogram, decodingOptions).text

def writeTextToOutputFile(text: str):
    text_file = open("Output.txt", "w")
    text_file.write(text)
    text_file.close()

@app.route('/',methods=['GET', 'POST','OPTIONS'])
def recieveByteStream():
    lock.acquire()
    _handleByteStream()
    lock.release()

    return text

def _handleByteStream():
    result = performTranscription()

    global text
    text += result

    writeTextToOutputFile(text)



if __name__ == "__main__":
    app.run(debug=True)