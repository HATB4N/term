FROM golang:1.24.1

RUN apt update && apt install -y python3 python3-pip iputils-ping

WORKDIR /app

COPY . .

RUN [ ! -f go.mod ] && go mod init scan || true

RUN go mod tidy

RUN go build -o scan scan.go

ENTRYPOINT ["python3", "main.py"]