# Golang песочница

Здесь я разбираю паттерны и другие темы на Go.
У каждого примера есть реализация и демонстрация запуска.

## Примеры

| Паттерн | Варианты | Что делает |
| --- | --- | --- |
| [Worker Pool](concurrency/workerpool) | basic, results, cancellable | Считает сумму квадратов несколькими воркерами |
| [Pipeline](concurrency/pipeline) | basic, cancellable | Передает числа по этапам: квадраты → фильтр четных |
| [Semaphore](concurrency/semaphore) | basic, cancellable, weighted | Ограничивает число операций или их общий вес |

## Запуск

Нужен Go 1.26.4 или новее. Команды выполняются из корня репозитория:

```sh
go run ./concurrency/workerpool
go run ./concurrency/pipeline
go run ./concurrency/semaphore
```

Каждая команда запускает все варианты своего примера.
В Worker Pool и Pipeline входные числа случайные, от 10 до 99.
Этап вычисления квадратов в Pipeline работает в нескольких воркерах,
поэтому порядок результатов может меняться.
Примеры с контекстом используют короткий таймаут для отмены.
В семафорах отменяется ожидание разрешений, а запущенные задачи завершают работу.

## Структура

- `concurrency/` — примеры работы с горутинами и каналами.
- `main.go` — запуск вариантов паттерна.
- `example.go` — сценарий демонстрации и вывод результата.
- `workerpool.go`, `pipeline.go`, `semaphore.go` — реализации.
  WorkerPool и Pipeline выводят результаты внутри реализации.
- `pkg/number` — общий генератор чисел для демонстраций.

## Проверки

```sh
go vet ./...
go test -race ./...
go run -race ./concurrency/workerpool
go run -race ./concurrency/pipeline
go run -race ./concurrency/semaphore
```
