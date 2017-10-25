from golang:1.9

workdir /go/src/tim

copy . .

run go get ./...
run go build -o /app/tim

workdir /app

CMD ["./tim"]