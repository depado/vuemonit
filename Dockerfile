FROM node:19.6.0-alpine AS front_builder

ADD ./front /front
WORKDIR /front
RUN npm install && ./node_modules/.bin/quasar build

# Backend Build Step
FROM golang:1.20.0-alpine AS builder

# Prerequisites
RUN apk update && apk add --no-cache upx

# Dependencies
WORKDIR $GOPATH/src/github.com/Depado/vuemonit
COPY . .
RUN go mod download
RUN go mod verify
RUN go get github.com/rakyll/statik

# Copy frontend build
COPY --from=front_builder /front/dist $GOPATH/src/github.com/Depado/vuemonit/front/dist/

# Build
ARG build
ARG version
RUN mv front/dist/spa/index.html front/dist/spa/main.html
RUN statik -src=./front/dist/spa/ -f
RUN mv front/dist/spa/main.html front/dist/spa/index.html
RUN CGO_ENABLED=0 go build -ldflags="-s -w -X main.Version=${version} -X main.Build=${build}" -o /tmp/vuemonit
RUN upx --best /tmp/vuemonit

# Final Step
FROM gcr.io/distroless/static
COPY --from=builder /tmp/vuemonit /go/bin/vuemonit
VOLUME [ "/data" ]
WORKDIR /data
EXPOSE 8080
ENTRYPOINT ["/go/bin/vuemonit", "--server.port", "8080", "--server.host", "0.0.0.0"]
