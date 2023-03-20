package main

import (
	"bytes"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/pion/rtp/codecs"
	"github.com/pion/webrtc/v3"
	"github.com/pion/webrtc/v3/pkg/media"
	"github.com/pion/webrtc/v3/pkg/media/oggwriter"

	lksdk "github.com/livekit/server-sdk-go"
	"github.com/livekit/server-sdk-go/pkg/samplebuilder"
)

const (
	maxAudioLate uint16 = 200
	host         string = "ws://localhost:7880"
	apiKey       string = "devkey"
	apiSecret    string = "secret"
	roomName     string = "my-first-room"
	identity     string = "botuser"
	serverUrl    string = "http://127.0.0.1:5000"

	sampleRate   uint32 = 48000
	channelCount uint16 = 2
)

func main() {
	room, err := lksdk.ConnectToRoom(host, lksdk.ConnectInfo{
		APIKey:              apiKey,
		APISecret:           apiSecret,
		RoomName:            roomName,
		ParticipantIdentity: identity,
	}, &lksdk.RoomCallback{
		ParticipantCallback: lksdk.ParticipantCallback{
			OnTrackSubscribed: onTrackSubscribed,
		},
	})
	if err != nil {
		panic(err)
	}

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT)

	<-sigChan
	room.Disconnect()
}

func onTrackSubscribed(track *webrtc.TrackRemote, publication *lksdk.RemoteTrackPublication, rp *lksdk.RemoteParticipant) {
	NewTrackWriter(track)
}

type TrackWriter struct {
	sb     *samplebuilder.SampleBuilder
	writer media.Writer
	track  *webrtc.TrackRemote
}

func NewTrackWriter(track *webrtc.TrackRemote) (*TrackWriter, error) {
	var (
		sb     *samplebuilder.SampleBuilder
		writer media.Writer
		err    error
	)

	var byteBuffer bytes.Buffer

	if strings.EqualFold(track.Codec().MimeType, "audio/opus") {
		sb = samplebuilder.New(maxAudioLate, &codecs.OpusPacket{}, track.Codec().ClockRate)
		writer, err = oggwriter.NewWith(&byteBuffer, sampleRate, track.Codec().Channels)
	} else {
		return nil, nil
	}

	if err != nil {
		return nil, err
	}

	t := &TrackWriter{
		sb:     sb,
		writer: writer,
		track:  track,
	}

	go t.start(&byteBuffer)
	return t, nil
}

func (t *TrackWriter) start(byteBuffer *bytes.Buffer) {
	defer t.writer.Close()

	var counter = 0

	for {
		pkt, _, err := t.track.ReadRTP()
		if err != nil {
			break
		}
		t.sb.Push(pkt)

		for _, p := range t.sb.PopPackets() {
			t.writer.WriteRTP(p)
			counter++
		}

		if counter > 500 {

			byteBufferToSend := byteBuffer.Bytes()

			request, err := http.NewRequest("POST", serverUrl, bytes.NewBuffer(byteBufferToSend))
			if err != nil {
				panic(err)
			}

			client := &http.Client{}
			_, err = client.Do(request)
			if err != nil {
				panic(err)
			}

			byteBuffer.Reset()
			counter = 0

			writer, err := oggwriter.NewWith(byteBuffer, sampleRate, channelCount)
			if err != nil {
				panic(err)
			}

			t.writer = writer
		}
	}
}
