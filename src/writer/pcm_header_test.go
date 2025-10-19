package writer

import (
	"os"
	"testing"
)

func TestWritePCMData(t *testing.T) {
	// Prepare sample data
	mock_data := []byte{0x01, 0x02, 0x03, 0x04}

	// Create a temporary file
	tmpFile, err := os.CreateTemp("", "pcmdata-*.wav")
	if err != nil {
		t.Fatalf("failed to create temp file: %v", err)
	}
	// Ensure cleanup
	fname := tmpFile.Name()
	defer os.Remove(fname)
	tmpFile.Close()

	// Re-open for writing
	f, err := os.OpenFile(fname, os.O_WRONLY|os.O_TRUNC, 0666)
	if err != nil {
		t.Fatalf("failed to open temp file for write: %v", err)
	}
	defer f.Close()

	/*Call the function under test. Pass nil as the err parameter since
	the function currently accepts an interface{} for err.
	*/
	WritePCMData(mock_data, f, nil)

	// Read back and compare
	got, err := os.ReadFile(fname)
	if err != nil {
		t.Fatalf("failed to read temp file: %v", err)
	}

	if len(got) != len(mock_data) {
		t.Fatalf("written length mismatch: got %d, want %d", len(got), len(mock_data))
	}
	for i := range got {
		if got[i] != mock_data[i] {
			t.Fatalf("byte mismatch at %d: got 0x%02x, want 0x%02x", i, got[i], mock_data[i])
		}
	}
}
