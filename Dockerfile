FROM golang

LABEL key="MINT"

ENV GOPROXY=https://goproxy.io
ENV GO111MODULE=on

RUN mkdir /dim-step-core-api
WORKDIR /dim-step-core-api
COPY go.mod .
COPY go.sum .

RUN mkdir /upload

RUN go mod download
COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags "-X main.BuildEnv=prod" main.go

RUN rm -rf .vscode .git model protocol service utils
RUN rm -f .gitignore go.mod go.sum main.go 

EXPOSE 9089

ENTRYPOINT ["./main"]
