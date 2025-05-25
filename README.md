# What is Obsermon?

Obsermon (*OBSERvation & MONitoring*) is my educational project.

This is small [`server`](https://github.com/stepkareserva/obsermon/tree/main/cmd/server) to monitor client's state and collect statistics.

Package also contains  [`agent`](https://github.com/stepkareserva/obsermon/tree/main/cmd/agent) as sample of client compatible with server

## Service usage

Run server: `go run cmd/server/main.go`

Params:

| `CLI`| `ENV` | `type` | `default` | **Description** |
|:-----|:------|:-------|:----------|:----------------|
|`-a`  | `ADDRESS` | `string` | `localhost:8080` |  server endpoint tcp address, like `:8080`, `127.0.0.1:80`, `localhost:22`
|`-i`  | `STORE_INTERVAL` | `int` | `300` | server state file storing interval, s, 0 for sync storing 
|`-f`  | `FILE_STORAGE_PATH` | `string` | `obsermon/storage.json` in appdata (depends on os) | path to server state storage file
|`-r`  | `RESTORE` | `bool` | `false` | restore server state from storage file
|`-d`  | `DATABASE_DSN` | `string` | `""` | database connection string
|`-k` | `KEY` | `string` | `""` | key to sing requests via SHA256
|`-m`  | `MODE` | `string` | `prod` | app mode, `quiet` (no logs), `dev` (human-readable logs), `prod` (machine-readable logs)|


## Agent usage

Run server: `go run cmd/agent/main.go`

Params: 

| `CLI`| `ENV` | `type` | `default` | **Description** |
|:-----|:------|:-------|:----------|:----------------|
|`-a`  | `ADDRESS` | `string` | `localhost:8080` | server endpoint tcp address, like `:8080`, `127.0.0.1:80`, `localhost:22`
|`-p`  | `POLL_INTERVAL` | `int` | `2` | poll (local metrics update) interval, in seconds, positive integer 
|`-r` | `REPORT_INTERVAL` | `int` | `10` | report (send metrics to server) interval, in seconds, positive integer
|`-k` | `KEY` | `string` | `""` | key to sing requests via SHA256
|`-l` | `RATE_LIMIT` | `int` | `1` | max count of requests on the same time

## Service API

**WIP** learn REST description rules

- `POST /update/counter/name/value` - update counter, value is int
- `POST /update/gauge/name/value` - update gauge, value is float
- `POST /update` - update counter or gauge
- `POST /updates` - update batch of metrics (counters and gauges)
- `GET /value/counter/name` - get counter value, 404 if not exists
- `GET /value/gauge/name` - get gauge value, 404 if not exists
- `POST /value` - GET(lol) counter or gauge
- `GET /` - html page with all counters and gauges
- `GET /ping` - check database status

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