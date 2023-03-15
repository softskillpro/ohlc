# Start from golang base image
FROM golang:1.19.5-bullseye

# Setup Folders
RUN mkdir /app
WORKDIR /app

# Copy the source from the current directory to the working Directory inside the container
COPY ./ ./

ENV GONOPROXY="github.com/*"
# ENV GOPROXY="https://goproxy.io,direct"

# RUN go get -d -v ./../cmd/api/
RUN go mod tidy

# Build the Go app
RUN go build -o . ./...

# Expose port 8080 to the outside world
EXPOSE 8080

# Run the executable
CMD [ "./http" ]