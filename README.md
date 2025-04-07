# Ballerina Update Tool
Ballerina Update Tool client implementation to manage Ballerina versions

## Build
Use the following commands to build for different platforms and architectures.

For macOS (AMD64)
```bash
GOOS=darwin GOARCH=amd64 go build -o ballerina-command-darwin-amd64 main.go
```

For macOS (ARM64 - M1/M2)
```bash
GOOS=darwin GOARCH=arm64 go build -o ballerina-command-darwin-arm64 main.go
```

For Windows (64-bit)
```bash
GOOS=windows GOARCH=amd64 go build -o ballerina-command-windows-amd64.exe main.go
```

For Windows (32-bit)
```bash
GOOS=windows GOARCH=386 go build -o ballerina-command-windows-386.exe main.go
```

For Linux (64-bit)
```bash
GOOS=linux GOARCH=amd64 go build -o ballerina-command-linux-amd64 main.go
```

For Linux (ARM64 - e.g., Raspberry Pi)
```bash
GOOS=linux GOARCH=arm64 go build -o ballerina-command-linux-arm64 main.go
```

For Linux (ARM - older Raspberry Pi)
```bash
GOOS=linux GOARCH=arm go build -o ballerina-command-linux-arm main.go
```
