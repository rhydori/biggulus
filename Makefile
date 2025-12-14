OUTPUT = server.exe
MAIN_PATH = ./cmd/main.go
LINK = .\startserver.lnk

server:
	@@ cd cmd/ && go run .

build:
	@@ go build -o $(OUTPUT) $(MAIN_PATH)
	@@ explorer.exe "$(LINK)"