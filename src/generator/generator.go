package generator

import (
	"encoding/binary"
	"errors"
	"math"
)

// Generate wave given args
// TODO anyway to reduce repitition here?
/*
waveData params:
numSamples uint, sampleRate uint, frequency float64, amplitude float64, data []byte,
*/
func GenerateWave(waveType string, numSamples uint, sampleRate uint, frequency float64, amplitude float64, data []byte) ([]byte, error) {
	if waveType == "Sin" {
		return generateSinWave(numSamples, sampleRate, frequency, amplitude, data), nil
	} else if waveType == "Triangle" {
		return generateTriangleWave(numSamples, sampleRate, frequency, amplitude, data), nil
	}
	return nil, errors.New("invalid wave type")
}
func generateSinWave(numSamples uint, sampleRate uint, frequency float64, amplitude float64, data []byte) []byte {
	for i := uint(0); i < numSamples; i++ {
		t := float64(i) / float64(sampleRate)
		sample := amplitude * math.Sin(2*math.Pi*frequency*t)
		intSample := int16(sample * 32767) //TODO why is this here??
		binary.LittleEndian.PutUint16(data[i*2:], uint16(intSample))
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
		where phase is [0,1)*/
		v := 4*math.Abs(phase-0.5) - 1.0

		// scale by amplitude and convert to 16-bit signed
		sample := int16(v * amplitude * maxInt16)

		// write little-endian
		off := i * 2
		binary.LittleEndian.PutUint16(data[off:off+2], uint16(sample))

		phase += phaseInc
		if phase >= 1.0 {
			phase -= math.Floor(phase) // robust wrap for large phaseInc
		}
	}
	return data
}

func generateSquareWave(numSamples uint, sampleRate uint, frequency float64, amplitude float64, data []byte) []byte {
	return data
}
