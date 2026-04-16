## Запуск

```bash
make run
```

## Make targets

```bash
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
