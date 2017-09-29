build:
	GOOS=linux GOARCH=386 go build -o mapmyrun-to-garmin_linux-386 ./cmd/translate/main.go && \
  GOOS=windows GOARCH=386 go build -o mapmyrun-to-garmin_windows-386 ./cmd/translate/main.go \
	go build -o mapmyrun-to-garmin-web ./cmd/web/main.go
