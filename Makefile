NAME := aezeed-recover
OUTPUT_DIR := ./bin

.PHONY: build-all
build-all: build-windows-64 build-darwin-arm64 build-darwin-amd64 build-linux-amd build-linux-386 build-linux-arm

.PHONY: checksum-all
checksum-all: checksum-windows-64 checksum-darwin-arm64 checksum-darwin-amd64 checksum-linux-amd checksum-linux-386 checksum-linux-arm

.PHONY: build-windows-64
build-windows-64:
	@echo "Building for Windows 64-bit..."
	@GOOS=windows GOARCH=amd64 go build -o $(OUTPUT_DIR)/$(NAME)_$${VERSION}_windows_64.exe ./main.go

.PHONY: build-darwin-arm64
build-darwin-arm64:
	@echo "Building for macOS arm64..."
	@GOOS=darwin GOARCH=arm64 go build -o $(OUTPUT_DIR)/$(NAME)_$${VERSION}_darwin_arm64 ./main.go

.PHONY: build-darwin-amd64
build-darwin-amd64:
	@echo "Building for macOS x86-64..."
	@GOOS=darwin GOARCH=amd64 go build -o $(OUTPUT_DIR)/$(NAME)_$${VERSION}_darwin_amd64 ./main.go

.PHONY: build-linux-amd
build-linux-amd:
	@echo "Building for Linux amd..."
	@GOOS=linux GOARCH=amd64 go build -o $(OUTPUT_DIR)/$(NAME)_$${VERSION}_linux_amd ./main.go

.PHONY: build-linux-386
build-linux-386:
	@echo "Building for Linux x86..."
	@GOOS=linux GOARCH=386 go build -o $(OUTPUT_DIR)/$(NAME)_$${VERSION}_linux_386 ./main.go

.PHONY: build-linux-arm
build-linux-arm:
	@echo "Building for Linux arm..."
	@GOOS=linux GOARCH=arm go build -o $(OUTPUT_DIR)/$(NAME)_$${VERSION}_linux_arm ./main.go

.PHONY: checksum-windows-64
checksum-windows-64:
	@echo "Generating checksum for Windows 64-bit..."
	@shasum -a 256 $(OUTPUT_DIR)/$(NAME)_$${VERSION}_windows_64.exe > $(OUTPUT_DIR)/$(NAME)_$${VERSION}_windows_64.exe.sha256

.PHONY: checksum-darwin-arm64
checksum-darwin-arm64:
	@echo "Generating checksum for macOS arm64..."
	@shasum -a 256 $(OUTPUT_DIR)/$(NAME)_$${VERSION}_darwin_arm64 > $(OUTPUT_DIR)/$(NAME)_$${VERSION}_darwin_arm64.sha256

.PHONY: checksum-darwin-amd64
checksum-darwin-amd64:
	@echo "Generating checksum for macOS x86-64..."
	@shasum -a 256 $(OUTPUT_DIR)/$(NAME)_$${VERSION}_darwin_amd64 > $(OUTPUT_DIR)/$(NAME)_$${VERSION}_darwin_amd64.sha256

.PHONY: checksum-linux-amd
checksum-linux-amd:
	@echo "Generating checksum for Linux amd..."
	@shasum -a 256 $(OUTPUT_DIR)/$(NAME)_$${VERSION}_linux_amd > $(OUTPUT_DIR)/$(NAME)_$${VERSION}_linux_amd.sha256

.PHONY: checksum-linux-386
checksum-linux-386:
	@echo "Generating checksum for Linux x86..."
	@shasum -a 256 $(OUTPUT_DIR)/$(NAME)_$${VERSION}_linux_386 > $(OUTPUT_DIR)/$(NAME)_$${VERSION}_linux_386.sha256

.PHONY: checksum-linux-arm
checksum-linux-arm:
	@echo "Generating checksum for Linux arm..."
	@shasum -a 256 $(OUTPUT_DIR)/$(NAME)_$${VERSION}_linux_arm > $(OUTPUT_DIR)/$(NAME)_$${VERSION}_linux_arm.sha256

.PHONY: clean
clean:
	@echo "Cleaning build files..."
	@find $(OUTPUT_DIR) -type f ! -name ".gitkeep" -delete

record-demo:
	terminalizer record demo -c ./terminalizer.yml -d "go run main.go --consent --dev"

render-demo:
	terminalizer render demo  -q 25 -o doc/demo.gif
