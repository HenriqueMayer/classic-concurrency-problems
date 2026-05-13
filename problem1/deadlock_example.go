package main

import (
	"fmt"
	"sync"
	"time"
)

func mesa(id int, garfoEsquerda, garfoDireita chan struct{}, wg *sync.WaitGroup) {
	defer wg.Done()
	for i := 0; i < comida; i++ {
		fmt.Printf("Filósofo %d escolheu.\n", id)

		//cada Garfo é um canal com 1 espaço
		<-garfoEsquerda
		fmt.Printf("Filosofo %d pegou o garfo da esquerda.\n", id)
		time.Sleep(time.Millisecond * 200)

		<-garfoDireita
		fmt.Printf("Filosofo %d pegou o garfo da direita.\n", id)
		time.Sleep(time.Millisecond * 200)

		//liberamos os garfos
		garfoDireita <- struct{}{}
		garfoEsquerda <- struct{}{}

		fmt.Printf("Filosofo %d terminou de comer e largou os garfos.\n", id)
	}
}

func main() {

	//cria uma fila
	var wg sync.WaitGroup
	garfo := make([]chan struct{}, filosofos)

	//cria todos os garfos com um buffer de tamanho 1 e coloca 1 valor vazio no garfo
	for i := 0; i < filosofos; i++ {
		garfo[i] = make(chan struct{}, 1)
		garfo[i] <- struct{}{}
	}

	//cria uma fila com 5 pessoas
	wg.Add(filosofos)

	//cria os filosofos
	for i := 0; i < filosofos; i++ {
		go mesa(i, garfo[i], garfo[(i+1)%filosofos], &wg)
	}

	//pausa espera todos os filosofos comerem
	wg.Wait()
	fmt.Println("Todos os filosofos terminaram de comer.")
}

const (
	//numero de filosofos e quantidade de comidas
	filosofos = 5
	comida    = 3
)
