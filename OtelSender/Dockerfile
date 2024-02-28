FROM mcr.microsoft.com/oss/go/microsoft/golang:1.21-bookworm

# Set the working directory within the container  
WORKDIR /app  
  
# Copy the source code to the container  
COPY . .  

# Build the Go application  
RUN go build -o otelsender .  
  
# Set the entry point for the container  
CMD ["./otelsender"]  