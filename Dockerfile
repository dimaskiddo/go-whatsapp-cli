# Builder Image
# ---------------------------------------------------
FROM dimaskiddo/alpine:go-1.12 AS go-builder

WORKDIR /usr/src/app

COPY . ./

RUN go mod download \
    && CGO_ENABLED=0 GOOS=linux go build -a -o go-whatsapp-cli cmd/main/main.go


# Final Image
# ---------------------------------------------------
FROM dimaskiddo/alpine:base
MAINTAINER Dimas Restu Hidayanto <dimas.restu@student.upi.edu>

WORKDIR /usr/app

COPY --from=go-builder /usr/src/app/go-whatsapp-cli ./go-whatsapp-cli

CMD ["./go-whatsapp-cli"]
