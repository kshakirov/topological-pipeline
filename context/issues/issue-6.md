# Issue #6 — Первая оценка эксперимента


**State:** OPEN
**Author:** @kshakirov
**Created:** 2026-08-26T14:34:01Z
**Updated:** 2026-08-31T14:28:08Z
**URL:** https://github.com/kshakirov/topological-pipeline/issues/6

---

# Parallel Node: ownership, termination and flow semantics

## Контекст

Текущий макет намеренно раскладывает **один параллельный узел** на три наблюдаемые части:

```text
External Input
      |
      v
Parallel Composition
   /   |   \
 W1   W2   ... Wn
   \   |   /
      v
   Smoother
      |
      v
External Output
```

В дальнейшем это единая сущность — аналог `Bolt` в Storm.

Степень параллелизма определяется количеством workers `n`.
Round-Robin пока является первой простой политикой диспетчеризации.
Go channels — временная локальная модель будущего распределённого транспорта.

## Текущее состояние

Гонка от конкурентного `Dump()` устранена: workers пишут в `LocalSplitBuffer.InChan`, а `Interleave()` единолично изменяет slice.

Обнаружена другая проблема: `LocalSplitBuffer` копируется по значению внутрь `LocalSplitNode`.

Channel после копирования остаётся общим, но slice header копируется. Поэтому `Interleave()` может изменять состояние, которое основной `lsp.Buffer` корректно не наблюдает.

Отсюда исследовательский вопрос:

> Достаточен ли `single writer`, или внутреннее состояние Parallel Node требует `single owner + single observation point`?

Не исправлять это механически через pointer до определения семантики владения.

## Следующий эксперимент

Пока **не трогать ordering** и не вводить sequence ID, сортировку или reorder buffer.

Добиться:

* завершения pipeline по событию вместо `time.Sleep`;
* ровно 128 входов → 128 результатов;
* 64 результатов `0` и 64 результатов `1`;
* отсутствия data races;
* явного протокола завершения внутренних процессов Parallel Node.

## Зафиксировать

* Round-Robin сейчас означает **next worker**, а не **next available worker**.
* Smoother — внутренняя часть Parallel Node, а не самостоятельный топологический узел.
* Изменение `n` должно менять степень параллелизма, не меняя внешний контракт узла.
* Backpressure сейчас моделируется unbuffered channels; в распределённой версии он должен стать свойством протокола `Wire`.
* Известные локальные дефекты (`make([]byte, 256)`, незакрытый `Buffer.InChan`, неиспользуемый `Node.OutChan`, пустой worker list / nil channels) исправлять отдельно, не вводя скрытую семантику порядка.

---

## Comments


### @kshakirov — 2026-08-26T15:14:30Z

Дропнули первый кусок termination protocol. Workers теперь живут под `WaitGroup`; координатор ждёт завершения всех отправителей и только после этого закрывает `Buffer.InChan`. Канал закрывает тот, кто точно знает, что новых sends уже не будет.

Заодно убрали стартовые нули: slice теперь создаётся с `len=0`.

Зверь живёт: `go test ./bin` и `go vet ./bin` проходят.

Ordering не трогали. Следующая зацепка — завершение самого Smoother и единая точка наблюдения результата; копия slice header в `lsp.Buffer` пока остаётся открытым вопросом.


### @kshakirov — 2026-08-31T14:28:08Z

Дропнули следующий кусок живого тракта: `Smoother -> Sink`.

`LocalSplitBuffer` больше не растит slice: принимает результат worker через `InChan` и сразу отдаёт его в `OutChan`. После полного drain входа Smoother закрывает собственный выход; Sink только читает и в закрытие чужого канала не лезет. Все каналы пока unbuffered — backpressure честно долетает от Sink до Source, capacity отдельно не угадывали.

Зверь прошёл сквозную проверку: 128 входов -> 128 значений в Sink, 64 нуля / 64 единицы, `go test`, `go vet` и запуск под `-race` чистые.

Что осталось:
- перевести Source из конечной пачки в Eternal Streaming;
- заменить `time.Sleep` явным lifecycle: вечная работа в нормальном режиме и отдельный проверяемый shutdown;
- проверить длинные паузы: empty != closed;
- позже отдельно решить capacity каждого участка, не ломая backpressure;
- после стабилизации убрать мёртвые `OutChan`, `Dump` и прочие хвосты прототипа.

Ordering, sequence ID и distributed Wire пока не трогаем.
