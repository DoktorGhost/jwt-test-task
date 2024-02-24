// Создание базы данных "test".
db = db.getSiblingDB('jwtDB');

// Создание коллекции "users".
db.createCollection('users');

// Добавление тестовых данных в коллекцию "users".
db.users.insertMany([
  { guid: '1111'},
  { guid: '2222',
  refresh_token: Binary.createFromBase64('JDJhJDEwJEYxM0EuLjBuZ2Q3NHdCclNoL2ZBVHVOUUh5NlhMa0VveC5FR0VLWkZUc3VZRmxnWXcxaHBP', 0)}
]);