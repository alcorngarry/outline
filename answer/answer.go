package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/signal"
	"strings"
	"time"

	"github.com/pion/webrtc"
	"github.com/sacOO7/gowebsocket"
)

func main() {
	// webcam, _ := gocv.OpenVideoCapture(0)
	// window := gocv.NewWindow("Hello")
	// img := gocv.NewMat()

	// for {
	// 	webcam.Read(&img)
	// 	window.IMShow(img)
	// 	gocv.WaitKey(1)
	// }

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)
	socket := gowebsocket.New("ws://192.168.1.21:9090")

	config := webrtc.Configuration{
		ICEServers: []webrtc.ICEServer{
			{
				URLs: []string{"stun:stun.l.google.com:19302"},
			},
		},
	}

	peerConnection, err := webrtc.NewPeerConnection(config)
	if err != nil {
		panic(err)
	}

	// Set the handler for ICE connection state
	// This will notify you when the peer has connected/disconnected
	peerConnection.OnICEConnectionStateChange(func(connectionState webrtc.ICEConnectionState) {
		fmt.Printf("ICE Connection State has changed: %s\n", connectionState.String())
	})

	// Register data channel creation handling
	peerConnection.OnDataChannel(func(d *webrtc.DataChannel) {
		fmt.Printf("New DataChannel %s %d\n", d.Label(), d.ID())

		// Register channel opening handling
		d.OnOpen(func() {
			fmt.Printf("Data channel '%s'-'%d' open. Random messages will now be sent to any connected DataChannels every 5 seconds\n", d.Label(), d.ID())

			for range time.NewTicker(5 * time.Second).C {
				message := "hello world"
				fmt.Printf("Sending '%s'\n", message)

				// Send the message as text
				sendTextErr := d.SendText(message)
				if sendTextErr != nil {
					panic(sendTextErr)
				}
			}
		})

		// Register text message handling
		d.OnMessage(func(msg webrtc.DataChannelMessage) {
			fmt.Printf("Message from DataChannel '%s': '%s'\n", d.Label(), string(msg.Data))
		})
	})

	socket.Connect()
	socket.OnTextMessage = func(message string, socket gowebsocket.Socket) {
		var offer webrtc.SessionDescription
		err := json.NewDecoder(strings.NewReader(message)).Decode(&offer)
		if err != nil {
			panic(err)
		}

		err = peerConnection.SetRemoteDescription(offer)
		if err != nil {
			panic(err)
		}
	}

	// Create answer
	answer, err := peerConnection.CreateAnswer(nil)

	// Sets the LocalDescription, and starts our UDP listeners
	err = peerConnection.SetLocalDescription(answer)

	fmt.Println(answer)

	b := new(bytes.Buffer)
	json.NewEncoder(b).Encode(answer)
	socket.SendText(b.String())

	//interrupting sequence
	for {
		select {
		case <-interrupt:
			log.Println("interrupt")
			socket.Close()
			return
		}
	}
	// Block forever
	select {}

}
