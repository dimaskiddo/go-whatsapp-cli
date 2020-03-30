# Builder Image
# ---------------------------------------------------
FROM dimaskiddo/alpine:go-1.12 AS go-builder

WORKDIR /usr/src/app

COPY . ./

RUN go mod download \
    && CGO_ENABLED=0 GOOS=linux go build -a -o dist/go-whatsapp *.go


# Final Image
# ---------------------------------------------------
FROM dimaskiddo/alpine:base
MAINTAINER Dimas Restu Hidayanto <dimas.restu@student.upi.edu>

ARG SERVICE_NAME="go-whatsapp-cli"
ENV PATH $PATH:/opt/${SERVICE_NAME}

WORKDIR /opt/${SERVICE_NAME}

COPY --from=go-builder /usr/src/app/dist/go-whatsapp ./go-whatsapp
COPY share/ ./share

RUN chmod 775 share

CMD ["go-whatsapp", "daemon"]
