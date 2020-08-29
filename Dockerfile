FROM golang:1.15

LABEL maintainer="Michael Senn <michael@morrolan.ch>"

# Create non-privileged user
RUN adduser --home /go/src/probius --system probius

WORKDIR /go/src/probius

# Download dependencies before compiling code, to allow caching if only the
# code changed.
COPY go.mod go.sum /go/src/probius/
RUN go get -d -v ./...

COPY . .
RUN go install -v ./...

CMD ["probius"]
