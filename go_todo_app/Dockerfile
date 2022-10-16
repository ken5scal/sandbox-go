FROM golang:1.18-bullseye as builder

ENV GO111MODULE=on CGO_ENABLED=0 GOOS=linux GOARCH=amd64

RUN apt-get update

RUN addgroup --system --gid 50001 app && \
    adduser --system --ingroup app \
    --disabled-password --disabled-login --no-create-home \
    --uid 50001 app

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build -trimpath -ldflags "-w -s" -o app

FROM scratch as deploy
COPY --from=builder /usr/share/zoneinfo /usr/share/zoneinfo
ENV TZ=Asia/Tokyo

# copy sytem setting
COPY --from=builder /etc/group /etc/passwd /etc/
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt

# copy app settings
COPY --from=builder /app/app /bin/app
#RUN #pwd
#RUN #ls /bin/
#COPY --from=builder /app/pkg/api/config/*.toml /app/pkg/api/config/

USER app
EXPOSE 80
ENTRYPOINT ["/bin/app"]

# local
FROM golang:1.18 as dev
WORKDIR /app
RUN go install github.com/cosmtrek/air@latest
ENTRYPOINT ["air"]