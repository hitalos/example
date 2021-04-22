FROM golang:alpine AS builder

WORKDIR /app
ADD . .
RUN CGO_ENABLED=0 go build -ldflags '-s -w'

FROM scratch
LABEL maintainer="Hítalo Silva <hitalos@gmail.com>"

WORKDIR /app
COPY --from=builder /app/example .

CMD ["/app/example"]