package input

type Note struct {
	Midi      int64
	Label     string
	Frequency float32
}

var MidiNotes = map[int64]Note{
	127: Note{Midi: 127, Label: "G9", Frequency: 12543.85},
	126: Note{Midi: 126, Label: "F#9/Gb9", Frequency: 11839.82},
	125: Note{Midi: 125, Label: "F9", Frequency: 11175.30},
	124: Note{Midi: 124, Label: "E9", Frequency: 10548.08},
	123: Note{Midi: 123, Label: "D#9/Eb9", Frequency: 9956.06},
	122: Note{Midi: 122, Label: "D9", Frequency: 9397.27},
	121: Note{Midi: 121, Label: "C#9/Db9", Frequency: 8869.84},
	120: Note{Midi: 120, Label: "C9", Frequency: 8372.02},
	119: Note{Midi: 119, Label: "B8", Frequency: 7902.13},
	118: Note{Midi: 118, Label: "A#8/Bb8", Frequency: 7458.62},
	117: Note{Midi: 117, Label: "A8", Frequency: 7040.00},
	116: Note{Midi: 116, Label: "G#8/Ab8", Frequency: 6644.88},
	115: Note{Midi: 115, Label: "G8", Frequency: 6271.93},
	114: Note{Midi: 114, Label: "F#8/Gb8", Frequency: 5919.91},
	113: Note{Midi: 113, Label: "F8", Frequency: 5587.65},
	112: Note{Midi: 112, Label: "E8", Frequency: 5274.04},
	111: Note{Midi: 111, Label: "D#8/Eb8", Frequency: 4978.03},
	110: Note{Midi: 110, Label: "D8", Frequency: 4698.64},
	109: Note{Midi: 109, Label: "C#8/Db8", Frequency: 4434.92},
	108: Note{Midi: 108, Label: "C8", Frequency: 4186.01},
	107: Note{Midi: 107, Label: "B7", Frequency: 3951.07},
	106: Note{Midi: 106, Label: "A#7/Bb7", Frequency: 3729.31},
	105: Note{Midi: 105, Label: "A7", Frequency: 3520.00},
	104: Note{Midi: 104, Label: "G#7/Ab7", Frequency: 3322.44},
	103: Note{Midi: 103, Label: "G7", Frequency: 3135.96},
	102: Note{Midi: 102, Label: "F#7/Gb7", Frequency: 2959.96},
	101: Note{Midi: 101, Label: "F7", Frequency: 2793.83},
	100: Note{Midi: 100, Label: "E7", Frequency: 2637.02},
	99:  Note{Midi: 99, Label: "D#7/Eb7", Frequency: 2489.02},
	98:  Note{Midi: 98, Label: "D7", Frequency: 2349.32},
	97:  Note{Midi: 97, Label: "C#7/Db7", Frequency: 2217.46},
	96:  Note{Midi: 96, Label: "C7", Frequency: 2093.00},
	95:  Note{Midi: 95, Label: "B6", Frequency: 1975.53},
	94:  Note{Midi: 94, Label: "A#6/Bb6", Frequency: 1864.66},
	93:  Note{Midi: 93, Label: "A6", Frequency: 1760.00},
	92:  Note{Midi: 92, Label: "G#6/Ab6", Frequency: 1661.22},
	91:  Note{Midi: 91, Label: "G6", Frequency: 1567.98},
	90:  Note{Midi: 90, Label: "F#6/Gb6", Frequency: 1479.98},
	89:  Note{Midi: 89, Label: "F6", Frequency: 1396.91},
	88:  Note{Midi: 88, Label: "E6", Frequency: 1318.51},
	87:  Note{Midi: 87, Label: "D#6/Eb6", Frequency: 1244.51},
	86:  Note{Midi: 86, Label: "D6", Frequency: 1174.66},
	85:  Note{Midi: 85, Label: "C#6/Db6", Frequency: 1108.73},
	84:  Note{Midi: 84, Label: "C6", Frequency: 1046.50},
	83:  Note{Midi: 83, Label: "B5", Frequency: 987.77},
	82:  Note{Midi: 82, Label: "A#5/Bb5", Frequency: 932.33},
	81:  Note{Midi: 81, Label: "A5", Frequency: 880.00},
	80:  Note{Midi: 80, Label: "G#5/Ab5", Frequency: 830.61},
	79:  Note{Midi: 79, Label: "G5", Frequency: 783.99},
	78:  Note{Midi: 78, Label: "F#5/Gb5", Frequency: 739.99},
	77:  Note{Midi: 77, Label: "F5", Frequency: 698.46},
	76:  Note{Midi: 76, Label: "E5", Frequency: 659.26},
	75:  Note{Midi: 75, Label: "D#5/Eb5", Frequency: 622.25},
	74:  Note{Midi: 74, Label: "D5", Frequency: 587.33},
	73:  Note{Midi: 73, Label: "C#5/Db5", Frequency: 554.37},
	72:  Note{Midi: 72, Label: "C5", Frequency: 523.25},
	71:  Note{Midi: 71, Label: "B4", Frequency: 493.88},
	70:  Note{Midi: 70, Label: "A#4/Bb4", Frequency: 466.16},
	69:  Note{Midi: 69, Label: "A4 concert pitch", Frequency: 440.00},
	68:  Note{Midi: 68, Label: "G#4/Ab4", Frequency: 415.30},
	67:  Note{Midi: 67, Label: "G4", Frequency: 392.00},
	66:  Note{Midi: 66, Label: "F#4/Gb4", Frequency: 369.99},
	65:  Note{Midi: 65, Label: "F4", Frequency: 349.23},
	64:  Note{Midi: 64, Label: "E4", Frequency: 329.63},
	63:  Note{Midi: 63, Label: "D#4/Eb4", Frequency: 311.13},
	62:  Note{Midi: 62, Label: "D4", Frequency: 293.66},
	61:  Note{Midi: 61, Label: "C#4/Db4", Frequency: 277.18},
	60:  Note{Midi: 60, Label: "C4 (middle C)", Frequency: 261.63},
	59:  Note{Midi: 59, Label: "B3", Frequency: 246.94},
	58:  Note{Midi: 58, Label: "A#3/Bb3", Frequency: 233.08},
	57:  Note{Midi: 57, Label: "A3", Frequency: 220.00},
	56:  Note{Midi: 56, Label: "G#3/Ab3", Frequency: 207.65},
	55:  Note{Midi: 55, Label: "G3", Frequency: 196.00},
	54:  Note{Midi: 54, Label: "F#3/Gb3", Frequency: 185.00},
	53:  Note{Midi: 53, Label: "F3", Frequency: 174.61},
	52:  Note{Midi: 52, Label: "E3", Frequency: 164.81},
	51:  Note{Midi: 51, Label: "D#3/Eb3", Frequency: 155.56},
	50:  Note{Midi: 50, Label: "D3", Frequency: 146.83},
	49:  Note{Midi: 49, Label: "C#3/Db3", Frequency: 138.59},
	48:  Note{Midi: 48, Label: "C3", Frequency: 130.81},
	47:  Note{Midi: 47, Label: "B2", Frequency: 123.47},
	46:  Note{Midi: 46, Label: "A#2/Bb2", Frequency: 116.54},
	45:  Note{Midi: 45, Label: "A2", Frequency: 110.00},
	44:  Note{Midi: 44, Label: "G#2/Ab2", Frequency: 103.83},
	43:  Note{Midi: 43, Label: "G2", Frequency: 98.00},
	42:  Note{Midi: 42, Label: "F#2/Gb2", Frequency: 92.50},
	41:  Note{Midi: 41, Label: "F2", Frequency: 87.31},
	40:  Note{Midi: 40, Label: "E2", Frequency: 82.41},
	39:  Note{Midi: 39, Label: "D#2/Eb2", Frequency: 77.78},
	38:  Note{Midi: 38, Label: "D2", Frequency: 73.42},
	37:  Note{Midi: 37, Label: "C#2/Db2", Frequency: 69.30},
	36:  Note{Midi: 36, Label: "C2", Frequency: 65.41},
	35:  Note{Midi: 35, Label: "B1", Frequency: 61.74},
	34:  Note{Midi: 34, Label: "A#1/Bb1", Frequency: 58.27},
	33:  Note{Midi: 33, Label: "A1", Frequency: 55.00},
	32:  Note{Midi: 32, Label: "G#1/Ab1", Frequency: 51.91},
	31:  Note{Midi: 31, Label: "G1", Frequency: 49.00},
	30:  Note{Midi: 30, Label: "F#1/Gb1", Frequency: 46.25},
	29:  Note{Midi: 29, Label: "F1", Frequency: 43.65},
	28:  Note{Midi: 28, Label: "E1", Frequency: 41.20},
	27:  Note{Midi: 27, Label: "D#1/Eb1", Frequency: 38.89},
	26:  Note{Midi: 26, Label: "D1", Frequency: 36.71},
	25:  Note{Midi: 25, Label: "C#1/Db1", Frequency: 34.65},
	24:  Note{Midi: 24, Label: "C1", Frequency: 32.70},
	23:  Note{Midi: 23, Label: "B0", Frequency: 30.87},
	22:  Note{Midi: 22, Label: "A#0/Bb0", Frequency: 29.14},
	21:  Note{Midi: 21, Label: "A0", Frequency: 27.50},
	20:  Note{Midi: 20, Label: " ", Frequency: 25.96},
	19:  Note{Midi: 19, Label: " ", Frequency: 24.50},
	18:  Note{Midi: 18, Label: " ", Frequency: 23.12},
	17:  Note{Midi: 17, Label: " ", Frequency: 21.83},
	16:  Note{Midi: 16, Label: " ", Frequency: 20.60},
	15:  Note{Midi: 15, Label: " ", Frequency: 19.45},
	14:  Note{Midi: 14, Label: " ", Frequency: 18.35},
	13:  Note{Midi: 13, Label: " ", Frequency: 17.32},
	12:  Note{Midi: 12, Label: " ", Frequency: 16.35},
	11:  Note{Midi: 11, Label: " ", Frequency: 15.43},
	10:  Note{Midi: 10, Label: " ", Frequency: 14.57},
	9:   Note{Midi: 9, Label: " ", Frequency: 13.75},
	8:   Note{Midi: 8, Label: " ", Frequency: 12.98},
	7:   Note{Midi: 7, Label: " ", Frequency: 12.25},
	6:   Note{Midi: 6, Label: " ", Frequency: 11.56},
	5:   Note{Midi: 5, Label: " ", Frequency: 10.91},
	4:   Note{Midi: 4, Label: " ", Frequency: 10.30},
	3:   Note{Midi: 3, Label: " ", Frequency: 9.72},
	2:   Note{Midi: 2, Label: " ", Frequency: 9.18},
	1:   Note{Midi: 1, Label: " ", Frequency: 8.66},
	0:   Note{Midi: 0, Label: " ", Frequency: 8.18},
}