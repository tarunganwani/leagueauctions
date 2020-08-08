# Use the official image as a parent image.
FROM golang:1.14

# Set the working directory.
WORKDIR /usr/src/app

#set required environment variables
ENV CERTS_DIR /usr/src/app/server/certs

# Copy the file from your host to your current location.

RUN mkdir server
COPY server server/
COPY go.mod .
COPY go.sum .
RUN cd server

# Run the command inside your image filesystem.
RUN go mod download

# Add metadata to the image to describe which port the container is listening on at runtime.
EXPOSE 8081

# Run the specified command within the container.
# CMD [ "go", "run", "servicemain/main/run.go" ]
