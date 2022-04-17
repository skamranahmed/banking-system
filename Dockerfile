# build stage
FROM golang:1.18-alpine AS builder

# setup the working directory inside the container
WORKDIR /app

# copy the files from the host computer to the container
# everything from the root directory of the project will be copied
COPY . .

# build go binary
RUN go build -o main main.go

# run stage
FROM alpine:3.15

# setup the working directory inside the container
WORKDIR /app

# copy the binary build from the builder to the current stage container
COPY --from=builder /app/main .

# copy the config directory from the host computer to the current stage container
COPY ./config ./config

# tell the docker container to expose the mentioned port
EXPOSE 8080

# define the commands that are to executed when the container starts
CMD [ "/app/main" ]