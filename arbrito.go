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
	arbitro := make(chan bool, NumFilosofos-1) // Limita a 4 pessoas na mesa

	for i := 0; i < NumFilosofos; i++ {
		garfos[i] = make(chan bool, 1)
		garfos[i] <- true
	}

	wg.Add(NumFilosofos)
	for i := 0; i < NumFilosofos; i++ {
		go func(id int) {
			defer wg.Done()
			esq, dir := garfos[id], garfos[(id+1)%NumFilosofos]

			for j := 0; j < Iteracoes; j++ {
				arbitro <- true // Pede permissão ao árbitro
				<-esq
				<-dir
				fmt.Printf("Filósofo %d comeu (iteração %d)\n", id, j)
				esq <- true
				dir <- true
				<-arbitro // Libera vaga na mesa
				time.Sleep(time.Millisecond * 10)
			}
		}(i)
	}
	wg.Wait()
}