FROM golang:1.24

# Set the working directory in the container to /app
# This is where the application code will be stored
WORKDIR /app

# Copy the go.mod and go.sum files to the /app directory in the container
COPY go.mod ./
COPY go.sum ./
# This command downloads the dependencies specified in the go.mod file
# and stores them in the module cache
RUN go mod download

# Copy all the files from the root directory to the /app directory in the container
COPY . ./

RUN go build -o main ./cmd/main.go

EXPOSE 1000

CMD ["./main"]
