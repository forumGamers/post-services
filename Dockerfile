FROM golang:1.20.1-alpine

# Set working directory ke dalam folder "bin"
WORKDIR /app/bin

# Copy file main ke dalam container
COPY ./ ./

# Build aplikasi Go
RUN go mod tidy

RUN go build main.go

# Jalankan aplikasi saat container dijalankan
CMD ["./main"]