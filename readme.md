# PoC - Speech to text with browser recording - Whisper model


<div align="center">

![Python](https://img.shields.io/badge/python-3670A0?logo=python&logoColor=ffdd54)
![JavaScript](https://img.shields.io/badge/javascript-%23323330.svg?logo=javascript&logoColor=%23F7DF1E)
[![Open Source Love](https://badges.frapsoft.com/os/v1/open-source.png?v=103)](https://github.com/ellerbrock/open-source-badges/)
![Edge](https://img.shields.io/badge/Edge-0078D7)
![Firefox](https://img.shields.io/badge/Firefox-FF7139)
![Google Chrome](https://img.shields.io/badge/Google%20Chrome-4285F4)

</div>

## Description of the project

This is a small proof of concept to test how we can record sound from the navigator and perform speech to text from whisper. Whisper is a multilingual speech to text model from OpenAI.


<div align="center">

[Whisper blog from openAI](https://openai.com/blog/whisper/) • [Github repo of Whisper](https://github.com/openai/whisper)

</div>

## How to run the code

Run the backend

``` bash
python app.py
```

Run the frontend

``` bash
python -m http.server 8000
```

## Live demonstration

TODO: ADD GIF HERE

### Documentation

Convert blobs into mp3: https://medium.com/jeremy-gottfrieds-tech-blog/javascript-tutorial-record-audio-and-encode-it-to-mp3-2eedcd466e78
Example for livekit : https://github.com/livekit/server-sdk-go/blob/main/examples/filesaver/main.go
