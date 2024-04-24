FROM golang:alpine
WORKDIR /app
COPY . .
RUN go build -o veetro cmd/veetro/veetro.go
CMD [ "./veetro" ]