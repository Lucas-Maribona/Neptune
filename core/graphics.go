package main
import "device/uefi"
import "unsafe"

type GraphicsModeInfo struct {
	Version			  uint32
	Width			  uint32
	Height			  uint32
	PixelFormat		  uint32
	PixelInformation  [4]uint32
	PixelsPerScanLine uint32
}

type GraphicsMode struct {
	MaxMode			uint32
	CurrentMode		uint32
	Information     *GraphicsModeInfo
	InformationSize uint64
	Framebuffer		uint64
	FramebufferSize uint64
}

type GraphicsProtocol struct {
	QueryMode uintptr
	SetMode   uintptr
	Blt		  uintptr
	Mode	  *GraphicsMode
}

type Andromeda struct {
	Pixels		[]uint32
	Width		uint32
	Height		uint32
	Stride		uint32
	PixelFormat uint32
}

var screen Andromeda
func andromedaInit() bool {
	guid := uefi.EFI_GUID{
		Data1: 0x9042a9de,
		Data2: 0x23dc,
		Data3: 0x4a38,
		Data4: [8]byte{0x96, 0xfb, 0x7a, 0xde, 0xd0, 0x80, 0x51, 0x6a},
	}

	var graphics *GraphicsProtocol;
	status := uefi.BS().LocateProtocol(&guid, nil, unsafe.Pointer(&graphics))

	if status != uefi.EFI_SUCCESS || graphics == nil || graphics.Mode == nil {
		return false
	}

	mode := graphics.Mode
	if mode.Information == nil || mode.Framebuffer == 0 || mode.Information.PixelFormat > 1 {
		return false
	}

	screen = Andromeda{
		Pixels:      unsafe.Slice((*uint32)(unsafe.Pointer(uintptr(mode.Framebuffer))), int(mode.FramebufferSize / 4)),
		Width:       mode.Information.Width,
		Height:      mode.Information.Height,
		Stride:      mode.Information.PixelsPerScanLine,
		PixelFormat: mode.Information.PixelFormat,
	}

	return true
}

func andromedaDrawPixel(x, y, color uint32) {
    if x >= screen.Width || y >= screen.Height {
        return
    }

    if screen.PixelFormat == 0 {
        color = (color & 0x00FF00) |
            ((color & 0x0000FF) << 16) |
            ((color & 0xFF0000) >> 16)
    }

    screen.Pixels[y*screen.Stride+x] = color
}

func andromedaDrawRectangle(x, y, width, height, color uint32) {
    for row := uint32(0); row < height; row++ {
        for column := uint32(0); column < width; column++ {
            andromedaDrawPixel(x+column, y+row, color)
        }
    }
}

func andromedaClear(color uint32) {
    andromedaDrawRectangle(0, 0, screen.Width, screen.Height, color)
}
