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
			esq, dir := garfos[id], garfos[(id+1)%NumFilosofos]
			for {
				fmt.Printf("Filósofo %d tentando comer...\n", id)
				<-esq // Pega esquerda
				<-dir // Pega direita (Vai travar aqui se todos pegarem a esquerda)
				fmt.Printf("Filósofo %d COMENDO.\n", id)
				esq <- true
				dir <- true
			}
		}(i)
	}
	wg.Wait()
}