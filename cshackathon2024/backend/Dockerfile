FROM --platform=linux/amd64 golang:alpine as packager

WORKDIR /app

COPY ./go.mod .
COPY ./go.sum .

RUN go mod download

FROM --platform=linux/amd64 golang:alpine as builder

WORKDIR /app

COPY --from=packager /go/pkg/mod /go/pkg/mod

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o utilService

FROM --platform=linux/amd64 alpine as runner

WORKDIR /app

COPY --from=builder /app/utilService .

EXPOSE 8000

ENV GIN_MODE=release

RUN apk add --no-cache tzdata

ENV TZ="Asia/Bangkok"

ENTRYPOINT [ "/app/utilService" ]