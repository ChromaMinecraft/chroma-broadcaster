FROM golang:alpine AS build_base

ENV GO111MODULE=on \
  CGO_ENABLED=0 \
  GOOS=linux \
  GOARCH=amd64

RUN apk add --no-cache git 

WORKDIR /tmp/chroma-broadcaster

COPY go.mod .
COPY go.sum .

RUN go mod download

COPY . .

RUN go build -o app .

# FROM alpine:3.9 

# COPY --from=build_base /tmp/chroma-broadcaster/app /app/chroma-broadcaster

EXPOSE 8199

CMD ["./app"]
