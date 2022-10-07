FROM golang:1.19-alpine as build
WORKDIR /go/src/app
COPY go.mod ./
COPY go.sum ./
RUN go mod download
# Copy the code into the container
COPY . .
ENV SERVICE queue
RUN CGO_ENABLED=0 go build -o /go/bin/app services/${SERVICE}/main.go

FROM gcr.io/distroless/static
COPY --from=build /go/bin/app /
CMD [ "/app" ]
