package closer

import (
	"log"
	"os"
	"os/signal"
	"sync"
)

//closer - сущность, которая копит функции закрытия каких-либо ресурсов и вызывается в конце

var globalCloser = New()

// Add adds `func() error` callback to the globalCloser
func Add(f ...func() error) {
	globalCloser.Add(f...)
}

// Wait ...
func Wait() {
	globalCloser.Wait()
}

// CloseAll ...
func CloseAll() {
	globalCloser.CloseAll()
}

type Closer struct {
	mu sync.Mutex // Мьютекс - это примитив синхронизации,
	// обеспечивающий взаимное исключение исполнения критических участков кода.

	once sync.Once //Once — это конструкция низкоуровневой синхронизации в Golang,
	// которая позволяет определить задачу для однократного выполнения за всё время работы программы.
	//Она содержит одну-единственную функцию Do, позволяющую передавать другую функцию для однократного применения.

	done  chan struct{}
	funcs []func() error
}

// New reurnts new Closer, if []os.Signal is specified Closer will automatically call CloseALL when one of signals is received from os
func New(sig ...os.Signal) *Closer { // принимает переменное количество сигналов
	c := &Closer{done: make(chan struct{})}
	if len(sig) > 0 {
		go func() {
			ch := make(chan os.Signal, 1)
			signal.Notify(ch, sig...) // ловим сигналы
			<-ch
			signal.Stop(ch)
			c.CloseAll() // метод вызова всех функций закрытия
		}()
	}

	return c
}

// Add func to closer
func (c *Closer) Add(f ...func() error) {
	c.mu.Lock() // блокируемся на  мьютексе
	c.funcs = append(c.funcs, f...)
	c.mu.Unlock()
}

// Wait blocks until all closer functions are done - Wait висит и ждет, пока что-нибудь не прилетит в канал
func (c *Closer) Wait() {
	<-c.done
}

// CloseAll calls all closer functions
func (c *Closer) CloseAll() {
	c.once.Do(func() { // т.к. с.CloseAll - выполняется ТОЛЬКО ОДИН РАЗ
		defer close(c.done)

		c.mu.Lock()
		funcs := c.funcs
		c.funcs = nil
		c.mu.Unlock()

		// call all Closer funcs async
		errs := make(chan error, len(funcs))
		for _, f := range funcs {
			go func(f func() error) {
				errs <- f() // собираем ошибки при выполнении функций в канал
			}(f)
		}

		for i := 0; i < cap(errs); i++ {
			if err := <-errs; err != nil { // идем по каналу и если есть ошибки - выводим в logg'и
				log.Printf("error returned from Closer: %v\n", err)
			}
		}
	})
}
