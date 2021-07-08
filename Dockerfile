#######################################################################################
# Builder stage is separate: Golang compiles down to a single binary file, so we don't
# need to keep all the dependencies and source code in the final image.
FROM golang:1.15-alpine3.13 AS builder

# Set the working directory outside $GOPATH to enable the support for modules.
WORKDIR /src

# Fetch dependencies first; they are less susceptible to change on every build
# and will therefore be cached for speeding up the next build
COPY ./go.mod ./go.sum ./
RUN go mod download

# Copy the rest of the source code
COPY . ./

ENV CGO_ENABLED=0 \
        GOOS=linux \
        GOARCH=amd64

# And compile the project
RUN go build -o ./out ./

#######################################################################################
# Includes *just* the executable binary.
FROM scratch AS final

COPY --from=builder /src/out /websocketserver

ENTRYPOINT ["/websocketserver"]
