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
