package main
import "device/uefi"

func main() {
    if !memoryInit() {
        println("[Neptune ERROR]: Memory initialization failed!")
        return
    }

    if !andromedaInit() {
        println("[Neptune ERROR]: Graphics initialization failed!")
        return
    }

    memory := heapAlloc(128)
    if memory == nil {
        println("[Neptune ERROR]: Memory allocation failed!")
        return
    }

    memory[0] = 42
    andromedaClear(0x000000)
    andromedaDrawRectangle(100, 100, 300, 200, 0x2563EB)
    println("[Neptune INFO]: Neptune memory initialized!")

    for {
        uefi.CpuPause()
    }
}
