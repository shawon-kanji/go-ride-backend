#!/usr/bin/env bash
# Stops every go-ride backend process started by start-servers.sh. Also
# cleans up an instance found bound to a service's port even if it wasn't
# started by this script (e.g. launched manually in another terminal).
# Postgres/Kafka/Redis containers are left running; pass --infra to stop
# those too.
set -euo pipefail

source "$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)/_servers-common.sh"

stop_pid() {
  local pid="$1" name="$2"
  echo "Stopping $name (pid $pid)..."
  pkill -TERM -P "$pid" 2>/dev/null || true
  kill -TERM "$pid" 2>/dev/null || true
  for _ in $(seq 1 10); do
    pid_alive "$pid" || return 0
    sleep 0.3
  done
  if pid_alive "$pid"; then
    echo "  still alive, force killing"
    pkill -KILL -P "$pid" 2>/dev/null || true
    kill -KILL "$pid" 2>/dev/null || true
  fi
}

echo "$SERVICES" | while IFS='|' read -r name reldir runcmd port gowork; do
  [ -z "$name" ] && continue
  pidfile="$RUN_DIR/$name.pid"

  if [ -f "$pidfile" ]; then
    pid="$(cat "$pidfile")"
    if pid_alive "$pid"; then
      stop_pid "$pid" "$name"
    else
      echo "$name: not running (stale pid file)"
    fi
    rm -f "$pidfile"
  elif [ -n "$port" ]; then
    extpid="$(port_owner_pid "$port")"
    if [ -n "$extpid" ]; then
      echo "$name: no pid file, but found untracked process on port $port"
      stop_pid "$extpid" "$name"
    else
      echo "$name: not running"
    fi
  else
    echo "$name: not running"
  fi
done

if [ "${1:-}" = "--infra" ]; then
  echo "== Infrastructure =="
  (cd "$PROJECTS_DIR/go-ride-kafka-consumers" && docker compose down)
  (cd "$PROJECTS_DIR/go-ride-backend" && docker compose down)
else
  echo "(Postgres/Kafka/Redis containers left running; pass --infra to stop those too)"
fi

print_status
