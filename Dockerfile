FROM golang:1.20-alpine

WORKDIR /app

ENV GO111MODULE=on
COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY *.go ./
COPY . ./

RUN go build -o /TestManager

EXPOSE 8081

CMD [ "/TestManager" ]
