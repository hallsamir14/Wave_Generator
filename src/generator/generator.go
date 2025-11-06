package generator

import (
	"encoding/binary"
	"errors"
	"math"
)

type Waveform int

const (
	Sine Waveform = iota
	Square
	Triangle
)

// Generate wave given args
// TODO anyway to reduce repitition here?
/*
waveData params:
numSamples uint, sampleRate uint, frequency float64, amplitude float64, data []byte,
*/
func GenerateWave(waveType Waveform, numSamples uint, sampleRate uint, frequency float64, amplitude float64, data []byte) ([]byte, error) {
	clampAmplitude(&amplitude)

	switch waveType {
	case Sine:
		return generateSinWave(numSamples, sampleRate, frequency, amplitude, data), nil
	case Square:
		return generateSquareWave(numSamples, sampleRate, frequency, amplitude, data), nil
	case Triangle:
		return generateTriangleWave(numSamples, sampleRate, frequency, amplitude, data), nil
	}
	return nil, errors.New("invalid wave type")
}

func clampAmplitude(amplitude *float64) {
	// clamp amplitude to [0,1] to avoid clipping
	if *amplitude < 0 {
		*amplitude = 0
	}
	if *amplitude > 1 {
		*amplitude = 1
	}
}
func generateSinWave(numSamples uint, sampleRate uint, frequency float64, amplitude float64, data []byte) []byte {

	maxInt16 := float64(math.MaxInt16)
	for i := uint(0); i < numSamples; i++ {
		t := float64(i) / float64(sampleRate)
		sample := amplitude * math.Sin(2*math.Pi*frequency*t)
		intSample := int16(sample * maxInt16)
		offset := i * 2
		binary.LittleEndian.PutUint16(data[offset:], uint16(intSample))
	}
	return data
}
func generateTriangleWave(numSamples uint, sampleRate uint, frequency float64, amplitude float64, data []byte) []byte {

	//Use a phase accumlator to avoid expensive math.Mod each sample
	phase := 0.0
	phaseInc := frequency / float64(sampleRate) // increment per sample (wraps at 1.0)
	maxInt16 := float64(math.MaxInt16)

	for i := uint(0); i < numSamples; i++ {
		/*triangle wave (range -1..1):
		triangle(phase) = 4*abs(phase - 0.5) - 1
		where phase is [0,1)
		*/
		waveformValue := 4*math.Abs(phase-0.5) - 1.0

		// scale by amplitude and convert to 16-bit signed
		intSample := int16(waveformValue * amplitude * maxInt16)

		// write little-endian
		offset := i * 2
		binary.LittleEndian.PutUint16(data[offset:offset+2], uint16(intSample))

		// robust wrap for large phaseInc
		phase += phaseInc
		if phase >= 1.0 {
			phase -= math.Floor(phase)
		}
	}
	return data
}

func generateSquareWave(numSamples uint, sampleRate uint, frequency float64, amplitude float64, data []byte) []byte {

	phase := 0.0
	phaseInc := frequency / float64(sampleRate) // increment per sample (wraps at 1.0)
	maxInt16 := float64(math.MaxInt16)

	for i := uint(0); i < numSamples; i++ {
		// Square: +1 for first half of period, -1 for second half
		var waveformValue float64
		if phase < 0.5 {
			waveformValue = 1.0
		} else {
			waveformValue = -1.0
		}

		// Scale by amplitude and convert to 16-bit signed
		intSample := int16(waveformValue * amplitude * maxInt16)

		// write little-endian
		offset := i * 2
		binary.LittleEndian.PutUint16(data[offset:offset+2], uint16(intSample))

		// increment & robust wrap for large phaseInc
		phase += phaseInc
		if phase >= 1.0 {
			phase -= math.Floor(phase)
		}
	}

	return data
}
