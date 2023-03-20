// collect DOMs
const display = document.querySelector('.display')
const controllerWrapper = document.querySelector('.controllers')

let mediaRecorder, chunks = []

// mediaRecorder setup for audio
if (navigator.mediaDevices && navigator.mediaDevices.getUserMedia) {
    console.log('mediaDevices supported..')

    navigator.mediaDevices.getUserMedia({
        audio: true
    }).then(stream => {

        console.log("Starting the program")

        mediaRecorder = new MediaRecorder(stream)

        mediaRecorder.ondataavailable = (e) => {
            chunks.push(e.data)
        }

        mediaRecorder.onstop = () => {
            let blob = new Blob(chunks, { 'type': 'audio/wav;' })
            chunks = []

            uploadBlob(blob)
        }

        record();

        // Create the loop
        window.setInterval(function () {
            console.log('Send a text to the backend')
            stopRecording();
            record();
        }, 5000);

    }).catch(error => {
        console.log('Following error has occured : ', error)
    })
} else {
    console.log('Failed to start the program')
}

const record = () => {
    mediaRecorder.start()
}

const stopRecording = () => {
    mediaRecorder.stop()
}

function uploadBlob(blob) {

    // Creating a new blob  
    // Hostname and port of the local server
    fetch('http://localhost:5000', {
        // HTTP request type
        method: "POST",
        mode: 'no-cors',
        body: blob
    })
        .then(response => console.log('Blob Uploaded'))
        .catch(err => alert(err));
}