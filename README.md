# Weather Lab

Небольшой JSON API для домашнего задания по погоде.

Используются:

- `net/http` и `chi`
- `Open-Meteo` для геокодинга и текущей погоды
- `CountriesNow` для списка городов страны
- `REST Countries` для ISO-кода страны

## Запуск

```bash
make run
```

Приложение стартует на `http://localhost:8080`.

## Make targets

```bash
make fmt
make build
make tidy
make curl-health
make curl-city CITY=Almaty
make curl-country COUNTRY=Kazakhstan
make curl-top COUNTRY=Kazakhstan
```

## API

```bash
GET /weather/{city}
GET /weather/country/{country}
GET /weather/country/{country}/top
```

Примеры:

```bash
curl http://localhost:8080/weather/Almaty
curl http://localhost:8080/weather/country/Kazakhstan
curl http://localhost:8080/weather/country/Kazakhstan/top
```

## Переменные окружения

```bash
HTTP_PORT=8080
REQUEST_TIMEOUT_SEC=10
```
