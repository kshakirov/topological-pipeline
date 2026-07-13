package main

import (
	"fmt"
	_ "gostorm.org/go_storm/lib/spout"
	_ "gostorm.org/go_storm/lib/bolt"
	_ "gostorm.org/go_storm/lib/tuple"

	_ "sync"
	"time"
)

// Data — универсальный объект нашей категории (пока просто interface{})
type Data interface{}

// BoxFunc — чистая пользовательская функция. 
// Принимает объект, возвращает трансформированный объект.
type BoxFunc func(in Data) Data

// Box — инфраструктурный контейнер (Универсальный Узел)
type Box struct {
	ID          string
	UserFunc    BoxFunc
	
	// Внутренние каналы для физики потоков
}

func NewBox(id string, fn BoxFunc) *Box {
	return &Box{
		ID:          id,
		UserFunc:    fn,

	}
}

// Start запускает схемотехнику Ящика
func (b *Box) Start(msg Data ) Data{
	// Создаем индивидуальные каналы для каждого воркера

	
	res := b.UserFunc(msg)
	// Кто первый обработал — тот сразу плюет в Mux (outChan)
	return res

	// 2. Входной Demux (Строгий Раунд-Робин)

}

func main() {
	//	var wg sync.WaitGroup


	// Наша чистая функция: имитирует тяжелую математику (тугодумы vs шустрики)
	heavyCompute := func(in Data) Data {
		val := in.(int)
		// Искусственный разброс по времени, чтобы проверить FIFO на выходе
		if val%2 == 0 {
			time.Sleep(50 * time.Millisecond) // Четные — тугодумы
		} else {
			time.Sleep(10 * time.Millisecond) // Нечетные — шустрики
		}
		return fmt.Sprintf("Result(%d)", val*10)
	}

	box := NewBox("QuantumProcessor", heavyCompute)


	// Эмулируем подачу данных в Ящик
	
	for i := 1; i <= 6; i++ {
		fmt.Printf("📥 Подано на вход: %d\n", i)
		res:= box.Start(i)
		fmt.Printf("Получено со входа %s\n",res)
	}

	
	//	wg.Wait()

	// Читаем из выходной трубы (Mux)
	fmt.Println("🏁 Конвейер полностью остановлен.")
}
