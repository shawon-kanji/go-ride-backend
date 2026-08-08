# Shared state/helpers for start-servers.sh and stop-servers.sh.
# Not meant to be run directly.

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
ROOT_DIR="$(cd "$SCRIPT_DIR/.." && pwd)"
PROJECTS_DIR="$(cd "$ROOT_DIR/.." && pwd)"
RUN_DIR="$SCRIPT_DIR/.run"
mkdir -p "$RUN_DIR"

# name|dir (relative to PROJECTS_DIR)|go run subcommand|port (empty = none)|GOWORK=off needed
SERVICES='
go-ride-backend|go-ride-backend|./cmd/api|8080|no
location-producers|go-ride-kafka-consumers/services/location-producers|./cmd/api|8081|no
location-consumers|go-ride-kafka-consumers/services/location-consumers|./cmd/consumer||no
cab-request-handler|go-ride-kafka-consumers/services/cab-request-handler|./cmd/api|8082|no
driver-request-handler|go-ride-kafka-consumers/services/driver-request-handler|./cmd/api|8084|no
trip-dispatch-worker-api|go-ride-kafka-consumers/services/trip-dispatch-worker|./cmd/api||no
trip-dispatch-worker-consumer|go-ride-kafka-consumers/services/trip-dispatch-worker|./cmd/consumer||no
websocket-gateway|go-ride-kafka-consumers/services/websocket-gateway|./cmd/api|8083|no
driver-location-worker|go-ride-kafka-consumers|./cmd/driver-location-worker||yes
'

pid_alive() {
  kill -0 "$1" 2>/dev/null
}

port_owner_pid() {
  lsof -ti "tcp:$1" 2>/dev/null | head -1 || true
}

print_status() {
  printf '\n%-30s %-10s %-8s %s\n' "SERVICE" "STATUS" "PORT" "PID"
  echo "$SERVICES" | while IFS='|' read -r name reldir runcmd port gowork; do
    [ -z "$name" ] && continue
    pidfile="$RUN_DIR/$name.pid"
    pid=""
    status="stopped"
    if [ -f "$pidfile" ]; then
      pid="$(cat "$pidfile")"
      if pid_alive "$pid"; then status="running"; else pid=""; fi
    fi
    if [ "$status" = "stopped" ] && [ -n "$port" ]; then
      extpid="$(port_owner_pid "$port")"
      if [ -n "$extpid" ]; then status="running*"; pid="$extpid"; fi
    fi
    printf '%-30s %-10s %-8s %s\n' "$name" "$status" "${port:--}" "${pid:--}"
  done
  echo "(* = something is bound to that port but wasn't started by this script)"
}
