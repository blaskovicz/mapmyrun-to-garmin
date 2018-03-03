FROM golang:1.9
WORKDIR /go/src/github.com/blaskovicz/mapmyrun-to-garmin
COPY . .
RUN go-wrapper install ./...
EXPOSE 3091
ENV ENVIRONMENT=production PORT=3091 COOKIE_KEY=test-key CSRF_KEY=test-key-2
CMD ["web"]
