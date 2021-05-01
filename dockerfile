# Start from base image
FROM golang:alpine

# Set the current working directory inside the container
WORKDIR /backend

# Copy go mod and sum files
COPY go.mod go.sum ./
COPY wait-for.sh ./

# Download all dependencies
RUN go mod download

# Copy source from current directory to working directory
COPY . .

# Build the application
RUN go build

# Expose necessary port
EXPOSE 9000

# ENV WAIT_HOSTS=mysql:3306

# Run the created binary executable after wait for mysql container to be up
CMD ["./wait-for.sh" , "mysql:3306" , "--timeout=300" , "--" , "./komodo-backend"]
