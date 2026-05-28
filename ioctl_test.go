package ioctl

import "testing"

func TestIOC(t *testing.T) {
	tests := []struct {
		name     string
		dir      uintptr
		typ      uintptr
		nr       uintptr
		size     uintptr
		expected uintptr
	}{
		{"none zero", _IOC_NONE, 0, 1, 0, 0x00000001},
		{"read with size", _IOC_READ, 0, 1, 4, 0x80040001},
		{"write with size", _IOC_WRITE, 0, 1, 4, 0x40040001},
		{"dir bits masked", _IOC_READ | _IOC_WRITE, 0, 1, 4, 0xC0040001},
		{"max nr", _IOC_READ, 0, 0xFF, 4, 0x800400FF},
		{"max type", _IOC_READ, 0xFF, 1, 4, 0x8004FF01},
		{"max size", _IOC_READ, 0, 1, 0x3FFF, 0xBFFF0001},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ioc(tt.dir, tt.typ, tt.nr, tt.size)
			if got != tt.expected {
				t.Errorf("ioc(%#x, %#x, %#x, %#x) = %#x, want %#x",
					tt.dir, tt.typ, tt.nr, tt.size, got, tt.expected)
			}
		})
	}
}

func TestMacros(t *testing.T) {
	if got, want := IO('T', 1), uintptr(0x5401); got != want {
		t.Errorf("IO = %#x, want %#x", got, want)
	}
	if got, want := IOR('T', 1, 4), uintptr(0x80045401); got != want {
		t.Errorf("IOR = %#x, want %#x", got, want)
	}
	if got, want := IOW('T', 1, 4), uintptr(0x40045401); got != want {
		t.Errorf("IOW = %#x, want %#x", got, want)
	}
	if got, want := IOWR('T', 1, 4), uintptr(0xC0045401); got != want {
		t.Errorf("IOWR = %#x, want %#x", got, want)
	}
}

func TestDecode(t *testing.T) {
	tests := []struct {
		name string
		cmd  uintptr
		want IOC
	}{
		{"IO", 0x5401, IOC{_IOC_NONE, 'T', 1, 0}},
		{"IOR", 0x80045401, IOC{_IOC_READ, 'T', 1, 4}},
		{"IOW", 0x40045401, IOC{_IOC_WRITE, 'T', 1, 4}},
		{"IOWR", 0xC0045401, IOC{_IOC_READ | _IOC_WRITE, 'T', 1, 4}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Decode(tt.cmd)
			if got != tt.want {
				t.Errorf("Decode(%#x) = %+v, want %+v", tt.cmd, got, tt.want)
			}
		})
	}
}

func TestDecodeRoundtrip(t *testing.T) {
	tests := []struct {
		name string
		dir  uintptr
		typ  uintptr
		nr   uintptr
		size uintptr
	}{
		{"IO", _IOC_NONE, 'K', 1, 0},
		{"IOR", _IOC_READ, 'K', 1, 8},
		{"IOW", _IOC_WRITE, 'K', 1, 16},
		{"IOWR", _IOC_READ | _IOC_WRITE, 'K', 1, 32},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := ioc(tt.dir, tt.typ, tt.nr, tt.size)
			d := Decode(cmd)
			if d.Dir != tt.dir || d.Type != tt.typ || d.NR != tt.nr || d.Size != tt.size {
				t.Errorf("roundtrip failed: encoded %#x, decoded %+v, want dir=%#x type=%#x nr=%#x size=%#x",
					cmd, d, tt.dir, tt.typ, tt.nr, tt.size)
			}
		})
	}
}
