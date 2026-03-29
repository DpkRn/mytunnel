# Linux
GOOS=linux GOARCH=amd64 go build -o mytunnel-linux

# Mac
GOOS=darwin GOARCH=amd64 go build -o mytunnel-mac

# Windows
GOOS=windows GOARCH=amd64 go build -o mytunnel.exe