package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/eiannone/keyboard"
	"github.com/faiface/beep"
	"github.com/faiface/beep/mp3"
	"github.com/faiface/beep/speaker"
	"github.com/faiface/beep/wav"
)

type Music struct {
	Path     string
	Format   string
	Duration time.Duration
}

var (
	ctrl      *beep.Ctrl
	format    beep.Format
	streamer  beep.StreamSeeker
	playlist  []Music
	isPaused  bool
	isPlaying bool
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: popmusic [music_directory...]")
		os.Exit(1)
	}

	// Scan directories for music files
	for _, dir := range os.Args[1:] {
		scanDirectory(dir)
	}

	if len(playlist) == 0 {
		fmt.Println("No music files found in the specified directories")
		os.Exit(1)
	}

	// Initialize keyboard
	if err := keyboard.Open(); err != nil {
		log.Fatal(err)
	}
	defer keyboard.Close()

	fmt.Println("\nControls:")
	fmt.Println("Space: Play/Pause")
	fmt.Println("q: Quit")
	fmt.Println("\nPlaylist:")
	for i, music := range playlist {
		fmt.Printf("%d. %s\n", i+1, filepath.Base(music.Path))
	}

	// Start playing the first song
	playMusic(0)

	// Handle keyboard events
	for {
		char, key, err := keyboard.GetKey()
		if err != nil {
			log.Fatal(err)
		}

		switch {
		case char == 'q' || char == 'Q':
			fmt.Println("\nGoodbye!")
			return
		case key == keyboard.KeySpace:
			if isPlaying {
				if isPaused {
					ctrl.Paused = false
					isPaused = false
					fmt.Println("\nResumed")
				} else {
					ctrl.Paused = true
					isPaused = true
					fmt.Println("\nPaused")
				}
			}
		}
	}
}

func scanDirectory(root string) {
	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			ext := strings.ToLower(filepath.Ext(path))
			if ext == ".mp3" || ext == ".wav" {
				playlist = append(playlist, Music{
					Path:   path,
					Format: ext[1:],
				})
			}
		}
		return nil
	})
	if err != nil {
		log.Printf("Error scanning directory %s: %v\n", root, err)
	}
}

func playMusic(index int) {
	if index >= len(playlist) {
		return
	}

	f, err := os.Open(playlist[index].Path)
	if err != nil {
		log.Printf("Error opening file %s: %v\n", playlist[index].Path, err)
		return
	}

	var err2 error
	switch playlist[index].Format {
	case "mp3":
		streamer, format, err2 = mp3.Decode(f)
	case "wav":
		streamer, format, err2 = wav.Decode(f)
	}

	if err2 != nil {
		log.Printf("Error decoding file %s: %v\n", playlist[index].Path, err2)
		return
	}

	speaker.Init(format.SampleRate, format.SampleRate.N(time.Second/10))
	ctrl = &beep.Ctrl{Streamer: streamer}
	speaker.Play(ctrl)
	isPlaying = true

	fmt.Printf("\nNow playing: %s\n", filepath.Base(playlist[index].Path))
}
