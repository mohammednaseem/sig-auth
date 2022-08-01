FROM golang:1.18 as build-env
WORKDIR /app/
COPY . ./
RUN go mod download
RUN go get -d -v ./... 
RUN go vet -v ./...
RUN go test -v ./...
RUN CGO_ENABLED=0 go build -o sig-auth app/main.go
FROM gcr.io/distroless/static
LABEL "microservice.name"="Signature Authentication"
COPY --from=build-env /app/sig-auth /
COPY --from=build-env /app/config.json /
CMD ["/sig-auth"]