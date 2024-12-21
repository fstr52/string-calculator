<h1 align=center>
    <b>
        String-Calculator
    <b>
</h1>

## О проекте

Веб-сервис для вычисления арифметических выражений через HTTP-запрос с методом POST

## Запуск

1. Установите [Go](https://go.dev/doc/install)
2. Установите [Git](https://git-scm.com/downloads)
3. Склонируйте проект через команду:
    ```console
    git clone https://github.com/fstr52/string-calculator
    ```
4. Перейдите в дирректорию проекта
5. Запустите приложение через команду:
    ```console
    go run ./cmd/app
    ```
6. Сервис доступен по адресу: `localhost:8080/api/v1/calculate`


## Конфигурация запуска

Для смены порта запуска измените параметр в файле config.json

## Примеры использования 

```shell
curl --location "localhost:8080/api/v1/calculate" --header "Content-Type: application/json" --data "{\"expression\": \"12*(1+2*(1+2)+3)+1\"}"
```

## Все возможные результаты запросов

### Результат калькуляции и код `200 OK`:

*Ниже приведены curl-запросы, но рекомендую использовать [Postman](https://www.postman.com/downloads/) для удобства*

Запрос: 
```shell
curl --location "localhost:8080/api/v1/calculate" --header "Content-Type: application/json" --data "{\"expression\": \"2+3\"}"
```
Ответ:
```shell
{"result":"5"}
200 OK
```

### Ошибка и response status code
1. **Неверное выражение** <br>
    Запрос: 
    ```shell
    curl --location "localhost:8080/api/v1/calculate" --header "Content-Type: application/json" --data "{\"expression\": \"(2+3\"}"
    ```
    Ответ:
    ```shell
    {"error":"Expression is not valid"}
    422 Unprocessable Entity
    ```
2. **Неверный формат ввода**<br>
    Запрос: 
    ```shell
    curl --location "localhost:8080/api/v1/calculate" --header "Content-Type: text/plain" --data "{\"expression\": \"2+3\"}"
    ```
    Ответ:
    ```shell
    {"error":"Expression is not valid. Expected JSON format input"}
    400 Bad Request
    ```
3. **Неверный метод запроса**<br>
    Запрос: 
    ```shell
    curl --location --request GET  "localhost:8080/api/v1/calculate"  --header "Content-Type: application/json"  --data "{\"expression\": \"2+3\"}"
    ```
    Ответ:
    ```shell
    {"error":"Method not allowed"}
    405 Method Not Allowed
    ```
4. **Непредвиденная ошибка**<br>
    Ответ:
    ```shell
    {"error":"Internal server error"}
    500 Internal server error
    ```

