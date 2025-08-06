# Crypto Observer Service

Микросервис для мониторинга и хранения цен криптовалют с возможностью исторического анализа.

## Функционал

- **Добавление криптовалюты в мониторинг**
`/currency/add` - начинает регулярный сбор цен с указанным интервалом
- **Удаление криптовалюты из мониторинга**
`/currency/remove` - прекращает сбор цен для указанной криптовалюты
- **Получение исторической цены**
`/currency/price` - возвращает цену на запрошенный момент времени

## Технологии

- **Backend**: Go 1.23
- **База данных**: PostgreSQL 17
- **API документация**: SwaggerUI (доступно по `/api/doc/`) или папка `docs` в корне проекта
- **Контейнеризация**: Docker + Docker Compose

## Запуск проекта

### Требования
- Docker 20.10+
- Docker Compose 2.0+

### Инструкция

1. Клонировать репозиторий:
```bash
git clone https://github.com/yourusername/cryptoObserver.git
cd cryptoObserver
```

2. Создать файл `.env`:
```bash
echo "DB_USER=your_db_user" > .env
```

Для теста предлагаю сразу готовый рабочий конфиг (с активным trial ключем от API CoinGecko):
```bash
DB_HOST=localhost
DB_USER=botkai
DB_PASSWORD=passwordfjdk
DB_NAME=botdb
DB_PORT=5433
SERVER_PORT=8080
CRYPTO_API_KEY=
WORKER_POOL_SIZE=10
WORKER_POOL_UPDATE_TIME=60
CRYPTO_API_KEY=CG-kq8Ee8QmdRMM4MA32myqrqxN
```

3. Запустить сервисы:
```bash
docker-compose up --build
```

4. Сервис будет доступен на `http://localhost:8080`

## API Endpoints

| Метод | Путь                | Описание                          |
|-------|---------------------|-----------------------------------|
| POST  | /currency/add       | Добавить криптовалюту в мониторинг|
| POST  | /currency/remove    | Удалить криптовалюту из мониторинга|
| GET   | /currency/price     | Получить историческую цену        |


## Дополнительно

- Пул воркеров для параллельного сбора цен
- Валидация входящих запросов
- Логирование операций
- Health-check эндпоинты

Для доступа к полной документации API после запуска сервиса посетите:
`http://localhost:8080/swagger/index.html`

Полный список id валют доступен в файле `aviable-ids.json`.

Краткий список доступных валют:
1. bitcoin
2. ethereum
3. dogecoin
4. 1ex
