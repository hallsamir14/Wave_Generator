package writer

import (
	"encoding/binary"
	"io"
	"os"
	"testing"
)

func TestWriteWavHeader_WithPipe(t *testing.T) {
	type tc struct {
		name          string
		dataSize      uint32
		sampleRate    int
		channels      int16
		bitsPerSample int16
	}

	cases := []tc{
		{
			name:          "stereo 16-bit",
			dataSize:      1234,
			sampleRate:    44100,
			channels:      2,
			bitsPerSample: 16,
		},
		{
			name:          "mono 8-bit",
			dataSize:      42,
			sampleRate:    8000,
			channels:      1,
			bitsPerSample: 8,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			r, w, err := os.Pipe()
			if err != nil {
				t.Fatalf("os.Pipe() error: %v", err)
			}

			/*Write the header into the writer end of the pipe.
			Using an os.Pipe provides us a *os.File (writer) to inject,
			while allowing us to read the bytes off the reader for assertions.
			*/
			WriteWavHeader(w, tc.dataSize, tc.sampleRate, tc.channels, tc.bitsPerSample)

			// close writer so ReadAll will get EOF after the written bytes
			if err := w.Close(); err != nil {
				t.Fatalf("closing writer: %v", err)
			}

			// Read produced bytes
			got, err := io.ReadAll(r)
			if err != nil {
				t.Fatalf("reading from pipe: %v", err)
			}
			if err := r.Close(); err != nil {
				t.Fatalf("closing reader: %v", err)
			}

			if len(got) < 44 {
				t.Fatalf("header too short: got %d bytes, want >= 44", len(got))
			}

			// Validate fixed ASCII fields
			if string(got[0:4]) != "RIFF" {
				t.Fatalf("missing RIFF, got %q", string(got[0:4]))
			}
			if string(got[8:12]) != "WAVE" {
				t.Fatalf("missing WAVE, got %q", string(got[8:12]))
			}
			if string(got[12:16]) != "fmt " {
				t.Fatalf("missing fmt, got %q", string(got[12:16]))
			}
			if string(got[36:40]) != "data" {
				t.Fatalf("missing data, got %q", string(got[36:40]))
			}

			// Parse numeric fields (little-endian)
			chunkSize := binary.LittleEndian.Uint32(got[4:8])
			subchunk1Size := binary.LittleEndian.Uint32(got[16:20])
			audioFormat := binary.LittleEndian.Uint16(got[20:22])
			numChannels := binary.LittleEndian.Uint16(got[22:24])
			sampleRate := binary.LittleEndian.Uint32(got[24:28])
			byteRate := binary.LittleEndian.Uint32(got[28:32])
			blockAlign := binary.LittleEndian.Uint16(got[32:34])
			bitsPerSample := binary.LittleEndian.Uint16(got[34:36])
			dataSize := binary.LittleEndian.Uint32(got[40:44])

			// Expected values
			expectedChunk := uint32(36) + tc.dataSize
			expectedSubchunk1Size := uint32(16)
			expectedAudioFormat := uint16(1)
			expectedNumChannels := uint16(tc.channels)
			expectedSampleRate := uint32(tc.sampleRate)
			expectedByteRate := uint32(tc.sampleRate * int(tc.channels) * int(tc.bitsPerSample) / 8)
			expectedBlockAlign := uint16(tc.channels * tc.bitsPerSample / 8)
			expectedBitsPerSample := uint16(tc.bitsPerSample)
			expectedDataSize := tc.dataSize

			if chunkSize != expectedChunk {
				t.Fatalf("chunkSize: got %d want %d", chunkSize, expectedChunk)
			}
			if subchunk1Size != expectedSubchunk1Size {
				t.Fatalf("subchunk1Size: got %d want %d", subchunk1Size, expectedSubchunk1Size)
			}
			if audioFormat != expectedAudioFormat {
				t.Fatalf("audioFormat: got %d want %d", audioFormat, expectedAudioFormat)
			}
			if numChannels != expectedNumChannels {
				t.Fatalf("numChannels: got %d want %d", numChannels, expectedNumChannels)
			}
			if sampleRate != expectedSampleRate {
				t.Fatalf("sampleRate: got %d want %d", sampleRate, expectedSampleRate)
			}
			if byteRate != expectedByteRate {
				t.Fatalf("byteRate: got %d want %d", byteRate, expectedByteRate)
			}
			if blockAlign != expectedBlockAlign {
				t.Fatalf("blockAlign: got %d want %d", blockAlign, expectedBlockAlign)
			}
			if bitsPerSample != expectedBitsPerSample {
				t.Fatalf("bitsPerSample: got %d want %d", bitsPerSample, expectedBitsPerSample)
			}
			if dataSize != expectedDataSize {
				t.Fatalf("dataSize: got %d want %d", dataSize, expectedDataSize)
			}
		})
	}
}
