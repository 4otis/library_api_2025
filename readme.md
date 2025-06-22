# Library REST API

REST API для управления библиотекой книг

## Функции

### Книги
- `GET /books` - Список всех книг
- `GET /books/:id` - Получить книгу по ID
- `POST /books` - Добавить новую книгу
- `PUT /books/:id` - Обновить книгу
- `DELETE /books/:id` - Удалить книгу

### Авторы
- `GET /authors` - Список всех авторов
- `GET /authors/:id` - Получить автора по ID
- `POST /authors` - Добавить нового автора
- `PUT /authors/:id` - Обновить автора
- `DELETE /authors/:id` - Удалить автора


## QuickStart

### Требования
- Docker + Docker Compose
- Go 1.24+

```bash
# 1. Клонировать репозиторий
git clone https://github.com/4otis/library_api_2025
cd library_api_2025

# 2. Запустить инфраструктуру
make init

# 3. Установить зависимости
go mod download

# 4. Запустить приложение
make run

# 5. Запустить тесты
make tests

# 6. Собрать документацию для Swagger
make clean
make docs
make run
# Открыть в браузере страницу 
http://localhost:1323/swagger/index.html
```
## Технологии

| Компонент       | Версия    |
|-----------------|----------|
| Go              | 1.24+    |
| Echo            | v4       |
| GORM            | v2       |
| PostgreSQL      | 13+      |
| Swagger         | 2.0      |
