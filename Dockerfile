####
# Base Go build
####
FROM golang:1.17 as build
ENV GO111MODULE on

# Warm up the module cache.
# Only copy in go.mod and go.sum to increase Docker cache hit rate.
COPY go.mod go.sum /src/
WORKDIR /src
RUN go mod download

COPY . /src

RUN go build -v -o app

####
# Final build
####
FROM gcr.io/distroless/base-debian10:debug

# Copy recurses by default with slash at the end
COPY --from=build /src/app/ /app/
COPY css/ /app/css/
COPY templates/ /app/templates/
COPY md-file-viewer-docker.conf /app/md-file-viewer.conf

WORKDIR /app

EXPOSE 8080

ENTRYPOINT ["./app"]
