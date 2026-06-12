# Audio Fingerprinting (Shazam-like)

> Shazam-style music detection in Go

---

- [x] **Go module setup** (`go.mod`, `go.sum`, dependencies)
- [x] **WAV decoder** (`decode/decode.go`) — RIFF header parsing, PCM reading, stereo-to-mono, normalization
- [x] **FFT Spectrogram** (`fingerprint/spectogram.go`) — Hann window, FFT via gonum, magnitude computation
- [x] **Spectrogram image export** (`fingerprint/spectogram.go` — `SpectrogramImage`)
- [x] **Sine wave generator** (`fingerprint/spectogram.go` — `GenerateSineWave`)
- [x] **Basic frequency analysis** (`fingerprint/spectogram.go` — `AnalyzeFrequency`)
- [ ] **Peak finding / constellation map** — local maxima in time-frequency grid
- [ ] **Fingerprint hashing** — pair peaks, create hash, store anchor time
- [ ] **Storage (in-memory / JSON)** — song catalog + hash lookup tables
- [ ] **Matching / offset alignment scoring** — group by song ID, find dominant offset
- [ ] **Mic recording** (`portaudio`) — capture audio from microphone
- [ ] **CLI** — `add-song`, `listen`, `list-songs` commands
- [ ] **Polish** — noise filtering, peak density tuning, DB export/import

---

**Legend:** `[x]` = done, `[ ]` = needs work
