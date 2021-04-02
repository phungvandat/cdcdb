#!/bin/bash
echo "host replication all 0.0.0.0/0 md5" >> "$PGDATA/pg_hba.conf"
set -e
psql -v ON_ERROR_STOP=1 --username "$POSTGRES_USER" --dbname "$POSTGRES_DB" <<-EOSQL
CREATE USER $PG_REP_USER REPLICATION LOGIN CONNECTION LIMIT 100 ENCRYPTED PASSWORD '$PG_REP_PASSWORD';
CREATE PUBLICATION $PG_PUBLICATION_NAME FOR ALL TABLES;
EOSQL
cat >> ${PGDATA}/postgresql.conf <<EOF
wal_level = logical
EOF