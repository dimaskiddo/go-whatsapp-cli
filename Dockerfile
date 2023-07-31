# Builder Image
# ---------------------------------------------------
FROM dimaskiddo/alpine:go-1.19 AS go-builder

WORKDIR /usr/src/app

COPY . ./

RUN go mod download \
    && CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -a -o main cmd/main/main.go


# Final Image
# ---------------------------------------------------
FROM dimaskiddo/alpine:base
MAINTAINER Dimas Restu Hidayanto <dimas.restu@student.upi.edu>

ARG SERVICE_NAME="go-whatsapp-cli"
ENV PATH="$PATH:/usr/app/${SERVICE_NAME}"

WORKDIR /usr/app/${SERVICE_NAME}

COPY --from=go-builder /usr/src/app/config/ ./config
COPY --from=go-builder /usr/src/app/main ./go-whatsapp-cli

RUN chmod 777 ./config/stores

VOLUME ["/usr/app/${SERVICE_NAME}/config/stores"]
CMD ["go-whatsapp-cli","daemon"]
