FROM golang:1.9
WORKDIR /go/src/github.com/blaskovicz/mapmyrun-to-garmin
COPY . .
#RUN go-wrapper download ./...
RUN go-wrapper install ./...
EXPOSE 80
ENV ENVIRONMENT=production PORT=80
#ENV COOKIE_KEY=replace_me CSRF_KEY=replace_me
CMD ["web"]
