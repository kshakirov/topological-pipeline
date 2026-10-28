# #8 — Go: Streaming Framer & DATA v1 Decoder

State: open
URL: https://github.com/kshakirov/topological-pipeline/issues/8

Epic: #7 — DATA Protocol & Distributed Acknowledgement.

## Задача

Реализовать потоковое чтение бинарных кадров в Go-узлах и декодирование DATA v1. ACK и генерация потомков вне этого тикета.

## Формат кадра

```text
Magic    uint16  "ST" (53 54)
Length   uint16  big-endian
Opcode   uint8
Body     Length - 1 bytes
```

Length включает Opcode + Body; полный кадр занимает 4 + Length байт.

## DATA v1 — Opcode 0x04

```text
Version   uint8 = 1
Flags     uint8
RootID    uint64 BE
MessageID uint64 BE
ParentID  uint64 BE
Payload   bytes
```

Для корневого сообщения ParentID = 0. Минимальный Length DATA v1 = 27; Payload — оставшиеся байты.

## Реализация

- [ ] Framer собирает полные кадры из произвольных TCP-фрагментов.
- [ ] Поддерживает несколько кадров в одном чтении.
- [ ] Сохраняет неполный хвост между чтениями.
- [ ] Проверяет Magic, Length и границы буфера.
- [ ] Decoder распознаёт DATA v1, выделяет метаданные и Payload.
- [ ] Передаёт декодированное сообщение обработчику узла.

## Ограничения

- Не изменять существующие Opcode 01–04.
- Не реализовывать ACK и генерацию потомков.
- Не менять маршрутизацию BEAM.
- Не вводить новую сериализацию Payload.

## Acceptance Criteria

- [ ] Один кадр, разбитый на произвольные TCP-фрагменты.
- [ ] Несколько кадров в одном TCP-чтении.
- [ ] Неполный последний кадр.
- [ ] Некорректные Magic, Length и Version обнаруживаются.
- [ ] Корректно извлекаются RootID, MessageID, ParentID и Payload.
- [ ] Существующий тракт Go → BEAM → Go продолжает работать.

## Перед реализацией

Определить политику при неверных Magic/Length/Version: завершение соединения или восстановление потока. Семантику Flags пока не расширяем.

Проверить совместимость с текущим relay: прежний BEAM DATA-handler принимает однобайтовое тело. DATA v1 требует прозрачной передачи всего тела; это нужно сверить до реализации, не подменяя ограничение «не менять маршрутизацию BEAM» скрытым расширением задачи.


