package main

import (
	"fmt"
	"sync"
	"time"
)

const (
	NumFilosofos = 5
	NumIteracoes = 1000
)

func main() {
	// 1. VERSÃO COM DEADLOCK
	fmt.Println("\n=== Iniciando Versão com DEADLOCK ===")
	fmt.Println("(O programa deve travar aqui. Aperte Ctrl+C para interromper)")
	// Importante: Esta versão vai travar o terminal. 
	// executarJantar("deadlock")

	// 2. ESTRATÉGIA DE HIERARQUIA
	fmt.Println("\n=== Iniciando Estratégia de Hierarquia (Dijkstra) ===")
	executarJantar("hierarquia")

	// 3. ESTRATÉGIA DE ÁRBITRO
	fmt.Println("\n=== Iniciando Estratégia de Árbitro (Limitador) ===")
	executarJantar("arbitro")
}

func executarJantar(estrategia string) {
	var wg sync.WaitGroup
	counts := make([]int, NumFilosofos)
	garfos := make([]chan bool, NumFilosofos)

	for i := 0; i < NumFilosofos; i++ {
		garfos[i] = make(chan bool, 1)
		garfos[i] <- true
	}

	arbitro := make(chan bool, NumFilosofos-1)

	wg.Add(NumFilosofos)
	for i := 0; i < NumFilosofos; i++ {
		esq := garfos[i]
		dir := garfos[(i+1)%NumFilosofos]

		switch estrategia {
		case "deadlock":
			go filosofoDeadlock(i, esq, dir, &wg)
		case "hierarquia":
			go filosofoHierarquia(i, esq, dir, &wg, counts)
		case "arbitro":
			go filosofoArbitro(i, esq, dir, arbitro, &wg, counts)
		}
	}

	// No caso do deadlock, o Wait() nunca será liberado.
	wg.Wait()

	if estrategia != "deadlock" {
		fmt.Println("Resultados (Vezes que cada um comeu):")
		for id, qtd := range counts {
			fmt.Printf("Filósofo %d: %d vezes\n", id, qtd)
		}
	}
}

// --- ESTRATÉGIA COM DEADLOCK ---
func filosofoDeadlock(id int, esq, dir chan bool, wg *sync.WaitGroup) {
	defer wg.Done()
	for {
		fmt.Printf("Filósofo %d tentando pegar garfo ESQUERDO...\n", id)
		<-esq 
		fmt.Printf("Filósofo %d pegou ESQUERDO. Tentando DIREITO...\n", id)
		
		// Todos ficarão travados aqui esperando o canal 'dir'
		<-dir 
		
		fmt.Printf("Filósofo %d COMENDO.\n", id)
		esq <- true
		dir <- true
	}
}

// --- ESTRATÉGIA DE HIERARQUIA ---
func filosofoHierarquia(id int, esq, dir chan bool, wg *sync.WaitGroup, counts []int) {
	defer wg.Done()
	for i := 0; i < NumIteracoes; i++ {
		primeiro, segundo := esq, dir
		if id == NumFilosofos-1 {
			primeiro, segundo = dir, esq
		}
		<-primeiro
		<-segundo
		counts[id]++
		primeiro <- true
		segundo <- true
		time.Sleep(time.Millisecond * 5)
	}
}

// --- ESTRATÉGIA DE ÁRBITRO ---
func filosofoArbitro(id int, esq, dir chan bool, arbitro chan bool, wg *sync.WaitGroup, counts []int) {
	defer wg.Done()
	for i := 0; i < NumIteracoes; i++ {
		arbitro <- true 
		<-esq
		<-dir
		counts[id]++
		esq <- true
		dir <- true
		<-arbitro 
		time.Sleep(time.Millisecond * 5)
	}
}