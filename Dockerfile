FROM golang:1.9 as gobuild
WORKDIR /go/src/github.com/blaskovicz/mapmyrun-to-garmin
COPY . .
RUN go-wrapper install ./...

FROM node:8.11
WORKDIR /go/src/github.com/blaskovicz/mapmyrun-to-garmin
COPY --from=gobuild /go/src/github.com/blaskovicz/mapmyrun-to-garmin .
COPY --from=gobuild /go/bin .
RUN yarn install && yarn build

EXPOSE 3091
ENV PORT=3091
CMD ["./web"]
