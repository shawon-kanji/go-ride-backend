#!/usr/bin/env bash
# Starts Postgres/Kafka/Redis, runs DB migrations, ensures Kafka topics exist,
# then starts every go-ride backend process (go-ride-backend + all
# go-ride-kafka-consumers services). Safe to re-run: already-running
# services (tracked or found bound to their port) are left alone.
set -euo pipefail

source "$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)/_servers-common.sh"

echo "== Infrastructure =="
(cd "$PROJECTS_DIR/go-ride-backend" && docker compose up -d)
(cd "$PROJECTS_DIR/go-ride-kafka-consumers" && docker compose up -d)

echo "Waiting for Postgres..."
tries=0
until (cd "$PROJECTS_DIR/go-ride-backend" && docker compose exec -T postgres pg_isready -U postgres -d go_ride) >/dev/null 2>&1; do
  tries=$((tries + 1))
  [ "$tries" -gt 30 ] && { echo "Postgres did not become ready in time" >&2; break; }
  sleep 1
done

echo "Waiting for Kafka..."
tries=0
until (cd "$PROJECTS_DIR/go-ride-kafka-consumers" && docker compose exec -T kafka kafka-topics.sh --list --bootstrap-server localhost:9092) >/dev/null 2>&1; do
  tries=$((tries + 1))
  [ "$tries" -gt 30 ] && { echo "Kafka did not become ready in time" >&2; break; }
  sleep 1
done

echo "== DB migrations =="
(
  cd "$PROJECTS_DIR/go-ride-db-schema"
  export DB_HOST=localhost DB_PORT=5432 DB_USER=postgres DB_PASSWORD=postgres DB_NAME=go_ride DB_SSLMODE=disable
  go run ./cmd/migrate up
)

echo "== Kafka topics =="
(cd "$PROJECTS_DIR/go-ride-kafka-consumers" && make topic-create-all)

echo "== App servers =="
echo "$SERVICES" | while IFS='|' read -r name reldir runcmd port gowork; do
  [ -z "$name" ] && continue
  dir="$PROJECTS_DIR/$reldir"
  pidfile="$RUN_DIR/$name.pid"
  logfile="$RUN_DIR/$name.log"

  if [ -f "$pidfile" ] && pid_alive "$(cat "$pidfile")"; then
    echo "$name: already running (pid $(cat "$pidfile")), skipping"
    continue
  fi

  if [ -n "$port" ]; then
    extpid="$(port_owner_pid "$port")"
    if [ -n "$extpid" ]; then
      echo "$name: port $port already in use by pid $extpid (not started by this script), skipping"
      continue
    fi
  fi

  if [ ! -f "$dir/.env" ]; then
    if [ -f "$dir/.env.example" ]; then
      cp "$dir/.env.example" "$dir/.env"
      echo "$name: created $dir/.env from .env.example"
    else
      echo "$name: no .env or .env.example in $dir, skipping" >&2
      continue
    fi
  fi

  echo "Starting $name..."
  (
    cd "$dir"
    set -a
    # shellcheck disable=SC1091
    source .env
    set +a
    if [ "$gowork" = "yes" ]; then export GOWORK=off; fi
    exec go run "$runcmd"
  ) >"$logfile" 2>&1 &
  pid=$!
  disown "$pid" 2>/dev/null || true
  echo "$pid" >"$pidfile"
  sleep 0.3
  if pid_alive "$pid"; then
    echo "  -> pid $pid, log: $logfile"
  else
    echo "  -> FAILED to start, see $logfile" >&2
  fi
done

print_status
