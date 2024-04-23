# Use the latest Golang base image
FROM golang:latest

# Set the current working directory inside the container
WORKDIR /app

# Copy go.mod and go.sum files to the workspace
COPY go.mod go.sum ./

# Download all dependencies
RUN go mod download

# Update package lists
RUN apt-get update

# Install necessary packages for the DNS update tool
RUN apt-get install -y curl perl libwww-perl vim

# Copy the source from the current directory to the workspace
COPY . .

# Download and setup DNS update tool
RUN wget https://dinahosting.com/utilidades/estandar/aplicaciones/dinaIP-consola.tar.gz
RUN tar xzpf dinaIP-consola.tar.gz
WORKDIR /app/dinaIP-consola
RUN sh install.sh
WORKDIR /app

ARG DINAHOSTING_DOMAIN
ENV DINAHOSTING_DOMAIN=$DINAHOSTING_DOMAIN
ARG DINAHOSTING_USER
ENV DINAHOSTING_USER=$DINAHOSTING_USER
ARG DINAHOSTING_PASSWORD
ENV DINAHOSTING_PASSWORD=$DINAHOSTING_PASSWORD

# Build the Go app
RUN go build -o main .

# Expose port 80 to the outside world
EXPOSE 80

# Expose port 443 to the outside world
EXPOSE 443

# Command to run the executabl
CMD dinaip -u $DINAHOSTING_USER -p $DINAHOSTING_PASSWORD -a $DINAHOSTING_DOMAIN && ./main
