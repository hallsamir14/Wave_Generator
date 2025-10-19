# Wave Generator

This project is a simple Go program that generates a sine wave audio file in WAV format. It demonstrates how to synthesize audio data, write PCM samples, and construct a valid WAV file header from scratch.

## Features
- Generates a mono sine wave at a specified frequency, amplitude, and duration
- Outputs a standard 16-bit PCM WAV file (`tone.wav`)
- Modular code structure with separate packages for wave generation and WAV/PCM writing

## Usage

### Prerequisites
- Go 1.18 or later

### Build and Run

1. Clone the repository:
   ```bash
   git clone <repo-url>
   cd wave_generator/src
   ```
2. Build the program:
   ```bash
   go build -o wavegen main.go
   ```
3. Run the program:
   ```bash
   ./wavegen
   ```
   This will generate a `tone.wav` file in the current directory.

## Customization

You can change the following parameters in `main.go` to generate different tones:
- `sampleRate`: Audio sample rate (e.g., 44100 Hz)
- `frequency`: Frequency of the sine wave (e.g., 110.0 Hz)
- `amplitude`: Amplitude of the wave (0.0 to 1.0)
- `durationSec`: Duration of the tone in seconds

## License
MIT License
