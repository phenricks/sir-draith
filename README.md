# Sir Draith - Bot Medieval para Discord

Bot de Discord interativo e expansível ambientado em um universo medieval sombrio, desenvolvido em Go.

## 🎯 Visão Geral

Sir Draith é um bot que combina funcionalidades de gestão de servidor Discord com elementos de RPG medieval, oferecendo uma experiência única e imersiva para comunidades. Ele atua como mestre de RPG e administrador, permitindo que usuários participem de aventuras, batalhas com cartas, e gerenciem seus personagens.

## 🚀 Funcionalidades Principais

- Sistema de personagens com classes e profissões
- Batalhas no estilo card game
- Narrativa automática e ambientação medieval
- Sistema de roles e permissões temático
- Gestão de servidor com temática medieval
- Sistema de missões e recompensas

## 🛠️ Tecnologias

- Go 1.21+
- DiscordGo
- MongoDB
- Docker

## 📋 Pré-requisitos

- Go 1.21 ou superior
- Docker e Docker Compose
- MongoDB
- Token de bot do Discord

## 🔧 Instalação

1. Clone o repositório:
   ```bash
   git clone https://github.com/phenricks/sir-draith.git
   cd sir-draith
   ```

2. Copie o arquivo de configuração de exemplo:
   ```bash
   cp configs/config.example.yaml configs/config.yaml
   ```

3. Configure suas variáveis no arquivo `configs/config.yaml`

4. Execute com Docker:
   ```bash
   make docker-build
   make docker-run
   ```

   Ou localmente:
   ```bash
   make build
   make run
   ```

## 🔨 Comandos Disponíveis

Use `make help` para ver todos os comandos disponíveis:

- `make build` - Compila o projeto
- `make run` - Executa o bot
- `make test` - Executa os testes
- `make lint` - Executa o linter
- `make docker-build` - Constrói a imagem Docker
- `make docker-run` - Executa o container Docker

## 📚 Documentação

A documentação completa está disponível na pasta `docs/`:

- [Guia de Contribuição](docs/CONTRIBUTING.md)
- [Documentação da API](docs/API.md)
- [Guia de Desenvolvimento](docs/DEVELOPMENT.md)

## 🤝 Contribuindo

1. Fork o projeto
2. Crie sua branch de feature (`git checkout -b feature/AmazingFeature`)
3. Commit suas mudanças (`git commit -m 'Add some AmazingFeature'`)
4. Push para a branch (`git push origin feature/AmazingFeature`)
5. Abra um Pull Request

## 📝 Licença

Este projeto está sob a licença MIT - veja o arquivo [LICENSE](LICENSE) para detalhes.

## ✨ Agradecimentos

- Comunidade Go
- Desenvolvedores do DiscordGo
- Todos os contribuidores

## 📞 Contato

Henrique - [@phenricks](https://github.com/phenricks) 