package ioctl

const (
	_IOC_NRBITS   = 8
	_IOC_TYPEBITS = 8
	_IOC_SIZEBITS = 14
	_IOC_DIRBITS  = 2

	_IOC_NONE  = 0
	_IOC_WRITE = 1
	_IOC_READ  = 2

	_IOC_NRMASK   = (1 << _IOC_NRBITS) - 1
	_IOC_TYPEMASK = (1 << _IOC_TYPEBITS) - 1
	_IOC_SIZEMASK = (1 << _IOC_SIZEBITS) - 1
	_IOC_DIRMASK  = (1 << _IOC_DIRBITS) - 1
)

func ioc(dir, t, nr, size uintptr) uintptr {
	return (dir << 30) | (size << 16) | (t << 8) | nr
}

func IOR(t, nr, size uintptr) uintptr  { return ioc(_IOC_READ, t, nr, size) }
func IOW(t, nr, size uintptr) uintptr  { return ioc(_IOC_WRITE, t, nr, size) }
func IOWR(t, nr, size uintptr) uintptr { return ioc(_IOC_READ|_IOC_WRITE, t, nr, size) }
func IO(t, nr uintptr) uintptr         { return ioc(_IOC_NONE, t, nr, 0) }

type IOC struct {
	Dir  uintptr
	Type uintptr
	NR   uintptr
	Size uintptr
}

func Decode(cmd uintptr) IOC {
	return IOC{
		Dir:  (cmd >> 30) & _IOC_DIRMASK,
		Size: (cmd >> 16) & _IOC_SIZEMASK,
		Type: (cmd >> 8) & _IOC_TYPEMASK,
		NR:   cmd & _IOC_NRMASK,
	}
}
