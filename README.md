# "Часть сервиса аутентификации"

*Используемые технологии:*
1. Go
2. JWT
3. MongoDB

*Задание:*
  Написать часть сервиса аутентификации.
  Два REST маршрута:
- Первый маршрут выдает пару Access, Refresh токенов для пользователя с идентификатором (GUID) указанным в параметре запроса
- Второй маршрут выполняет Refresh операцию на пару Access, Refresh токенов

## Решение
1. Создайте файл **.env** со следующими переменными:
```golang
MONGO_INITDB_ROOT_USERNAME=root  //оставляем, если используем контейнер
MONGO_INITDB_ROOT_PASSWORD=example //оставляем, если используем контейнер
MONGO_PORT=27017 //оставляем, если используем контейнер
DB_NAME=jwtDB //оставляем, если используем контейнер
COLLECTION_NAME=users //оставляем, если используем контейнер
JWT_KEY="sUper-SecrEt-JWT-keY-!" //секретный ключ для JWT токена
URL_MONGO="" // если тестим на базе данных mongo не из контейнера, здесь указываем строку подключения
```
  *Если тестирование происходит не в контейнере, в файле **database.go** комментируем строки 20-22, убираем комментарий со строки 24.*

2. **Поднимаем контейнер**
```bash
docker-compose build
docker compose up
```
  В проекте предусмотерн docker-compose файл, который поднимает контейнер mongo, а так же скрипт init-mongo.js, который создает базу данных и записывает в коллекцию некоторые тестовые данные.

3. **Запускаем приложение**
```bash
go run .\main.go
```
  Если подключение к базе данных успешно - в командной строке выведется сообщение.

4. **Тестирование через Postman**

  В проекте предусмотрен файл **REST API basics- CRUD, test & variable.postman_collection.json**, содержащий коллекцию из двух запросов, на создание и обновление токенов.
