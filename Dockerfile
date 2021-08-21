# Use an official Go runtime as a parent image
FROM  golang:1.16
# Set the Current Working Directory inside the container
WORKDIR /app/weather_backend
# Copy all files from the current directory to the PWD (Present Working Directory) inside the container
COPY . .
# Download all the dependencies from go.mod file
RUN go mod download
# tidy up dependencies in go.sum file and project files
RUN go mod tidy
# Build the Go app
RUN go build -o ./app/out/weather_backend .
# This container exposes port 8080 to the outside world
EXPOSE 8080
# Run the binary program produced by `go build` from out folder
CMD ["./app/out/weather_backend"]