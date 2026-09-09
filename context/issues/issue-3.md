# Issue #3 — Перевод конвейера на IoC-интерфейс `Processor`


**State:** CLOSED
**Author:** @kshakirov
**Created:** 2026-08-03T16:44:31Z
**Updated:** 2026-08-12T14:11:09Z
**URL:** https://github.com/kshakirov/topological-pipeline/issues/3

---



## 1. Назначение задачи
Полностью устранить жесткую связанность (Coupling) между транспортным слоем (`Wire`) и бизнес-логикой. Перейти от передачи функциональных типов (`BoxFunc`) в метод `WireIn` к единому полиморфному контракту обратного вызова (Callback).

## 2. Критерии приемки (Definition of Done)
- [ ] В файле `types.go` объявить интерфейс `Processor` с методом `Process(in Set) Set`.
- [ ] Изменить сигнатуру метода `Wire.WireIn(p Processor)` в интерфейсе и его реализации `LocalWire`.
- [ ] Реализовать контракт `Processor` на структуре `ComputeNode` в файле `compute_node.go`.
- [ ] Инкапсулировать двухтактный цикл (`InBox` -> `OutBox`) внутри метода `ComputeNode.Process`.
- [ ] Изменить метод `ComputeNode.Prep()`, чтобы он передавал указатель на себя в провод: `s.Wire.WireIn(s)`.
- [ ] Убедиться, что `main.go` компилируется без ошибок и логи сквозной обработки остаются стабильными.


---

## Comments

