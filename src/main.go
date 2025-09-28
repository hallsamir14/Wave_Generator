package main

import (
	"log"
	"os"
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
	numSamples := sampleRate * durationSec
	data := make([]byte, numSamples*2)

	data = generator.GenerateWave(numSamples, sampleRate, frequency, amplitude, data)

	// Create WAV file
	file, err := os.Create("tone.wav")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	// Write WAV header
	writer.WriteWavHeader(file, uint32(len(data)), sampleRate, 1, 16)

	//Write PCM data
	writer.WritePCMData(data, file, err)
}
