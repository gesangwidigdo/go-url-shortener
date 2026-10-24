# build
FROM golang:1.27 AS build

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 go build -trimpath -ldflags="-s -w" -o /out/api ./cmd/api \
  && CGO_ENABLED=0 go build -trimpath -ldflags="-s -w" -o /out/migrate ./cmd/migrate


# run
FROM alpine:3.20
RUN adduser -D -u 10001 app
WORKDIR /app

COPY --from=build /out/api /out/migrate ./

USER app
EXPOSE 8080
CMD [ "./api" ]