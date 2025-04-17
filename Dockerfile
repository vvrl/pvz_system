# Используем официальный образ Go
FROM golang:1.23.8

# Устанавливаем рабочую директорию внутри контейнера
WORKDIR /app

# Копируем модули и загружаем зависимости
COPY go.mod ./
# COPY go.sum ./
RUN go mod download

# Копируем весь проект
COPY . .

# Собираем исполняемый файл
# RUN go build -o main ./cmd/pvzSystem

RUN go run ./cmd/pvzSystem/main.go

# Открываем порт
EXPOSE 8080

# Команда по умолчанию при запуске контейнера
CMD ["./main"]
