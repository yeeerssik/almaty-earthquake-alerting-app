# syntax=docker/dockerfile:1

FROM golang:1.22.1

# Set destination for COPY
WORKDIR /

# Download Go modules
COPY go.mod go.sum ./
RUN go mod download

# Copy the source code
COPY *.go ./

# Copy additional files
COPY config.env ./

# Build
RUN CGO_ENABLED=0 GOOS=linux go build -o /almaty_earthquake_alerting_app_test

EXPOSE 8000

# Run
CMD ["/almaty_earthquake_alerting_app_test"]