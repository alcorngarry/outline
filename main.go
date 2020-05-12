package main

import (
	"fmt"

	"github.com/pion/webrtc"
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

	dataChannel, err := peerConnection.CreateDataChannel("data", nil)
	if err != nil {
		panic(err) // Please handle your errors correctly!
	}

	peerConnection.OnDataChannel(func(dataChannel *webrtc.DataChannel) {
		fmt.Printf("New DataChannel %s %d\n", dataChannel.Label, dataChannel.ID)

		// Handle data channel
	})

	dataChannel.OnOpen(func() {
		fmt.Printf("Data channel '%s'-'%d' open.\n", dataChannel.Label(), dataChannel.ID())

		// Now we can start sending data.
	})

	// Register text message handling
	dataChannel.OnMessage(func(msg webrtc.DataChannelMessage) {
		fmt.Printf("Message from DataChannel '%s': '%s'\n", dataChannel.Label(), string(msg.Data))

		// Handle the message here
	})

	// Send the message as text
	err1 := dataChannel.SendText("hello")
	if err1 != nil {
		panic(err1) // Please handle your errors correctly!
	}

}
