FROM golang:alpine as builder

WORKDIR /go/src
copy . /go/src
run cd /go/src && env GOOS=linux GOARCH=amd64 go build


from alpine
WORKDIR /app
run cd /app
copy --from=builder /go/src/food-tracker .

ENTRYPOINT [ "./food-tracker" ]