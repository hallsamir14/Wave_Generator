package main

import (
	"wave_generator/generator"
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
	generator.GenerateWave(generator.Triangle, numSamples, sampleRate, frequency, amplitude, data)
	writeFile(sampleRate, data)
}
