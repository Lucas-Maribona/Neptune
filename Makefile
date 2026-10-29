all:
	@clear
	@rm -rf BUILD
	@mkdir -p BUILD/EFI/BOOT

	@tinygo build -target=uefi-amd64 -o BUILD/EFI/BOOT/BOOTX64.EFI ./core

	@cp /usr/share/edk2/x64/OVMF_VARS.4m.fd BUILD/OVMF_VARS.fd
	@qemu-system-x86_64 -machine q35 -m 512M \
		-drive if=pflash,format=raw,readonly=on,file=/usr/share/edk2/x64/OVMF_CODE.4m.fd \
		-drive if=pflash,format=raw,file=BUILD/OVMF_VARS.fd \
		-drive format=raw,file=fat:rw:BUILD

clean:
	@clear
	@rm -rf BUILD
