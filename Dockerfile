FROM golang:1.23.4

WORKDIR /myapi

COPY  go.mod ./
COPY go.sum  ./
RUN go mod download

COPY . .


CMD [ "go","run","main.go" ]




