# Topological Pipeline

Исследовательский CSP-прототип потокового вычислительного узла на Go.

## Живой контур

```text
Eternal Source
      |
      v
Round-Robin Dispatcher
   /   |   \
 W1   W2   ... Wn
   \   |   /
      v
   Smoother
      |
      v
     Sink
```

`Dispatcher + Workers + Smoother` образуют один `Parallel Node`. Количество workers меняет внутреннюю степень параллелизма, но не внешний контракт узла.

Все каналы текущего эксперимента небуферизованы. Медленный Sink блокирует Smoother, затем workers, Dispatcher и Source: backpressure проходит через весь тракт без неограниченного накопления в памяти.

## Жизненный цикл

Обычный режим работы не ограничен по времени. Отсутствие сообщений означает паузу, а не EOF.

`SIGINT` отменяет корневой `context.Context`. Контекст останавливает только Source; остальные процессы завершаются через каскад drain и закрытий:

```text
Source closes output
→ Dispatcher drains and closes worker inputs
→ Workers finish
→ Smoother drains and closes output
→ Sink drains and returns from Consume
→ main exits
```

Канал закрывает его отправляющая сторона либо координатор, который доказал завершение всех отправителей.

## Запуск

Обычный режим без диагностических логов:

```text
go run ./bin
```

Остановить процесс можно через `Ctrl-C`.

Диагностический режим:

```text
go run -tags debug ./bin
```

## Проверка

```text
go test ./...
go vet ./...
go test -race ./...
```

Конечный integration test передаёт 128 сообщений через полный тракт и проверяет 128 результатов: 64 нуля и 64 единицы. Отдельный тест подтверждает закрытие выхода отменённого Source.

## Граница прототипа

Пока намеренно не реализованы:

- ёмкости каналов;
- восстановление исходного ordering;
- sequence ID и reorder buffer;
- динамическое изменение числа workers;
- распределённый `Wire`;
- hard abort с потерей сообщений.

Текущая цель — минимальное локальное ядро с доказуемыми backpressure и graceful shutdown. История экспериментов и принятые термины сохранены в `context/`.
