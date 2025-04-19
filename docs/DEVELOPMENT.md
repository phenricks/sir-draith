# Guia de Desenvolvimento

Este guia fornece instruções detalhadas para configurar e trabalhar no ambiente de desenvolvimento do Sir Draith.

## 🛠️ Pré-requisitos

- Docker e Docker Compose
- Go 1.21 ou superior (para desenvolvimento local)
- Make
- Git

## 🚀 Configuração Inicial

1. Clone o repositório:
   ```bash
   git clone https://github.com/phenricks/sir-draith.git
   cd sir-draith
   ```

2. Configure as variáveis de ambiente:
   ```bash
   cp .env.example .env
   # Edite o arquivo .env com suas configurações
   ```

3. Configure o arquivo de configuração:
   ```bash
   cp configs/config.example.yaml configs/config.yaml
   # Edite o arquivo config.yaml conforme necessário
   ```

## 🐳 Desenvolvimento com Docker

### Iniciar o Ambiente
```bash
docker compose up -d
```

### Verificar Logs
```bash
docker compose logs -f bot        # Logs do bot
docker compose logs -f mongodb    # Logs do MongoDB
docker compose logs -f mongo-express  # Logs do Mongo Express
```

### Acessar Serviços
- MongoDB: localhost:27017
- Mongo Express: http://localhost:8081
- Bot: Através do Discord

### Parar o Ambiente
```bash
docker compose down  # Mantém os volumes
docker compose down -v  # Remove os volumes
```

## 💻 Desenvolvimento Local

### Instalar Dependências
```bash
make deps
```

### Hot Reload (Desenvolvimento)
```bash
make watch
```

### Executar Testes
```bash
make test
```

### Executar Linter
```bash
make lint
```

## 📦 MongoDB

### Backup Manual
```bash
./deployments/mongodb/backup.sh
```

### Restaurar Backup
```bash
# Substitua YYYYMMDD_HHMMSS pela data do backup
tar -xzf /backup/mongodb/YYYYMMDD_HHMMSS.tar.gz
mongorestore --uri "mongodb://user:pass@localhost:27017" YYYYMMDD_HHMMSS/
```

## 🔍 Troubleshooting

### Problemas Comuns

1. **Erro de Conexão com MongoDB**
   - Verifique se o serviço está rodando: `docker compose ps`
   - Verifique as credenciais no .env
   - Verifique os logs: `docker compose logs mongodb`

2. **Bot Não Conecta ao Discord**
   - Verifique o token no arquivo .env
   - Verifique os logs: `docker compose logs bot`
   - Verifique se o bot está habilitado no Discord Developer Portal

3. **Hot Reload Não Funciona**
   - Verifique se o Air está instalado
   - Verifique o arquivo .air.toml
   - Verifique permissões de arquivo

## 📝 Convenções

### Commits
- feat: Nova funcionalidade
- fix: Correção de bug
- docs: Documentação
- style: Formatação
- refactor: Refatoração
- test: Testes
- chore: Manutenção

### Branches
- main: Produção
- develop: Desenvolvimento
- feature/*: Novas funcionalidades
- fix/*: Correções
- release/*: Preparação para release

## 🤝 Contribuindo

1. Crie uma branch para sua feature
2. Desenvolva e teste suas mudanças
3. Execute os testes: `make test`
4. Execute o linter: `make lint`
5. Commit suas mudanças
6. Push para sua branch
7. Abra um Pull Request

## 📚 Recursos Adicionais

- [Documentação do DiscordGo](https://pkg.go.dev/github.com/bwmarrin/discordgo)
- [Documentação do MongoDB Go Driver](https://pkg.go.dev/go.mongodb.org/mongo-driver/mongo)
- [Docker Compose Reference](https://docs.docker.com/compose/reference/) 