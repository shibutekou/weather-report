FROM golang:1.19-alpine

WORKDIR ./app

COPY ./ ./

RUN go mod download && go mod tidy && go build cmd/weather-report/main.go

CMD [ "./main" ]