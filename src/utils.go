package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"
	"wave_generator/writer"
)

func writeFile(sampleRate int, data []byte) error {
	/*Create directory './tones'
	Append '.wav' extension after timestamp
	Output file example: ./tones/Triangle_20251026213909..wav
	*/
	dir, extension := "tones", ".wav"
	timestamp := time.Now().Format("20060102150405")
	filename := fmt.Sprintf("%v.%v", timestamp, extension)
	filePath := filepath.Join(dir, filename)

	fmt.Printf("Writing wav file %v", filename)

	err := os.MkdirAll(filepath.Dir(filePath), 0755)
	if err != nil {
		fmt.Println("Error occured creating directory")
	}

	// Create WAV file
	file, err := os.Create(filePath)

	defer file.Close()

	// Write WAV header
	writer.WriteWavHeader(file, uint32(len(data)), sampleRate, 1, 16)

	//Write PCM data
	writer.WritePCMData(data, file, err)

	if err != nil {
		log.Fatal(err)
	}

	return err

}
