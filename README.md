# Transfer service

В корневой директории есть *makefile* для упрощения сборки, запуска и простого тестирования.

Команды:
* ```make build``` — компиляция сервиса в бинарный файл,
* ```make run``` — запуск микросервиса (порт по умолчанию ```8081```),
* ```make test``` — быстрая проверка работоспособности сервиса (производится curl-запрос на ```/transfer```).

Входные данные (текущие координаты клиента) сервис берёт из тела запроса. Формат входных данных - ```json```.

Формат ответа (наименьшее время подачи машины) также ```json```.

Пример входных данных: ```{"lat": 17.986511, "lng": 63.441092}```.

Пример выходных данных: ```{"response": 4669}```.

Пример ответа с ошибкой: ```{"error": "cars service is unavailable"}```.

В папке ```/internal``` располагаются клиенты для обращения к сервисам *cars* и *predict*, сгенерированные при помощи *go-swagger*.

В папке ```/transfer``` располагается логика приложения.
