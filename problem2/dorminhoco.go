package main

import (
	"fmt"
	"math/rand"
	"time"
)

const numJogadores = 5

type Jogador struct {
	// Estrutura para definir o Jogador

	id  int   // id para identificação
	mao []int // lista para salvar as cartas
}

func distribuirCartas(numJogadores int) [][]int {
	// Função para distribuir as cartas para os jogadores
	//
	// Args:
	// 		numJogadores (int): número de jogadores
	//
	// Return:
	// 		list[int]

	// Cria um baralho vazio
	baralho := []int{}

	// Adicione três cartas do mesmo valor ao baralho
	for i := range numJogadores {
		baralho = append(baralho, i, i, i)
	}

	// Embaralha o baralho
	rand.Shuffle(len(baralho), func(i, j int) {
		baralho[i], baralho[j] = baralho[j], baralho[i]
	})

	// Distribui as cartas para cada jogador
	maos := make([][]int, numJogadores)
	for i := range numJogadores {
		maos[i] = baralho[i*3 : (i+1)*3]
	}

	return maos

}

func jogar(
	jogador Jogador,
	recebeCarta <-chan int,
	enviaCarta chan<- int,
	mesa chan<- int,
	alguemBateu <-chan bool) {
	// Função para simular a jogada de cada jogador
	//
	// Args:
	//		jogador (Jogador): objeto da classe Jogador
	//		recebeCarta (chan): canal para receber carta do jogador à esquerda
	//		enviaCarta (chan): canal para enviar carta ao jogador à direita
	//		mesa (chan): canal indicar à main que o jogador bateu
	//
	// Return:
	//		None

	id := jogador.id
	mao := jogador.mao

	for {
		// 1) Verificar se possui trinca
		if (mao[0] == mao[1]) && (mao[1] == mao[2]) {
			fmt.Printf("Jogador %d | Bateu! \n", id)
			mesa <- id
			return
		}

		// 2) Descartar uma carta ou bater
		var descarte int

		// Se tiver duas cartas iguais, descarta a carta diferente
		// Tentativa de simular um jogador real
		// Também, evitar uma partida infinitamente grande
		if mao[0] == mao[1] {
			descarte = mao[2]
			mao = []int{mao[0], mao[1]}
		} else if mao[0] == mao[2] {
			descarte = mao[1]
			mao = []int{mao[0], mao[2]}
		} else if mao[1] == mao[2] {
			descarte = mao[0]
			mao = []int{mao[1], mao[2]}
		} else {
			// Se todas cartas são diferentes, descarta a primeira
			descarte = mao[0]
			mao = []int{mao[1], mao[2]}
		}

		select {
		// Descar uma carta
		case enviaCarta <- descarte:
			fmt.Printf("Jogador %d | Descartou uma carta (%d). \n", id, descarte)

		// Se alguém bateu, bate também
		case <-alguemBateu:
			fmt.Printf("Jogador %d | Bateu em seguida (durante o descarte)!\n", id)
			time.Sleep(time.Duration(rand.Intn(1000)) * time.Millisecond)
			mesa <- id
			return
		}

		// 3) Pegar uma nova carta ou Bater
		select {
		// Comprar um carta
		case novaCarta := <-recebeCarta:
			mao = append(mao, novaCarta)
			fmt.Printf("Jogador %d | Recebeu uma carta (%d). \n", id, novaCarta)

		// Se alguém bateu, bate também
		case <-alguemBateu:
			fmt.Printf("Jogador %d | Bateu em seguida (durante a compra)!\n", id)
			time.Sleep(time.Duration(rand.Intn(1000)) * time.Millisecond)
			mesa <- id
			return
		}

	}

}

func main() {
	mesa := make(chan int, numJogadores)
	alguemBateu := make(chan bool)

	canais := make([]chan int, numJogadores)
	// https://stackoverflow.com/questions/21950244/is-there-a-way-to-iterate-over-a-range-of-integers
	for i := range numJogadores {
		canais[i] = make(chan int, 1)
	}

	fmt.Printf("Iniciando partida com %d jogadores\n", numJogadores)

	fmt.Println("Distribuindo cartas...")
	maosIniciais := distribuirCartas(numJogadores)

	for i := range numJogadores {

		novoJogador := Jogador{
			id:  i + 1,
			mao: maosIniciais[i],
		}
		fmt.Printf("Jogador %d | Mão inicial: %v \n", novoJogador.id, novoJogador.mao)

		recebeCarta := canais[i]
		enviaCarta := canais[(i+1)%numJogadores]

		go jogar(novoJogador, recebeCarta, enviaCarta, mesa, alguemBateu)
	}

	// Primeiro a bater é o vencedor
	vencedor := <-mesa

	// Sinaliza para as go routines que alguem bateu
	close(alguemBateu)

	// Último a bater é o perdedor
	var perdedor int
	for i := 0; i < numJogadores-1; i++ {
		perdedor = <-mesa
	}

	fmt.Printf("\n=====================\n")
	fmt.Printf("Jogo Encerrado! \n")
	fmt.Printf("Vencedor: Jogador %d \n", vencedor)
	fmt.Printf("Dorminhoco: Jogador %d\n", perdedor)
	fmt.Printf("=====================\n")

}
