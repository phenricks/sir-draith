#!/bin/bash

# Configurações
BACKUP_DIR="/backup/mongodb"
DATE=$(date +%Y%m%d_%H%M%S)
MONGO_HOST="mongodb"
MONGO_PORT="27017"
MONGO_DATABASE="sir_draith"
MONGO_USER="$MONGO_ROOT_USER"
MONGO_PASS="$MONGO_ROOT_PASSWORD"

# Criar diretório de backup se não existir
mkdir -p $BACKUP_DIR

# Realizar backup
mongodump \
  --host $MONGO_HOST \
  --port $MONGO_PORT \
  --db $MONGO_DATABASE \
  --username $MONGO_USER \
  --password $MONGO_PASS \
  --authenticationDatabase admin \
  --out $BACKUP_DIR/$DATE

# Compactar backup
cd $BACKUP_DIR
tar -czf $DATE.tar.gz $DATE
rm -rf $DATE

# Manter apenas os últimos 7 backups
find $BACKUP_DIR -name "*.tar.gz" -type f -mtime +7 -delete

echo "Backup completed: $BACKUP_DIR/$DATE.tar.gz" 