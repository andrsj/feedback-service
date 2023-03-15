FROM golang:1.19.7-alpine

# https://stackoverflow.com/questions/64462922/docker-multi-stage-build-go-image-x509-certificate-signed-by-unknown-authorit
ARG cert_location=/usr/local/share/ca-certificates

RUN apk add --no-cache git ca-certificates openssl \
    && openssl s_client -showcerts -connect github.com:443 </dev/null 2>/dev/null|openssl x509 -outform PEM > ${cert_location}/github.crt \
    && openssl s_client -showcerts -connect golang.org:443 </dev/null 2>/dev/null|openssl x509 -outform PEM > ${cert_location}/golang.org.crt \
    && openssl s_client -showcerts -connect go.uber.org:443 </dev/null 2>/dev/null|openssl x509 -outform PEM > ${cert_location}/go.uber.org.crt \
    && openssl s_client -showcerts -connect proxy.golang.org:443 </dev/null 2>/dev/null|openssl x509 -outform PEM >  ${cert_location}/proxy.golang.crt \
    && openssl s_client -showcerts -connect gorm.io:443 </dev/null 2>/dev/null|openssl x509 -outform PEM > ${cert_location}/gorm.io.crt \
    && openssl s_client -showcerts -connect gopkg.in:443 </dev/null 2>/dev/null|openssl x509 -outform PEM > ${cert_location}/gopkg.in.crt \
    && update-ca-certificates

ENV GOSUMDB=off \
    GOPRIVATE=*.com \
    GOPROXY=direct \
    CGO_ENABLED=0 \
    GO111MODULE="on"

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download && go mod verify

COPY . .

RUN go build -v -o ./server ./cmd/main.go

EXPOSE 8080

CMD ["./server"]
