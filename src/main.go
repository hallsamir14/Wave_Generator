package main

import (
	"wave_generator/generator"
)

func main() {
	const sampleRate = 44100
	const frequency = 220.0
	const amplitude = 0.3
	const durationSec = 4

	// Generate single wave sample
	numSamples := uint(sampleRate * durationSec)
	data := make([]byte, numSamples*2)
	generator.GenerateWave(generator.Triangle, numSamples, sampleRate, frequency, amplitude, data)
	writeFile(sampleRate, data)

	/*Generate continious stream of wave samples
	output := make([]int16, 0, 128)
	waveGenerator := generator.NewWaveGenerator(sampleRate, frequency, amplitude, generator.Sine)
	waveGenerator.GenerateInt16(64, output)
	*/
}
