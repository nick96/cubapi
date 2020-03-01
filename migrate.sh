#!/usr/bin/env bash
# Run migrations using 'goose'

# Ensure the goose migration tool is installed.
command goose &>/dev/null
if [ $? -ne 0 ]
then
    go get -u github.com/pressly/goose/cmd/goose
fi

db="postgres"
conn_string="user=$DB_USER dbname=$DB_NAME password=$DB_PASS host=$DB_HOST sslmode=$DB_SSL_MODE"

exec goose "$db" "$conn_string" "$@"
