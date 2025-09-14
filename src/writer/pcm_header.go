package writer

import (
	"log"
	"os"
)

// Write PCM data
func WritePCMData(data []byte, file *os.File, err interface{}) {
	_, err = file.Write(data)
	if err != nil {
		log.Fatal(err)
	}

}
