package decode_test

import (
	"encoding/binary"
	"math"
	"os"
	"testing"

	"github.com/Ashking-tech/audio/decode"
)

func makeWAV(channels, sampleRate, bitsPerSample, numSamples int) []byte {
	bytesPerSample := bitsPerSample / 8
	dataSize := numSamples * channels * bytesPerSample
	headerSize := 44
	buf := make([]byte, headerSize+dataSize)

	copy(buf[0:4], "RIFF")
	binary.LittleEndian.PutUint32(buf[4:8], uint32(36+dataSize))
	copy(buf[8:12], "WAVE")
	copy(buf[12:16], "fmt ")
	binary.LittleEndian.PutUint32(buf[16:20], 16)
	binary.LittleEndian.PutUint16(buf[20:22], 1)
	binary.LittleEndian.PutUint16(buf[22:24], uint16(channels))
	binary.LittleEndian.PutUint32(buf[24:28], uint32(sampleRate))
	binary.LittleEndian.PutUint32(buf[28:32], uint32(sampleRate*channels*bytesPerSample))
	binary.LittleEndian.PutUint16(buf[32:34], uint16(channels*bytesPerSample))
	binary.LittleEndian.PutUint16(buf[34:36], uint16(bitsPerSample))
	copy(buf[36:40], "data")
	binary.LittleEndian.PutUint32(buf[40:44], uint32(dataSize))

	for i := 0; i < numSamples*channels; i++ {
		val := int16(math.Sin(float64(i)/100) * 10000)
		binary.LittleEndian.PutUint16(buf[headerSize+i*2:], uint16(val))
	}
	return buf
}

func writeTempWAV(t *testing.T, data []byte) string {
	t.Helper()
	f, err := os.CreateTemp(t.TempDir(), "*.wav")
	if err != nil {
		t.Fatal(err)
	}
	if _, err := f.Write(data); err != nil {
		f.Close()
		t.Fatal(err)
	}
	f.Close()
	return f.Name()
}

func buildWAVWithExtraChunks() []byte {
	b := make([]byte, 0, 300)
	write := func(data ...byte) { b = append(b, data...) }
	writeStr := func(s string) { write([]byte(s)...) }
	writeU16 := func(v uint16) { p := make([]byte, 2); binary.LittleEndian.PutUint16(p, v); write(p...) }
	writeU32 := func(v uint32) { p := make([]byte, 4); binary.LittleEndian.PutUint32(p, v); write(p...) }

	writeStr("RIFF")
	writeU32(0)
	writeStr("WAVE")
	writeStr("fmt ")
	writeU32(16)
	writeU16(1)
	writeU16(1)
	writeU32(44100)
	writeU32(88200)
	writeU16(2)
	writeU16(16)
	writeStr("fact")
	writeU32(4)
	writeU32(0)
	writeStr("data")
	writeU32(200)
	for i := 0; i < 100; i++ {
		writeU16(uint16(i))
	}
	binary.LittleEndian.PutUint32(b[4:8], uint32(len(b)-8))
	return b
}

func TestParseHeader(t *testing.T) {
	data := makeWAV(1, 44100, 16, 1000)
	m, err := decode.ParseHeader(data)
	if err != nil {
		t.Fatal(err)
	}
	if m.Channels != 1 {
		t.Errorf("Channels = %d, want 1", m.Channels)
	}
	if m.SampleRate != 44100 {
		t.Errorf("SampleRate = %d, want 44100", m.SampleRate)
	}
	if m.BitsPerSample != 16 {
		t.Errorf("BitsPerSample = %d, want 16", m.BitsPerSample)
	}
	if m.DataOffset != 44 {
		t.Errorf("DataOffset = %d, want 44", m.DataOffset)
	}
}

func TestParseHeader_extraChunks(t *testing.T) {
	b := buildWAVWithExtraChunks()
	m, err := decode.ParseHeader(b)
	if err != nil {
		t.Fatal(err)
	}
	if m.DataOffset != 56 {
		t.Errorf("DataOffset = %d, want 56", m.DataOffset)
	}
}

func TestParseHeader_invalid(t *testing.T) {
	_, err := decode.ParseHeader([]byte("short"))
	if err == nil {
		t.Fatal("expected error for short input")
	}
}

func TestParseHeader_notRIFF(t *testing.T) {
	buf := make([]byte, 44)
	copy(buf[0:4], "NOT ")
	copy(buf[8:12], "WAVE")
	_, err := decode.ParseHeader(buf)
	if err == nil {
		t.Fatal("expected error for non-RIFF")
	}
}

func TestReadPCMSamples(t *testing.T) {
	data := makeWAV(1, 44100, 16, 100)
	m, _ := decode.ParseHeader(data)
	samples, err := decode.ReadPCMSamples(data, m)
	if err != nil {
		t.Fatal(err)
	}
	if len(samples) != 100 {
		t.Errorf("got %d samples, want 100", len(samples))
	}
}

func TestReadPCMSamples_unsupportedBits(t *testing.T) {
	m := decode.Metadata{BitsPerSample: 24}
	_, err := decode.ReadPCMSamples(nil, m)
	if err == nil {
		t.Fatal("expected error for 24-bit")
	}
}

func TestStereoToMono(t *testing.T) {
	samples := []int16{100, 200, 300, 400}
	mono := decode.StereoToMono(samples, 2)
	if len(mono) != 2 {
		t.Fatalf("got %d samples, want 2", len(mono))
	}
	if mono[0] != 150 {
		t.Errorf("mono[0] = %d, want 150", mono[0])
	}
	if mono[1] != 350 {
		t.Errorf("mono[1] = %d, want 350", mono[1])
	}
}

func TestStereoToMono_alreadyMono(t *testing.T) {
	samples := []int16{100, 200, 300}
	mono := decode.StereoToMono(samples, 1)
	if len(mono) != 3 {
		t.Fatalf("got %d samples, want 3", len(mono))
	}
	if mono[0] != 100 {
		t.Errorf("mono[0] = %d, want 100", mono[0])
	}
}

func TestNormalize(t *testing.T) {
	samples := []int16{0, 16384, 32767, -32768}
	norm := decode.Normalize(samples)
	if len(norm) != 4 {
		t.Fatalf("got %d samples, want 4", len(norm))
	}
	if norm[0] != 0.0 {
		t.Errorf("norm[0] = %f, want 0", norm[0])
	}
	if norm[1] != 16384.0/32768.0 {
		t.Errorf("norm[1] = %f, want %f", norm[1], 16384.0/32768.0)
	}
	if norm[2] != 32767.0/32768.0 {
		t.Errorf("norm[2] = %f, want %f", norm[2], 32767.0/32768.0)
	}
	if norm[3] != -1.0 {
		t.Errorf("norm[3] = %f, want -1", norm[3])
	}
}

func TestDecodeWav(t *testing.T) {
	data := makeWAV(2, 48000, 16, 30000)
	path := writeTempWAV(t, data)
	samples, err := decode.DecodeWav(path)
	if err != nil {
		t.Fatal(err)
	}
	if len(samples) == 0 {
		t.Fatal("got 0 samples")
	}
	for i, s := range samples {
		if s < -1.0 || s > 1.0 {
			t.Fatalf("sample %d = %f out of range", i, s)
		}
	}
}

func TestDecodeWav_invalidFile(t *testing.T) {
	_, err := decode.DecodeWav("/nonexistent/file.wav")
	if err == nil {
		t.Fatal("expected error for nonexistent file")
	}
}

func TestReadAudioFile(t *testing.T) {
	data := []byte{1, 2, 3}
	f := writeTempWAV(t, data)
	read, err := decode.ReadAudioFile(f)
	if err != nil {
		t.Fatal(err)
	}
	if len(read) == 0 {
		t.Fatal("read 0 bytes")
	}
}
