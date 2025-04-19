# Используем официальный образ Go
FROM golang:1.23.8

# Копируем модули и загружаем зависимости
COPY go.mod ./
COPY go.sum ./
RUN go mod download

# Копируем весь проект
COPY ./ ./

# Открываем порт
EXPOSE 8080

# Команда по умолчанию при запуске контейнера
CMD ["go", "run", "./cmd/pvzSystem"]
