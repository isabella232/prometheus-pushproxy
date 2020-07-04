# Dockerfile References: https://docs.docker.com/engine/reference/builder/

# Start from the latest golang base image
FROM golang:alpine AS builder

# Add Maintainer Info
LABEL maintainer="ming luo"

RUN apk --no-cache add build-base git

# Build Delve
RUN go get github.com/google/gops

WORKDIR /root/
ADD . /root
RUN cd /root/src && go build -o ppp

######## Start a new stage from scratch #######
FROM alpine

# RUN apk update
WORKDIR /root/bin
RUN mkdir /root/config/

# Copy the Pre-built binary file and default configuraitons from the previous stage
COPY --from=builder /root/src/ppp /root/bin
COPY --from=builder /root/config/default_config.yml /root/config/default_config.yml

# Copy debug tools
COPY --from=builder /go/bin/gops /root/bin

# Command to run the executable
ENTRYPOINT ["./ppp"]
