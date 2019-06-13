package main

import (
	"./message"
	"fmt"
	"github.com/golang/protobuf/proto"
	"net"
	"os"
)

const (
	keyUrl = "0.0.0.0:3333"
	cacheUrl = "0.0.0.0:3334"
	writeUrl = "0.0.0.0:3335"
	google = "https://www.googleapis.com/youtube/v3/channels?part=statistics&key=%s&id=%s"
	bufSize = 2000
)

var (
	path string
)

func init() {
	fmt.Println("Random poller started")
	path = os.Args[0]
	fmt.Println("Using path", path)
}

func getKey() string {
	conn, err := net.Dial("tcp4", keyUrl)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	defer func() {
		err := conn.Close()
		if err != nil {
			fmt.Println(err)
			os.Exit(2)
		}
	}()

	{
		_, err := conn.Write([]byte{1})
		if err != nil {
			fmt.Println(err)
			os.Exit(3)
		}
	}

	var key [24]byte
	{
		var tmp []byte
		n, err := conn.Read(tmp)
		if err != nil {
			fmt.Println(err)
			os.Exit(4)
		}

		if n != 24 {
			fmt.Println("Bad key size")
			os.Exit(5)
		}

		for i := 0; i < 24; i++ {
			key[i] = tmp[i]
		}
	}

	keyStr := string(key[:])
	fmt.Println("Using key", keyStr)
	return keyStr
}

func getChannels() message.ChannelMessage {
	conn, err := net.Dial("tcp4", keyUrl)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	defer func() {
		err := conn.Close()
		if err != nil {
			fmt.Println(err)
			os.Exit(2)
		}
	}()

	{
		_, err := conn.Write([]byte{1})
		if err != nil {
			fmt.Println(err)
			os.Exit(3)
		}
	}

	{
		var tmp []byte
		n, err := conn.Read(tmp)
		if err != nil {
			fmt.Println(err)
			os.Exit(5)
		}

		var msg message.ChannelMessage
		buf := tmp[:n]

		err = proto.Unmarshal(buf, &msg)
		if err != nil {
			fmt.Println(err)
			os.Exit(6)
		}

		fmt.Println("Retrieved channels:", msg)
		return msg
	}
}

func getMetrics(key string, msg message.ChannelMessage) []byte {
	return []byte{}
}

func sendPayload(bytes []byte) {

}

func main() {
	for {
		key := getKey()
		channels := getChannels()

		bytes := getMetrics(key, channels)
		sendPayload(bytes)
	}
}
