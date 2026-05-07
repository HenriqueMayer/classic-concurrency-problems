# Problema 2: O Jogo do Dorminhoco
> Implementar o jogo de cartas Dorminhoco com N jogadores dispostos em anel.

### Regras do jogo:
- Cada jogador recebe uma mão de cartas. O objetivo é formar uma trinca (3 cartas iguais).
- A cada rodada, cada jogador descarta uma carta para o vizinho da esquerda e recebe uma carta do vizinho da direita.
- Quando um jogador forma uma trinca, ele “bate” (sinaliza). Os demais jogadores devem reagir o mais rápido possível — o último a reagir perde a rodada.

### Requisitos:
- [ ] Cada jogador é uma goroutine. A comunicação entre jogadores vizinhos é feita por channels.
- [ ] O programa deve imprimir o andamento da partida (trocas de cartas, quem bateu, quem perdeu).
- [ ] Identificar os pontos da implementação onde deadlock poderia ocorrer e justificar como a arquitetura escolhida os previne.