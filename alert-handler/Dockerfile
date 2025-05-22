FROM golang:1.19
WORKDIR /app
COPY . .
RUN go mod download
RUN go build -o alert-handler .
CMD ["./alert-handler"]