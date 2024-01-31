# Use the official Golang base image
FROM golang:1.21.5

# Set the working directory inside the container
WORKDIR /app

# Copy the local code into the container at the working directory
COPY . .

# Build the Go application
RUN go build -o main .

# Expose port 8080 for the application
EXPOSE 8080

ARG DD_GIT_COMMIT_SHA
ENV DD_TAGS="git.repository_url:github.com/jon94/learn-golang,git.commit.sha:${DD_GIT_COMMIT_SHA}"

# Command to run the executable
CMD ["./main"]