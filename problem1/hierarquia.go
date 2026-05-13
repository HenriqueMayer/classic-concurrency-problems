package main

import (
	"fmt"
	"sync"
	"time"
)

const (
	NumFilosofos = 5
	Iteracoes    = 10
)

func main() {
	var wg sync.WaitGroup
	garfos := make([]chan bool, NumFilosofos)

	for i := 0; i < NumFilosofos; i++ {
		garfos[i] = make(chan bool, 1)
		garfos[i] <- true
	}

	wg.Add(NumFilosofos)
	for i := 0; i < NumFilosofos; i++ {
		go func(id int) {
			defer wg.Done()
			primeiro, segundo := garfos[id], garfos[(id+1)%NumFilosofos]
			
			// Inversão de hierarquia para o último filósofo
			if id == NumFilosofos-1 {
				primeiro, segundo = segundo, primeiro
			}

			for j := 0; j < Iteracoes; j++ {
				<-primeiro
				<-segundo
				fmt.Printf("Filósofo %d comeu (iteração %d)\n", id, j)
				primeiro <- true
				segundo <- true
				time.Sleep(time.Millisecond * 10)
			}
		}(i)
	}
	wg.Wait()
}