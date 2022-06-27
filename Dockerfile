FROM golang:1.18-alpine as builder

WORKDIR /terraform-controller-worker

COPY go.mod ./
COPY go.sum ./

RUN go mod download

COPY . .

RUN go build -o ./bin/terraform-controller-worker ./cmd/terraform-controller-worker/main.go

FROM alpine:latest

RUN apk --no-cache add curl

RUN apk add --update docker openrc

RUN rc-update add docker boot

RUN cd /usr/local/bin && \
    curl https://releases.hashicorp.com/terraform/1.2.3/terraform_1.2.3_linux_amd64.zip -o terraform.zip && \
    unzip terraform.zip && \
    rm terraform.zip

COPY --from=builder /terraform-controller-worker/bin/terraform-controller-worker /usr/local/bin

WORKDIR /opt/manifests

CMD ["/usr/local/bin/terraform-controller-worker"]