.DEFAULT_GOAL := all
DistDir = dist

$(DistDir):
	mkdir -p $(DistDir)

$(DistDir)/BusyBee_linux: $(DistDir)
	env GOOS=linux GOARCH=amd64 go build -a -o $@ ./cmd/BusyBee

$(DistDir)/BusyBee.exe: $(DistDir)
	env GOOS=windows GOARCH=amd64 go build -a -o $@ ./cmd/BusyBee

$(DistDir)/BusyBee_mac: $(DistDir)
	env GOOS=darwin GOARCH=amd64 go build -a -o $@ ./cmd/BusyBee

$(DistDir)/WindowsInstall.ps1: $(DistDir)
	cp WindowsInstall.ps1 $(DistDir)/.

$(DistDir)/UnixInstall.sh: $(DistDir)
	cp UnixInstall.sh $(DistDir)/.

.PHONY: all
all: $(DistDir)/WindowsInstall.ps1 $(DistDir)/UnixInstall.sh $(DistDir)/BusyBee_mac $(DistDir)/BusyBee.exe $(DistDir)/BusyBee_linux

.PHONY: clean
clean:
	rm -rf $(DistDir)
	go clean -r