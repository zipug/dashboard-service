FROM golang:1.23-alpine AS builder

ARG GOOS=linux
ARG GOARCH=amd64
WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download
COPY . .

WORKDIR /app/cmd/dashboard/
RUN go get -v ./... \
&& go install -v ./... \
&& go build -v -o server

# EXPOSE 8080

FROM scratch AS production

COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /app/cmd/dashboard/server /app/server

# HEALTHCHECK --interval=5s --timeout=3s --retries=3 \
# CMD ["CMD-SHELL", "curl --fail http://localhost:4400/api/v1/health | grep ok && echo 0 || echo 1"]

ENTRYPOINT [ "/app/server" ]
