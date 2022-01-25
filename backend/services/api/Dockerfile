FROM golang:1.17-alpine as build
WORKDIR /go/src/app
COPY go.mod ./
COPY go.sum ./
RUN go mod download
# Copy the code into the container
COPY . .
ENV SERVICE api
RUN CGO_ENABLED=0 go build -o /go/bin/app services/${SERVICE}/main.go

FROM gcr.io/distroless/static
COPY --from=build /go/bin/app /
# TODO: expose ports?
CMD [ "/app" ]
