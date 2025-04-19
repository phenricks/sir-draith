// Criar usuário para o banco de dados da aplicação
db.createUser({
    user: process.env.MONGO_APP_USER || 'sir_draith_user',
    pwd: process.env.MONGO_APP_PASSWORD || 'sir_draith_password',
    roles: [
        {
            role: 'readWrite',
            db: process.env.MONGO_DATABASE || 'sir_draith'
        }
    ]
});

// Criar coleções iniciais
db = db.getSiblingDB(process.env.MONGO_DATABASE || 'sir_draith');

// Coleção de personagens
db.createCollection('characters');
db.characters.createIndex({ "userId": 1 });
db.characters.createIndex({ "name": 1 }, { unique: true });
db.characters.createIndex({ "class": 1 });
db.characters.createIndex({ "level": 1 });

// Coleção de cartas
db.createCollection('cards');
db.cards.createIndex({ "name": 1 }, { unique: true });
db.cards.createIndex({ "type": 1 });
db.cards.createIndex({ "rarity": 1 });

// Coleção de eventos
db.createCollection('events');
db.events.createIndex({ "type": 1 });
db.events.createIndex({ "createdAt": 1 });

// Coleção de logs
db.createCollection('logs');
db.logs.createIndex({ "timestamp": 1 });
db.logs.createIndex({ "type": 1 });
db.logs.createIndex({ "userId": 1 });

// Criar TTL index para logs antigos (30 dias)
db.logs.createIndex({ "timestamp": 1 }, { expireAfterSeconds: 2592000 }); 