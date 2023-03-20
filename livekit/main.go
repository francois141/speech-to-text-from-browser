package main

import (
	"bytes"
	"fmt"
	"math/rand"
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

var (
	host, apiKey, apiSecret, roomName, identity string
)

func main() {
	host := "ws://localhost:7880"
	apiKey := "devkey"
	apiSecret := "secret"
	roomName := "my-first-room"
	identity := "botuser"

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
	fileName := fmt.Sprintf("%s-%s", rp.Identity(), track.ID())
	fmt.Println("write track to file ", fileName)
	NewTrackWriter(track, rp.WritePLI, fileName)
}

const (
	maxAudioLate = 200 // 4s for audio
)

type TrackWriter struct {
	sb     *samplebuilder.SampleBuilder
	writer media.Writer
	track  *webrtc.TrackRemote
}

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

func RandStringBytes(n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}
	return string(b)
}

var counter = 0
var byteBuffer bytes.Buffer

func NewTrackWriter(track *webrtc.TrackRemote, pliWriter lksdk.PLIWriter, fileName string) (*TrackWriter, error) {
	var (
		sb     *samplebuilder.SampleBuilder
		writer media.Writer
		err    error
	)

	switch {
	case strings.EqualFold(track.Codec().MimeType, "audio/opus"):
		sb = samplebuilder.New(maxAudioLate, &codecs.OpusPacket{}, track.Codec().ClockRate)
		writer, err = oggwriter.NewWith(&byteBuffer, 48000, track.Codec().Channels)

	default:
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
	go t.start()
	return t, nil
}

func (t *TrackWriter) start() {
	defer t.writer.Close()
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

		fmt.Println(counter)

		if counter > 100 {

			byteBufferPacket2 := make([]byte, 0, 0)
			byteBufferPacket2 = byteBuffer.Bytes()

			// HTTP endpoint

			posturl := "http://127.0.0.1:5000"

			request, err := http.NewRequest("POST", posturl, bytes.NewBuffer(byteBufferPacket2))
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

			writer, err := oggwriter.NewWith(&byteBuffer, 48000, 2)
			if err != nil {
				panic(err)
			}

			t.writer = writer
		}
	}
}
