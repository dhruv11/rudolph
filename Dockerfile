FROM golang

RUN mkdir -p /app

WORKDIR /app

ADD . /app

RUN go get github.com/adlio/trello

RUN go get github.com/nlopes/slack

RUN go build ./main.go

CMD ["./main"]