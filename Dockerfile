FROM golang:1.13-alpine AS go-build
WORKDIR /app
COPY . ./
RUN GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o server

FROM scratch
WORKDIR /opt/game
COPY --from=go-build /app/server /opt/game/server
EXPOSE 80
CMD ["/opt/game/server"]