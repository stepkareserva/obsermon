# What is Obsermon?

Obsermon (*OBSERvation & MONitoring*) is my educational project.

This is small [`server`](https://github.com/stepkareserva/obsermon/tree/main/cmd/server) to monitor client's state and collect statistics.

Package also contains  [`agent`](https://github.com/stepkareserva/obsermon/tree/main/cmd/agent) as sample of client compatible with server

## Service usage

### CLI

Run server: `go run cmd/server/main.go`

Params: 

- `-a` (string) - server endpoint tcp address, like `:8080`, `127.0.0.1:80`, `localhost:22` (default `localhost:8080`)

### Env

Env params overrides command line params, if exist:

- `ADDRESS` - same as CLI `-a` 

## Agent usage

### CLI

Run server: `go run cmd/agent/main.go`

Params: 

- `-a` (string) - server endpoint tcp address, like `:8080`, `127.0.0.1:80`, `localhost:22` (default `localhost:8080`)

- `-p` (int) - poll (local metrics update) interval, in seconds, positive integer (default 2)

- `-r` (int) - report (send metrics to server) interval, in seconds, positive integer (default 10)

### Env

Env params overrides command line params, if exist:

- `ADDRESS` - same as CLI `-a` 
- `POLL_INTERVAL` - same as CLI `-p` 
- `REPORT_INTERVAL` - same as CLI `-r`

## Service API

**WIP** learn REST description rules

- `POST /update/counter/name/value` - update counter, value is int
- `POST /update/gauge/name/value` - update gauge, value is float
- `GET /counter/name/value` - get counter value, 404 if not exists
- `GET /gauge/name/value` - get gauge value, 404 if not exists
- `GET /` - html page with all counters and gauges

## Monitoring page example

![monitoring](https://raw.githubusercontent.com/stepkareserva/obsermon/refs/heads/main/assets/metrics_sample.png)

## Обновление шаблона

Чтобы иметь возможность получать обновления автотестов и других частей шаблона, выполните команду:

```
git remote add -m main template https://github.com/Yandex-Practicum/go-musthave-metrics-tpl.git
```

Для обновления кода автотестов выполните команду:

```
git fetch template && git checkout template/main .github
```

Затем добавьте полученные изменения в свой репозиторий.

## Запуск автотестов

Для успешного запуска автотестов называйте ветки `iter<number>`, где `<number>` — порядковый номер инкремента. Например, в ветке с названием `iter4` запустятся автотесты для инкрементов с первого по четвёртый.

При мёрже ветки с инкрементом в основную ветку `main` будут запускаться все автотесты.

Подробнее про локальный и автоматический запуск читайте в [README автотестов](https://github.com/Yandex-Practicum/go-autotests).

##
![footer](https://raw.githubusercontent.com/stepkareserva/obsermon/refs/heads/main/assets/footer.svg)