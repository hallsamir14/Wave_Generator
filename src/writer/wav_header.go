package writer

import (
	"encoding/binary"
	"os"
)

/*
<WAVE-form> → RIFF('WAVE'

	<fmt-ck>            // Format of the file
	[<fact-ck>]         // Fact chunk
	[<cue-ck>]          // Cue points
	[<playlist-ck>]     // Playlist
	[<assoc-data-list>] // Associated data list
	<wave-data> )       // Wave data
*/
func WriteWavHeader(f *os.File, dataSize uint32, sampleRate int, channels int16, bitsPerSample int16) {
	var chunkSize = 36 + dataSize
	var byteRate = sampleRate * int(channels) * int(bitsPerSample) / 8
	var blockAlign = channels * bitsPerSample / 8

	// RIFF header
	f.Write([]byte("RIFF"))
	binary.Write(f, binary.LittleEndian, uint32(chunkSize))
	f.Write([]byte("WAVE"))

	// fmt subchunk
	f.Write([]byte("fmt "))
	binary.Write(f, binary.LittleEndian, uint32(16))            // Subchunk1Size (16 for PCM)
	binary.Write(f, binary.LittleEndian, uint16(1))             // AudioFormat (1 = PCM)
	binary.Write(f, binary.LittleEndian, channels)              // NumChannels
	binary.Write(f, binary.LittleEndian, uint32(sampleRate))    // SampleRate
	binary.Write(f, binary.LittleEndian, uint32(byteRate))      // ByteRate
	binary.Write(f, binary.LittleEndian, uint16(blockAlign))    // BlockAlign
	binary.Write(f, binary.LittleEndian, uint16(bitsPerSample)) // BitsPerSample

	// data subchunk
	f.Write([]byte("data"))
	binary.Write(f, binary.LittleEndian, uint32(dataSize))
}
