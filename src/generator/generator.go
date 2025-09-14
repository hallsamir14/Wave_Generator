package generator

import (
	"encoding/binary"
	"math"
)

// Generate wave given args
// TODO Correct usigned types for args
func GenerateWave(numSamples int, sampleRate int, frequency float64, amplitude float64, data []byte) []byte {
	for i := 0; i < numSamples; i++ {
		t := float64(i) / float64(sampleRate)
		sample := amplitude * math.Sin(2*math.Pi*frequency*t)
		intSample := int16(sample * 32767)
		binary.LittleEndian.PutUint16(data[i*2:], uint16(intSample))
	}
	return data
}
