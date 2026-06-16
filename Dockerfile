# Используем официальный образ Golang
FROM --platform=linux/amd64 golang:1.25 AS builder

# Устанавливаем рабочую директорию
WORKDIR /app

# Копируем go.mod и go.sum
COPY go.mod go.sum ./

# Загружаем зависимости
RUN go mod download

# Копируем исходный код
COPY . .

# Собираем приложение
RUN CGO_ENABLED=0 GOOS=linux go build -o /app cmd/account-service/main.go

# Используем более легкий образ для запуска
FROM alpine:latest

# Устанавливаем настройки для go
ENV GOARCH="amd64/whatever"
ENV GOOS="linux/whatever"

# Устанавливаем необходимые библиотеки
RUN apk --no-cache add ca-certificates

# Копируем собранное приложение из предыдущего этапа
COPY --from=builder /app .

RUN mkdir -p /migrations
COPY --from=builder /app/build/app/migrations /migrations

# Указываем команду для запуска приложения
CMD ["/main"]