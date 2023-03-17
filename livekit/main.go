package main

import (
	"fmt"
	"github.com/pion/webrtc/v3"
	"math/rand"
	"os"
	"time"

	lksdk "github.com/livekit/server-sdk-go"
)

func main() {
	host := "ws://localhost:7880"
	apiKey := "devkey"
	apiSecret := "secret"
	roomName := "my-first-room"
	identity := "botuser"
	roomCB := &lksdk.RoomCallback{
		ParticipantCallback: lksdk.ParticipantCallback{
			OnTrackSubscribed:   trackSubscribed,
			OnTrackUnsubscribed: trackUnSubscribed,
		},
	}
	_, err := lksdk.ConnectToRoom(host, lksdk.ConnectInfo{
		APIKey:              apiKey,
		APISecret:           apiSecret,
		RoomName:            roomName,
		ParticipantIdentity: identity,
	}, roomCB)
	if err != nil {
		panic(err)
	}

	time.Sleep(time.Hour)
}

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

func RandStringBytes(n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}
	return string(b)
}

func trackSubscribed(track *webrtc.TrackRemote, publication *lksdk.RemoteTrackPublication, rp *lksdk.RemoteParticipant) {

	currentBuffer := make([]byte, 0)
	counter := 0

	if !("audio/opus" == track.Codec().MimeType && rp.Name() == "user1") {
		return
	}

	fmt.Println(fmt.Sprintf("User : %s", rp.Name()))
	fmt.Print(fmt.Sprintf("Codec : %s", track.Codec().MimeType))

	filename := RandStringBytes(20)

	file, err := os.OpenFile(filename+".opus", os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0777)
	if err != nil {
		panic(err)
	}

	for {
		packet, _, err := track.ReadRTP()

		if err == nil {
			currentBuffer = append(currentBuffer, packet.Payload...)
			counter++

			_, err = file.Write(packet.Payload)
			if err != nil {
				panic(err)
			}

			if counter%1000 == 0 {
				file.Close()
				fmt.Println("Close")
				return
			}
		}
	}

}

func trackUnSubscribed(track *webrtc.TrackRemote, publication *lksdk.RemoteTrackPublication, rp *lksdk.RemoteParticipant) {
	fmt.Println(rp.Name())
}
