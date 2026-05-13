# Classic Concurrency Problems
> T1 - Fundamentos de Processamento Paralelo e Distribuído (98713-04)

### Integrantes
1. Henrique Mayer
2. Ighor Telles
3. Frederico Fragoso
4. João Henrique da Luz

### Opção A


### Estrutura
```
root/
    doc/
        pdf/t1-opcao-a.pdf
    
    problem1/
        doc.md --> documentação e controle de requisitos
    
    problem2/
        doc.md --> documentação e controle de requisitos
```

### Como executar o projeto
> TODO depois de finalziar o projeto

# Relatório

## Análise do Problema 1 (Jantar dos Filósofos)
---

### Análise de Deadlock e Starvation

#### Versão Base (Com Deadlock)
**Deadlock: Ocorre**
    - Como todos os filósofos seguem a mesma ordem (esquerda depois direita), satisfazemos a condição de Espera Circular. Se todos pegarem o garfo da esquerda simultaneamente, nenhum conseguirá o da direita, travando o sistema indefinidamente

**Starvation: Ocorre**
    - Uma vez que o sistema entra em deadlock, nenhum filósofo consegue progredir para o estado de "comendo", levando todos à inanição por bloqueio total.

---
#### Estratégia 1 (Hierarquia de Recursos)
**Deadlock: Não Ocorre**
    - Ao fazer o último filósofo inverter a ordem de acesso (direita depois esquerda), quebramos o ciclo de dependência. É impossível formar uma corrente circular de espera.

**Starvation: Não Ocorre**
    - O escalonador do Go (runtime scheduler) distribui o acesso aos canais de forma eficiente, garantindo que todos eventualmente consigam os dois garfos.

---
#### Estratégia 2 (Árbitro/Limitador)
**Deadlock: Não Ocorre**
    - Limitando a mesa a ***N-1*** comensais, garantimos que pelo menos um filósofo terá ambos os garfos disponíveis, pois haverá 5 garfos para 4 pessoas.


**Starvation: Não Ocorre**
    - O uso de um canal com buffer para o árbitro atua como uma fila, permitindo que todos os que estão esperando entrem na mesa conforme outros saem.

---
### Análise de Execução: O Jantar dos Filósofos
A tabela abaixo compara o comportamento das três implementações realizadas no código Go, focando no critério de ***Fairness (Justiça)*** e na ocorrência de travamentos.


| Estratégia            | Deadlock | Starvation | Fairness (Justiça) |
|-----------------------|-----------|-------------|--------------------|
| Versão Base           | Sim       | Sim         | Nula               |
| Hierarquia (Dijkstra) | Não       | Não         | Boa                |
| Árbitro (Limitador)   | Não       | Não         | Excelente          |

---

### Estratégia de Hierarquia (Dijkstra)

| Filósofo | Vezes que comeu |
|---|---|
| 0 | 1000 |
| 1 | 1000 |
| 2 | 1000 |
| 3 | 1000 |
| 4 | 1000 |

---

### Estratégia de Árbitro (Limitador)

| Filósofo | Vezes que comeu |
|---|---|
| 0 | 1000 |
| 1 | 1000 |
| 2 | 1000 |
| 3 | 1000 |
| 4 | 1000 |


---
### Explicação Detalhada do Fairness

#### 1. Versão Base (Deadlock)

- **Justiça pré-travamento:** Aleatória. Depende de qual *goroutine* o escalonador do Go executa primeiro.

- **Justiça pós-travamento:** Inexistente. Assim que o Deadlock ocorre, o contador de todos os filósofos para de subir.

- **Conclusão:** É a versão menos justa, pois a falha de sincronização impede que qualquer processo conclua sua tarefa, resultando em **Inanição Total**.

---

### 2. Hierarquia de Recursos

- **Comportamento:** Muito estável. Como a ordem de busca pelos garfos é alterada para o último filósofo, o sistema flui sem interrupções.

- **Dados de Execução:** Os contadores de refeições tendem a ser idênticos ou muito próximos (ex: todos terminam com exatamente 10 refeições), provando que ninguém foi "esquecido" pelo sistema.

---

### 3. Árbitro (Limitador de Comensais)

- **Comportamento:** É a estratégia mais robusta para garantir que todos tenham sua vez. O canal de controle (árbitro) funciona como uma fila de espera.

- **Dados de Execução:** Apresenta a maior equidade, pois a restrição de ocupação da mesa força uma rotatividade natural entre os filósofos pensantes e famintos.

---

## Análise do Problema 2 (Dorminhoco)

A implementação do jogo estruturada em topologia de anel apresenta riscos inerentes de interrupção ou travamento total se a concorrência não for gerenciada adequadamente. 

Abaixo estão listados três cenários em que há ocorrência de Deadlock:

### 1. Não utilizar canais bufferizados

    Causa: Em um canal não bufferizado (tamanho 0), a operação de envio bloqueia a execução imediatamente até que o vizinho receptor leia a mensagem. Se todos os N jogadores tentarem enviar suas cartas de descarte ao mesmo tempo, todos ficarão bloqueados aguardando a leitura. Como ninguém conseguirá avançar para a instrução de "recebimento", o sistema congela em uma espera circular irresolvível (deadlock).

    Solução: O uso de canais com buffer de tamanho 1 (make(chan int, 1)) permite que o envio ocorra de forma assíncrona, absorvendo a carta e liberando o jogador para realizar a leitura do próximo canal.

### 2. Não utilizar a estrutura select no envio e recebimento

    Causa: Quando um jogador forma a trinca, ele sinaliza a vitória e encerra sua rotina (return), abandonando permanentemente seu canal de leitura. Se o vizinho tentar repassar uma carta para esse jogador ausente através de um envio direto e bloqueante (enviaCarta <- descarte), ele ficará travado para sempre, pois o canal encherá e nunca mais será esvaziado. Esse travamento se propaga em cadeia (efeito cascata) para o resto da mesa, impedindo que os outros jogadores reajam ao fim do jogo.

    Solução: O uso do select cria rotas de fuga. Ao tentar enviar ou receber uma carta, a goroutine escuta simultaneamente um canal de aviso global (ex: alguemBateu). Se o jogo acabar durante uma transação, o select garante que o jogador abandone a operação de comunicação e reaja ao fim da partida.

### 3. Não utilizar a diretiva go ao chamar a função jogar()

    Causa: Omitir a palavra-chave go faz com que a função jogar() seja executada de forma totalmente síncrona, bloqueando a thread principal (main). Ao instanciar o primeiro jogador, a thread principal ficará presa dentro do loop infinito deste jogador aguardando que ele receba uma carta. Contudo, os demais jogadores sequer foram instanciados (o loop da main não avança para o índice 1), logo, a carta nunca será enviada. O programa paralisa instantaneamente em sua primeira iteração.

    Solução: A diretiva go garante que cada jogador seja lançado como uma goroutine independente, permitindo que a função main conclua o seu fluxo, instancie todos os membros da mesa e inicie o gerenciamento da partida.