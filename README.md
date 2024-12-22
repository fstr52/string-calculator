<h1 align=center>
    <b>
        String-Calculator
    <b>
</h1>

## О проекте

Данный проект представляет собой простой веб-сервис для вычисления арифметических выражений.

## Запуск

1. Установите [Go](https://go.dev/doc/install)
2. Установите [Git](https://git-scm.com/downloads) (при использовании далее способа с клонированием через git clone)
3. Склонируйте проект через команду:
    ```console
    git clone https://github.com/fstr52/string-calculator
    ```

    Или просто скачайте ZIP-архив проекта (зеленая кнопка Code над файлами проекта, затем Download ZIP)
4. Перейдите в директорию проекта
5. Запустите приложение через команду:
    ```console
    go run ./cmd/app
    ```
    Или скомпилируйте файл командой и запустите полученный исполняемый файл:
    ```console
    go build ./cmd/app/main.go
    ```
6. Сервис доступен по адресу: `localhost:8080/api/v1/calculate`


## Конфигурация запуска

Для смены порта запуска измените параметр в файле config.json

## Примеры использования 

Curl запрос:
```bash
curl --location "localhost:8080/api/v1/calculate" --header "Content-Type: application/json" --data "{\"expression\": \"12*(1+2*(1+2)+3)+1\"}"
```

Тело запроса (для простоты визуализации и понимания):
```json
{
    "expression": "12*(1+2*(1+2)+3)+1"
}
```

Ответ:
```json
{"result":"121"}
```
HTTP статус:
```
200 OK
```

## Все возможные результаты запросов

### Результат калькуляции и код `200 OK`:

*Ниже приведены curl-запросы, но рекомендую использовать [Postman](https://www.postman.com/downloads/) для удобства*

Запрос: 
```bash
curl --location "localhost:8080/api/v1/calculate" --header "Content-Type: application/json" --data "{\"expression\": \"2+3\"}"
```
Ответ:
```json
{"result":"5"}
```
HTTP статус:
```
200 OK
```

### Ошибка и HTTP status code
1. **Неверное выражение** <br>
    Запрос: 
    ```bash
    curl --location "localhost:8080/api/v1/calculate" --header "Content-Type: application/json" --data "{\"expression\": \"(2+3\"}"
    ```
    Ответ:
    ```json
    {"error":"Expression is not valid"}
    ```
    HTTP статус:
    ```
    422 Unprocessable Entity
    ```
2. **Неверный формат ввода**<br>
    Запрос: 
    ```bash
    curl --location "localhost:8080/api/v1/calculate" --header "Content-Type: text/plain" --data "{\"expression\": \"2+3\"}"
    ```
    Ответ:
    ```json
    {"error":"Expression is not valid. Expected JSON format input"}
    ```
    HTTP статус:
    ```
    400 Bad Request
    ```
3. **Неверный метод запроса**<br>
    Запрос: 
    ```bash
    curl --location --request GET  "localhost:8080/api/v1/calculate"  --header "Content-Type: application/json"  --data "{\"expression\": \"2+3\"}"
    ```
    Ответ:
    ```json
    {"error":"Method not allowed"}
    ```
    HTTP статус:
    ```
    405 Method Not Allowed
    ```
4. **Непредвиденная ошибка**<br>
    Ответ:
    ```json
    {"error":"Internal server error"}
    ```
    HTTP статус:
    ```
    500 Internal server error
    ```

## Примечание

- Поддерживаются стандартные арифметические операции
- Поддерживаются только POST запросы
- Поддерживается использование унарного минуса, но не плюса