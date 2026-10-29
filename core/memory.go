package main
import "device/uefi"
import "unsafe"

const (
	pageSize = 4096
	stackPages = 16
	heapPages = 256
)

var stackTop uintptr
var heap []byte
var heapUsed int

func memoryInit() bool {
	var stackAddress uefi.EFI_PHYSICAL_ADDRESS
	var heapAddress uefi.EFI_PHYSICAL_ADDRESS

	status := uefi.BS().AllocatePages(uefi.AllocateAnyPages, uefi.EfiLoaderData, stackPages, &stackAddress)
	if status != uefi.EFI_SUCCESS {
		return false
	}

	status = uefi.BS().AllocatePages(uefi.AllocateAnyPages, uefi.EfiLoaderData, heapPages, &heapAddress)
	if status != uefi.EFI_SUCCESS {
		uefi.BS().FreePages(stackAddress, stackPages)
		return false
	}

	stackTop = uintptr(stackAddress) + stackPages * pageSize
	heap = unsafe.Slice((*byte)(unsafe.Pointer(uintptr(heapAddress))), heapPages * pageSize)
	return true
}

func heapAlloc(size int) []byte {
    if size < 0 || size > len(heap) - heapUsed {
        return nil
    }

    alignedSize := (size + 15) &^ 15
    if alignedSize > len(heap)-heapUsed {
        return nil
    }

    start := heapUsed
    heapUsed += alignedSize

    return heap[start : start + size]
}
