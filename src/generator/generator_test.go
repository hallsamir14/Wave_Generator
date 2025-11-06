package generator

import (
	"encoding/binary"
	"errors"
	"math"
	"testing"
)

// helper: read int16 sample at index i from little-endian PCM buffer
func int16At(data []byte, i uint) int16 {
	off := i * 2
	u := binary.LittleEndian.Uint16(data[off : off+2])
	return int16(u)
}

func TestGenerateWave_TableDriven(t *testing.T) {
	const Max = int16(math.MaxInt16)

	type expectMap map[uint]int16

	type tc struct {
		name       string
		waveType   Waveform
		numSamples uint
		sampleRate uint
		freq       float64
		amp        float64
		wantErr    bool
		expect     expectMap          // exact sample index -> expected int16 value
		validator  func([]byte) error // optional custom check
	}

	cases := []tc{
		{
			name:       "sin basic deterministic samples",
			waveType:   Sine,
			numSamples: 8,
			sampleRate: 8,
			freq:       1.0,
			amp:        1.0,
			expect: expectMap{
				0: 0,           // sin(0) == 0
				2: Max,         // sin(pi/2) == 1.0
				4: 0,           // sin(pi) == 0
				6: int16(-Max), // sin(3pi/2) == -1.0
			},
		},
		{
			name:       "triangle 4-sample period values (implementation-specific)",
			waveType:   Triangle,
			numSamples: 4,
			sampleRate: 4,
			freq:       1.0,
			amp:        1.0,
			// current implementation: waveformValue := 4*abs(phase-0.5) - 1
			// phases: 0, 0.25, 0.5, 0.75 -> +1, 0, -1, 0
			expect: expectMap{
				0: Max,
				1: 0,
				2: int16(-Max),
				3: 0,
			},
		},
		{
			name:       "square full amplitude",
			waveType:   Square,
			numSamples: 4,
			sampleRate: 4,
			freq:       1.0,
			amp:        1.0,
			expect: expectMap{
				0: Max,
				1: Max,
				2: int16(-Max),
				3: int16(-Max),
			},
		},
		{
			name:       "square half amplitude",
			waveType:   Square,
			numSamples: 4,
			sampleRate: 4,
			freq:       1.0,
			amp:        0.5,
			expect: expectMap{
				0: int16(Max / 2),
				1: int16(Max / 2),
				2: int16(-Max / 2),
				3: int16(-Max / 2),
			},
		},
		{
			name:       "sin sample range smoke",
			waveType:   Sine,
			numSamples: 100,
			sampleRate: 44100,
			freq:       440.0,
			amp:        1.0,
			validator: func(b []byte) error {
				for i := uint(0); i < 100; i++ {
					v := int16At(b, i)
					if v > int16(math.MaxInt16) || v < int16(-math.MaxInt16) {
						return errors.New("sample out of int16 range")
					}
				}
				return nil
			},
		},
	}

	for _, c := range cases {
		c := c // capture
		t.Run(c.name, func(t *testing.T) {
			buf := make([]byte, c.numSamples*2)
			out, err := GenerateWave(c.waveType, c.numSamples, c.sampleRate, c.freq, c.amp, buf)
			if c.wantErr {
				if err == nil {
					t.Fatalf("expected error for case %q but got nil", c.name)
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error for case %q: %v", c.name, err)
			}
			if len(out) != int(c.numSamples*2) {
				t.Fatalf("unexpected output length: got %d want %d", len(out), c.numSamples*2)
			}

			if c.expect != nil {
				for idx, want := range c.expect {
					got := int16At(out, idx)
					if got != want {
						t.Fatalf("case %q: sample[%d] = %d, want %d", c.name, idx, got, want)
					}
				}
			}

			if c.validator != nil {
				if err := c.validator(out); err != nil {
					t.Fatalf("validator failed for case %q: %v", c.name, err)
				}
			}
		})
	}
}
