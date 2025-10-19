package writer

import (
	"bytes"
	"encoding/binary"
	"io"
	"os"
	"testing"
)

//TODO review suite

func TestWriteWavHeader(t *testing.T) {
	// Parameters to pass to WriteWavHeader
	var dataSize uint32 = 1000
	sampleRate := 44100
	var channels int16 = 1
	var bitsPerSample int16 = 16

	// Create an os.Pipe so we get two *os.File values.
	// We'll pass the writer end to WriteWavHeader and read the header from the reader end.
	r, w, err := os.Pipe()
	if err != nil {
		t.Fatalf("os.Pipe failed: %v", err)
	}
	// Ensure we close both ends
	defer r.Close()
	defer w.Close()

	// Call the function under test. It will write to the writer-end of the pipe.
	WriteWavHeader(w, dataSize, sampleRate, channels, bitsPerSample)

	// Read the header bytes from the reader end. WAV header written by the function
	// is 44 bytes for PCM (RIFF + fmt + data subchunk).
	expectedHeaderLen := 44
	buf := make([]byte, expectedHeaderLen)
	n, err := io.ReadFull(r, buf)
	if err != nil {
		t.Fatalf("failed to read header from pipe: %v (read %d bytes)", err, n)
	}

	// Validate ASCII identifiers
	if !bytes.Equal(buf[0:4], []byte("RIFF")) {
		t.Fatalf("expected 'RIFF' at bytes 0-3, got %q", buf[0:4])
	}
	if !bytes.Equal(buf[8:12], []byte("WAVE")) {
		t.Fatalf("expected 'WAVE' at bytes 8-11, got %q", buf[8:12])
	}
	if !bytes.Equal(buf[12:16], []byte("fmt ")) {
		t.Fatalf("expected 'fmt ' at bytes 12-15, got %q", buf[12:16])
	}
	if !bytes.Equal(buf[36:40], []byte("data")) {
		t.Fatalf("expected 'data' at bytes 36-39, got %q", buf[36:40])
	}

	// Helper to read little-endian values from buffer
	readUint32 := func(b []byte) uint32 {
		return binary.LittleEndian.Uint32(b)
	}
	readUint16 := func(b []byte) uint16 {
		return binary.LittleEndian.Uint16(b)
	}

	// chunkSize at bytes 4..7 should be 36 + dataSize
	gotChunkSize := readUint32(buf[4:8])
	wantChunkSize := uint32(36) + dataSize
	if gotChunkSize != wantChunkSize {
		t.Fatalf("chunkSize mismatch: got %d want %d", gotChunkSize, wantChunkSize)
	}

	// Subchunk1Size at bytes 16..19 (fmt subchunk size) should be 16 for PCM
	gotSubchunk1Size := readUint32(buf[16:20])
	if gotSubchunk1Size != 16 {
		t.Fatalf("fmt subchunk1 size mismatch: got %d want 16", gotSubchunk1Size)
	}

	// AudioFormat at bytes 20..21 (should be 1 for PCM)
	gotAudioFormat := readUint16(buf[20:22])
	if gotAudioFormat != 1 {
		t.Fatalf("audio format mismatch: got %d want 1 (PCM)", gotAudioFormat)
	}

	// NumChannels at bytes 22..23
	gotChannels := int16(readUint16(buf[22:24]))
	if gotChannels != channels {
		t.Fatalf("channels mismatch: got %d want %d", gotChannels, channels)
	}

	// SampleRate at bytes 24..27
	gotSampleRate := int(readUint32(buf[24:28]))
	if gotSampleRate != sampleRate {
		t.Fatalf("sampleRate mismatch: got %d want %d", gotSampleRate, sampleRate)
	}

	// ByteRate at bytes 28..31 should be sampleRate * channels * bitsPerSample / 8
	gotByteRate := readUint32(buf[28:32])
	wantByteRate := uint32(sampleRate * int(channels) * int(bitsPerSample) / 8)
	if gotByteRate != wantByteRate {
		t.Fatalf("byteRate mismatch: got %d want %d", gotByteRate, wantByteRate)
	}

	// BlockAlign at bytes 32..33 should be channels * bitsPerSample / 8
	gotBlockAlign := readUint16(buf[32:34])
	wantBlockAlign := uint16(channels * bitsPerSample / 8)
	if gotBlockAlign != wantBlockAlign {
		t.Fatalf("blockAlign mismatch: got %d want %d", gotBlockAlign, wantBlockAlign)
	}

	// BitsPerSample at bytes 34..35
	gotBitsPerSample := readUint16(buf[34:36])
	if int16(gotBitsPerSample) != bitsPerSample {
		t.Fatalf("bitsPerSample mismatch: got %d want %d", gotBitsPerSample, bitsPerSample)
	}

	// data subchunk size at bytes 40..43 should equal dataSize
	gotDataSize := readUint32(buf[40:44])
	if gotDataSize != dataSize {
		t.Fatalf("dataSize mismatch: got %d want %d", gotDataSize, dataSize)
	}
}
