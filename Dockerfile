FROM golang:1.20.2-alpine3.17 AS builder

# Add 'make' utility to the image
RUN apk add make
WORKDIR /scratch

# Pre download the go.mod dependencies
COPY src/go.mod .
COPY src/go.sum .
RUN go mod download -x
RUN rm go.mod go.sum

COPY tools/rest-gen/go.mod .
COPY tools/rest-gen/go.sum .
RUN go mod download -x
RUN rm go.mod go.sum

# Copy sources from local into container
COPY . .

# Build the stuff
RUN make local

# Prepare deployment image
FROM alpine:3.17 AS deployimg
EXPOSE 80-82
WORKDIR /root/
COPY --from=builder /scratch/app ./app/
CMD ["/root/app/apisrv/apisrv"]
