FROM golang:alpine
WORKDIR /app
COPY go.mod ./
COPY go.sum ./
RUN go mod download
# Copy the code into the container
COPY . .
ENV SERVICE queue
RUN go build -o app services/${SERVICE}/main.go

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=0 /app/app ./
CMD [ "./app" ]
