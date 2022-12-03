run:
	go run cmd/weather-report/main.go

build:
	go build -o bin/main cmd/weather-report/main.go

compile:
	echo "Compiling for every OS and Platform"
	GOOS=freebsd GOARCH=386 go build -o bin/main-freebsd-386 cmd/weather-report/main.go
	GOOS=linux GOARCH=386 go build -o bin/main-linux-386 cmd/weather-report/main.go
	GOOS=windows GOARCH=386 go build -o bin/main-windows-386 cmd/weather-report/main.go