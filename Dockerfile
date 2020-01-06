FROM golang

LABEL key="MINT"

RUN go env -w GO111MODULE=on
RUN go env -w GOPROXY=https://goproxy.cn,direct

RUN mkdir /dim-fs
WORKDIR /dim-fs
COPY go.mod .
COPY go.sum .

RUN mkdir /upload

RUN go mod download
COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags "-X main.BuildEnv=prod" main.go

EXPOSE 9089
EXPOSE 9088

ENTRYPOINT ["./main"]
