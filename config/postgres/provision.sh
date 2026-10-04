#!/usr/bin/env bash
# Run inside the shared PostgreSQL accessory as an administrator.
set -euo pipefail

app_name=${1:-}
if [[ ! "$app_name" =~ ^[a-z][a-z0-9_]{0,62}$ ]] ||
   [[ "$app_name" == postgres || "$app_name" == template0 || "$app_name" == template1 ]]; then
  echo "Usage: bash /opt/postgres/provision.sh <service_name> (lowercase letters, digits, underscores)" >&2
  exit 1
fi

if [[ -t 0 ]]; then
  IFS= read -r -s -p "Password for ${app_name}: " app_password
  echo
  IFS= read -r -s -p "Confirm password: " confirmation
  echo
  if [[ "$app_password" != "$confirmation" ]]; then
    echo "Passwords do not match." >&2
    exit 1
  fi
else
  # Allows local verification without passing a password in command arguments.
  IFS= read -r app_password
fi
if [[ ${#app_password} -lt 16 ]]; then
  echo "Use a password of at least 16 characters." >&2
  exit 1
fi
export APP_DB_NAME="$app_name" APP_DB_PASSWORD="$app_password"

# psql's identifier/literal quoting handles special characters safely. Fail on an
# existing login/database instead of changing credentials or taking ownership.
psql --no-psqlrc --set=ON_ERROR_STOP=1 --username="${POSTGRES_USER:-postgres}" --dbname=postgres <<'SQL'
\getenv app_name APP_DB_NAME
\getenv app_password APP_DB_PASSWORD
CREATE ROLE :"app_name" LOGIN NOSUPERUSER NOCREATEDB NOCREATEROLE NOREPLICATION PASSWORD :'app_password';
CREATE DATABASE :"app_name" OWNER :"app_name";
REVOKE ALL ON DATABASE :"app_name" FROM PUBLIC;
GRANT CONNECT, TEMPORARY ON DATABASE :"app_name" TO :"app_name";
\connect :app_name
REVOKE ALL ON SCHEMA public FROM PUBLIC;
GRANT ALL ON SCHEMA public TO :"app_name";
SQL
unset APP_DB_PASSWORD app_password
echo "Created database and login: ${app_name}"
