# ozon-fintech

<!-- ToC start -->
# Содержание

1. [Задание](#Описание-задачи)
2. [Реализация](#Реализация)
3. [API](#API)
4. [Сборка, запуск и тестирование](#Сборка-и-запуск)
<!-- ToC end -->

# Задание

Необходимо реализовать сервис, который должен предоставлять API по созданию сокращенных ссылок следующего формата:
- Ссылка должна быть уникальной и на один оригинальный URL должна ссылаться только одна сокращенная ссылка.
- Ссылка должна быть длинной 10 символов
- Ссылка должна состоять из символов латинского алфавита в нижнем и верхнем регистре, цифр и символа _ (подчеркивание)


Сервис должен быть написан на Go и принимать следующие запросы по http:
1. Метод Post, который будет сохранять оригинальный URL в базе и возвращать сокращённый
2. Метод Get, который будет принимать сокращённый URL и возвращать оригинальный URL

Решение должно быть предоставлено в «конечном виде», а именно:
- Сервис должен быть распространён в виде Docker-образа
- В качестве хранилища ожидается использовать 2 реализации: inmemory и postgresql.
Какое хранилище использовать указывается параметром при запуске сервиса.
- Покрыть реализованный функционал Unit-тестами

## Параметры запуска

http://localhost:8080/api/tokens

**Пример реализации сокращения ссылки**

Для проведения тестовых запросов можно использовать Postman, следует установить Content-Type и написать запрос к localhost:8080/api/tokens (localhost:8080/api/tokens/:token), либо сторонние утилиты.
```
Request:
Method: POST

http://localhost:8080/api/tokens
header: 'Content-Type: application/json'

{
    "base_url": "https://www.ozon.ru/product/smartfon-lenovo-k13-2-32gb-siniy-298657201/?sh=1gscsddsssal%2Fasaaaaa"
}


Response:
{
    "token": "WmXge0ssss"
}
```

**Получить оригинальный URL**
```
Request:
Method: GET

http://localhost:8080/api/tokens/WmXge0ssss
header: 'Content-Type: application/json'

Response:

{
    "base_url": "https://www.ozon.ru/product/smartfon-lenovo-k13-2-32gb-siniy-298657201/?sh=1gscsddsssal%2Fasaaaaa"
}

```


# Реализация

- Архитектурный стиль REST API.
- Clean Architecture & dependency injection
- Взаимодействие с РСУБД Postgres реализовано с использованием библиотеки sqlx, SQL запросы реализованы без использования сторонних фреймворков
с целью оптимизации
- Конфигурация, реализована при помощи библиотеки viper
- Запуск из Docker, либо непосредственно с самого приложения, для сборки и запуска следует использовать прописанные команды в Makefile (приведены ниже),
либо проводить непосредственно из консоли
- Юнит-тестирование уровней бизнес-логики с помощью mock object [golang/mock]

# API

> 1) Тело запроса/ответа - в формате JSON.
> 2) Обработка ошибок реализована путем возвращения соответствующего HTTP кода, в теле содержится описание возникшей ошибки

# Параметры запуска

## С помощью Makefile

### Запуск в режиме inmemory
```
make run_im

...или без makefile: go run ./cmd./main.go
```
Можно тестировать запросы, которые сохраняются в map самого приложения (inmemory)

### Запуск в режиме postgresql
```
make run_db

...или без makefile: go run ./cmd./main.go -db
*Флаг -db указывает на подключение к СУБД и используется в проверке условия подключения в main
```

1. Конфигурация для запуска на localhost (Или используйте соответствующий коммит с уже собранным конфигурационным файлом),
изменить (проверить) параметры в configs/config.yml:
```
PORT: "8080"

db:
  host: "localhost"
  port: "5436"
  user: "postgres"
  password: "mrv8336"
  dbname: "postgres"
  sslmode: "disable"
  
  
```  
2. Запустить докер контейнер:
```
  docker run --name=ozon-fintech -e POSTGRES_PASSWORD='mrv8336' -p 5436:5432 -d --rm postgres
```
4. Применить миграции к базе данных (миграции должны быть применены из файла с проектом):
```
 migrate -path ./schema -database 'postgres://postgres:mrv8336@localhost:5436/postgres?sslmode=disable' up
```
5.  Подключаемся к контейнеру:
```
docker exec -it *CONTAINER ID FROM "docker ps"* /bin/bash
```
6.  Можно тестировать запросы на базу данных postgres

## С помощью docker-compose
Для запуска образа Docker через docker-compose, следует привести файл конфигурации config/configs.yml к следующему виду:
```
PORT: "8080"

db:
  host: "db"
  port: "5432"
  user: "postgres"
  password: "mrv8336"
  dbname: "postgres"
  sslmode: "disable"
```
### Запуск в режиме inmemory
```
docker-compose up --build app_im
```
Последовательный ввод в терминале, сборка по docker-compose.yml
### Запуск в режиме postgresql
```
docker-compose up --build app_db
```
