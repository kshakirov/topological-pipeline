# TCP Wire v1 — переход к распределённому Pipeline

## Контекст

Локальный эксперимент дал нам работающий Parallel Node:

```text
External Choice
      |
      v
Parallel Composition
   /   |   \
 W1   W2   ... Wn
   \   |   /
      v
   Smoother
```

Сейчас взаимодействие внутри него построено на Go channels.

Следующий этап — **не переписывать Parallel Node под TCP**, а вынести TCP только на границы между топологическими узлами.

---

## 1. Parallel Node атомарен

Parallel Node рассматривается как единый топологический узел.

Нельзя разместить часть его workers на одной машине, а другую часть — на другой.

Весь узел:

```text
External Choice
      +
Parallel Composition
      +
Smoother
```

работает локально на одной машине.

Количество workers `n` — внутренняя степень параллелизма узла.

---

## 2. Где появляется TCP

TCP нужен **только между атомарными топологическими узлами**.

Внутри Parallel Node пока ничего не меняем:

```text
                 MACHINE A

TCP
 |
 v
External Choice
 |
 | Go channel
 v
Parallel Composition
 |
 | Go channels
 v
Smoother
 |
 v
TCP
```

Таким образом:

> local communication = Go channels  
> inter-node communication = TCP

---

## 3. Что меняем

### Input

External Choice раньше получал данные непосредственно из Go channel.

Теперь его внешняя сторона должна получать TCP byte stream, декодировать из него ST frames и передавать payload во внутреннюю локальную часть.

```text
TCP bytes
   ↓
ST frame decoder
   ↓
payload
   ↓
External Choice / local dispatch
```

### Output

Smoother получает результаты workers как и раньше.

На внешней границе результат кодируется в ST frame и отправляется через TCP.

```text
workers
   ↓
Smoother
   ↓
ST frame encoder
   ↓
TCP bytes
```

Внутреннюю Parallel Composition на этом этапе **не трогаем**.

---

## 4. ST Wire v1

Первый протокол намеренно минимален.

```text
+-------------------+--------------------+------------------------+
|  Magic (2 bytes)  | Length (2 bytes)   |  Payload (N bytes)     |
|   0x53 0x54       | uint16 BigEndian   |  raw bytes             |
+-------------------+--------------------+------------------------+
```

То есть:

```text
ST | Length | Payload
```

### Magic

```text
0x53 0x54
'S'  'T'
```

Маркер начала ST frame.

### Length

2 bytes, unsigned integer, Big Endian.

```text
uint16
```

Определяет количество байт в Payload.

### Payload

Ровно `Length` сырых байт.

Протокол пока ничего не знает о внутренней семантике payload.

---

## 5. Важное свойство TCP

TCP — byte stream.

Один `Read()` **не равен одному frame**.

Один frame может прийти несколькими чтениями:

```text
Read #1: ST 00
Read #2: 05 AA BB
Read #3: CC DD EE
```

И несколько frames могут прийти одним чтением.

Поэтому External Choice должен рассматривать вход как поток байтов и распознавать:

```text
Magic → Length → Payload → Magic → Length → Payload → ...
```

Это естественное место для нашего потокового распознавателя/автомата.

---

## 6. Что пока НЕ делаем

На первом этапе не добавляем:

- sequence ID;
- ordering/reordering;
- checksum;
- version;
- message type;
- ACK;
- compression;
- distributed workers;
- сложный control protocol;
- autoscaling.

Добавляем новое поле только тогда, когда эксперимент покажет необходимость.

---

## Первый эксперимент

Минимальная цель:

```text
Node A
  Smoother
     |
 ST encoder
     |
    TCP
     |
 ST decoder
     |
External Choice
  Node B
```

Добиться передачи последовательности ST frames через реальный TCP stream без изменения внутренней логики Parallel Node.

Главный принцип этапа:

> **Распределяем узлы, а не внутренности узла.**

TCP Wire заменяет только внешнюю связь между атомарными Parallel Nodes. Проверенную локальную CSP-модель внутри узла сохраняем.