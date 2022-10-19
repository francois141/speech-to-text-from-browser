import whisper
from flask import Flask,request

app = Flask(__name__)

model = whisper.load_model("base")

text = ""

@app.route('/',methods=['GET', 'POST','OPTIONS'])
def translate():

    print("Got request")
    raw_data = request.data
    file = open("sample.wav", "wb")
    file.write(raw_data)
    file.close()
    result = model.transcribe("sample.wav")

    global text
    text = text + " " + result["text"]

    text_file = open("Output.txt", "w")
    text_file.write(text)
    text_file.close()

    return '<h1>{}</h1>'.format(text)


if __name__ == "__main__":
    app.run(debug=True)