windows_bin:
	GOOS=windows GOARCH=amd64 go build -o bin/app.exe main.go

linux_bin:
	GOOS=linux GOARCH=amd64 go build -o bin/app.exe main.go