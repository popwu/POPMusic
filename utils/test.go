package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/faiface/beep"
	"github.com/faiface/beep/mp3"
	"github.com/faiface/beep/speaker"
)

func main() {
	filename1 := "a.mp3"
	filename2 := "b.mp3"

	// 播放 a.mp3 3秒后，播放 b.mp3
	go playMusic(filename1)
	time.Sleep(3 * time.Second)
	go playMusic(filename2)

	// 等待播放完成
	select {}
}

func playMusic(filename string) {
	f, err := os.Open(filename)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	var streamer beep.StreamSeekCloser
	var format beep.Format

	if filename[len(filename)-4:] == ".mp3" {
		streamer, format, err = mp3.Decode(f)
		if err != nil {
			log.Fatal(err)
		}
	} else {
		log.Fatal("Unsupported file format")
	}
	defer streamer.Close()

	speaker.Init(format.SampleRate, format.SampleRate.N(time.Second/10))
	done := make(chan bool)
	speaker.Play(beep.Seq(streamer, beep.Callback(func() {
		done <- true
	})))

	fmt.Printf("Now playing: %s\n", filename)
	<-done
}
