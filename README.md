Este repositório faz parte de um conjunto de três repositórios onde implementarei a mesma aplicação com o intuito de comparar cada implementação em questões como:
- Legibilidade do código
- Performance
- Facilidade de manutenção e extensão

A outra versão do projeto será desenvolvidas em Zig. O intuito é utilizar o mínimo de depêndencias externas possível, atendo-se ao disponibilizado pela linguagem e focando em criar implementações próprias para problemas comuns.

Todo:
- tornar validateTodoTitleLen e outras validações mais genéricas, permitindo reutilizações pra validações em outras partes do projeto
- implementar conexão com db
- escrever os testes automatizados
- escrever sistema de aplicação de migrations: tabela no banco registrando as migrations aplicadas pelo nome do arquivo e data de aplicação. Talvez transformar isso em uma lib separada. Talvez escrever em Zig e usar nesse projeto e na versão de Zig
- sanitizar os inputs
- implementar auth com jwt
- implementar sistema de middlewares