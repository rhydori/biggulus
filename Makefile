OUTPUT = server.exe
MAIN_PATH = ./cmd/main.go
LINK = .\startserver.lnk

server:
	@@ go run $(MAIN_PATH) -o ../

build:
	@@ go build -o $(OUTPUT) $(MAIN_PATH)

run: build
	@@ explorer.exe "$(LINK)"