# build stage
FROM golang:1.18-alpine AS builder

# setup the working directory inside the container
WORKDIR /app

# copy the files from the host computer to the container
# everything from the root directory of the project will be copied
COPY . .

# build go binary
RUN go build -o main main.go

# install curl in the builder
RUN apk add curl

# download golang-migrate
RUN curl -L https://github.com/golang-migrate/migrate/releases/download/v4.15.1/migrate.linux-amd64.tar.gz | tar xvz

# run stage
FROM alpine:3.15

# setup the working directory inside the container
WORKDIR /app

# copy the binary build from the builder to the current stage container
COPY --from=builder /app/main .

# copy the golang-migrate zip from the builder to the current stage container
COPY --from=builder /app/migrate ./migrate

# copy the migration files from the host computer to the current stage container
COPY db/migration ./migration

# copy the config directory from the host computer to the current stage container
COPY ./config ./config

# copy the entrypoint.sh file from the host computer to the current stage container
COPY entrypoint.sh .

# tell the docker container to expose the mentioned port
EXPOSE 8080

# define the command that are to executed when the container starts
CMD [ "/app/main" ]

# When the CMD instruction is used together with ENTRYPOINT, then the CMD acts as 
# additional parameters that will be passed to the ENTRYPOINT script
ENTRYPOINT [ "/app/entrypoint.sh" ]