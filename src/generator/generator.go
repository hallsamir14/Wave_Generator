package generator

import (
	"encoding/binary"
	"errors"
	"math"
)

// Waveform represents the waveform type to synthesize.
type Waveform int

const (
	Sine Waveform = iota
	Square
	Triangle
)

// WaveGenerator generates audio samples and maintains phase continuity between calls.
type WaveGenerator struct {
	SampleRate int     // samples per second, e.g. 44100
	Frequency  float64 // Hz
	Amplitude  float64 // 0..1
	Phase      float64 // 0..1 (fractional phase within period)
	Waveform   Waveform
}

// NewWaveGenerator returns a configured generator. Amplitude is clamped [0,1].
func NewWaveGenerator(sampleRate int, freq float64, amp float64, wf Waveform) *WaveGenerator {
	if sampleRate <= 0 {
		sampleRate = 44100
	}
	return &WaveGenerator{
		SampleRate: sampleRate,
		Frequency:  freq,
		Amplitude:  clamp(amp, 0.0, 1.0),
		Phase:      0.0,
		Waveform:   wf,
	}
}

/*
GenerateWave is a convenience wrapper that constructs a temporary WaveGenerator and writes
samples into the provided byte buffer. It's a simple stateless helper.
Stateless generator,generates single chuck of audio
*/

func GenerateWave(wf Waveform, numSamples uint, sampleRate uint, frequency float64, amplitude float64, data []byte) ([]byte, error) {
	// Validate byte buffer
	if int(numSamples)*2 > len(data) {
		return nil, errors.New("data buffer too small for requested number of samples")
	}
	g := NewWaveGenerator(int(sampleRate), frequency, amplitude, wf)
	if err := g.generateToPCMBytes(int(numSamples), data); err != nil {
		return nil, err
	}
	return data, nil
}

/*
GenerateInt16 fills the provided output slice with `numSamples` int16 samples.
len(out) must be >= numSamples. Returns an error on invalid inputs.
Generate audio in chucnks - continous audio
*/
func (g *WaveGenerator) GenerateInt16(numSamples int, out []int16) error {
	if numSamples < 0 {
		return errors.New("numSamples must be non-negative")
	}
	if g.SampleRate <= 0 {
		return errors.New("invalid sample rate")
	}
	if len(out) < numSamples {
		return errors.New("output slice too small")
	}

	phaseStep := g.Frequency / float64(g.SampleRate)
	maxInt16 := float64(math.MaxInt16)

	for i := 0; i < numSamples; i++ {
		var normalizedSample float64 // -1..+1

		switch g.Waveform {
		case Sine:
			// sine: sin(2*pi*phase)
			normalizedSample = math.Sin(2.0 * math.Pi * g.Phase)
		case Triangle:
			// triangle (range -1..+1): corrected formula
			normalizedSample = 1.0 - 4.0*math.Abs(g.Phase-0.5)
		case Square:
			// 50% duty square: +1 for first half, -1 for second half
			if g.Phase < 0.5 {
				normalizedSample = 1.0
			} else {
				normalizedSample = -1.0
			}
		default:
			return errors.New("unsupported waveform")
		}

		// scale by amplitude and convert to int16
		intSample := int16(normalizedSample * g.Amplitude * maxInt16)
		out[i] = intSample

		// advance phase and wrap robustly
		g.Phase += phaseStep
		if g.Phase >= 1.0 {
			g.Phase -= math.Floor(g.Phase)
		}
	}

	return nil
}

// GenerateToPCMBytes writes numSamples 16-bit little-endian PCM samples directly into out byte buffer.
// out must have length >= numSamples*2.
func (g *WaveGenerator) generateToPCMBytes(numSamples int, out []byte) error {
	if numSamples < 0 {
		return errors.New("numSamples must be non-negative")
	}
	if g.SampleRate <= 0 {
		return errors.New("invalid sample rate")
	}
	if len(out) < numSamples*2 {
		return errors.New("output byte buffer too small")
	}

	phaseStep := g.Frequency / float64(g.SampleRate)
	maxInt16 := float64(math.MaxInt16)

	for i := 0; i < numSamples; i++ {
		var normalizedSample float64

		switch g.Waveform {
		case Sine:
			normalizedSample = math.Sin(2.0 * math.Pi * g.Phase)
		case Triangle:
			normalizedSample = 1.0 - 4.0*math.Abs(g.Phase-0.5)
		case Square:
			if g.Phase < 0.5 {
				normalizedSample = 1.0
			} else {
				normalizedSample = -1.0
			}
		default:
			return errors.New("unsupported waveform")
		}

		intSample := int16(normalizedSample * g.Amplitude * maxInt16)
		offset := i * 2
		binary.LittleEndian.PutUint16(out[offset:offset+2], uint16(intSample))

		g.Phase += phaseStep
		if g.Phase >= 1.0 {
			g.Phase -= math.Floor(g.Phase)
		}
	}

	return nil
}

// clamp clamps v to [lo, hi].
func clamp(v, lo, hi float64) float64 {
	if v < lo {
		return lo
	}
	if v > hi {
		return hi
	}
	return v
}
