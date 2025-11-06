package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"
	"wave_generator/generator"
	"wave_generator/writer"
)

func main() {
	//Paramatize these
	const sampleRate = 44100
	const frequency = 220.0
	const amplitude = 0.3
	const durationSec = 4

	// 16-bit PCM (2 bytes per sample)
	numSamples := uint(sampleRate * durationSec)
	data := make([]byte, numSamples*2)
	data, generatorErr := generator.GenerateWave(generator.Triangle, numSamples, sampleRate, frequency, amplitude, data)

	if generatorErr != nil {
		fmt.Println("Error generating wave")
	}

	/*Create directory './tones'
	Append '.wav' extension after timestamp
	Output file example: ./tones/Triangle_20251026213909..wav
	*/
	dir, extension := "tones", ".wav"
	timestamp := time.Now().Format("20060102150405")
	filename := fmt.Sprintf("%v.%v", timestamp, extension)
	filePath := filepath.Join(dir, filename)

	err := os.MkdirAll(filepath.Dir(filePath), 0755)
	if err != nil {
		fmt.Println("Error occured creating directory")
	}

	// Create WAV file
	file, err := os.Create(filePath)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	// Write WAV header
	writer.WriteWavHeader(file, uint32(len(data)), sampleRate, 1, 16)

	//Write PCM data
	writer.WritePCMData(data, file, err)
}
