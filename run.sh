#!/usr/bin/env bash
clear
set -euo pipefail

PORT="${PORT:-8080}"
ADDR=":${PORT}"
BASE_URL="http://localhost:${PORT}"

SERVER_PID=""

pp() { if command -v jq >/dev/null 2>&1; then jq .; else cat; fi; }

wait_ready() {
  local url="$1"
  for i in {1..40}; do
    if curl -fsS "${url}/healthz" >/dev/null 2>&1; then
      return 0
    fi
    sleep 0.25
  done
  echo "Server did not become ready at ${url}" >&2
  return 1
}

free_port() {
  # Kill anything bound to $PORT using lsof, fuser, ss, or netstat (best available)
  local killed=0
  if command -v lsof >/dev/null 2>&1; then
    if pids=$(lsof -ti :"${PORT}" 2>/dev/null); then
      if [ -n "${pids}" ]; then
        echo "Freeing port ${PORT} via lsof (PIDs: ${pids})..."
        kill -9 ${pids} 2>/dev/null || true
        killed=1
      fi
    fi
  fi
  if [ $killed -eq 0 ] && command -v fuser >/dev/null 2>&1; then
    echo "Freeing port ${PORT} via fuser..."
    fuser -k "${PORT}/tcp" 2>/dev/null || true
    killed=1
  fi
  if [ $killed -eq 0 ] && command -v ss >/dev/null 2>&1; then
    # Parse PID from ss output (may require root to see all PIDs)
    pid=$(ss -ltnp 2>/dev/null | awk -v p=":${PORT}" '$4 ~ p {print $NF}' | sed -nE 's/.*pid=([0-9]+).*/\1/p' | head -n1)
    if [ -n "${pid:-}" ]; then
      echo "Freeing port ${PORT} via ss (PID: ${pid})..."
      kill -9 "${pid}" 2>/dev/null || true
      killed=1
    fi
  fi
  if [ $killed -eq 0 ] && command -v netstat >/dev/null 2>&1; then
    pid=$(netstat -ltnp 2>/dev/null | awk -v p=":${PORT}" '$4 ~ p {print $NF}' | sed -nE 's#.*/([0-9]+)$#\1#p' | head -n1)
    if [ -n "${pid:-}" ]; then
      echo "Freeing port ${PORT} via netstat (PID: ${pid})..."
      kill -9 "${pid}" 2>/dev/null || true
      killed=1
    fi
  fi
}

start_server() {
  echo "→ Starting server on ${ADDR}"
  ADDR="${ADDR}" go run ./... >/dev/null 2>&1 &
  SERVER_PID=$!
  sleep 0.1
  wait_ready "${BASE_URL}"
  echo "  Server ready (pid ${SERVER_PID})"
}

stop_server() {
  if [[ -n "${SERVER_PID}" ]] && kill -0 "${SERVER_PID}" 2>/dev/null; then
    echo "→ Stopping server (pid ${SERVER_PID})"
    kill "${SERVER_PID}" 2>/dev/null || true
    for _ in {1..50}; do
      kill -0 "${SERVER_PID}" 2>/dev/null || { echo "  Server stopped"; SERVER_PID=""; break; }
      sleep 0.1
    done
    if [[ -n "${SERVER_PID}" ]]; then
      echo "  Force killing..."
      kill -9 "${SERVER_PID}" 2>/dev/null || true
      SERVER_PID=""
    fi
  fi
  # Ensure port is freed even if the process changed/forked
  free_port
}

cleanup() { stop_server || true; }
trap cleanup EXIT

# ── Ensure port free at the very start ─────────────────────────────────────────
free_port
sleep 0.1

echo "STEP 1: Start the server"
start_server

echo "STEP 2: Run tests"
go test -v

echo "STEP 3: Stop the server"
stop_server

echo "STEP 4: Start the server"
start_server

echo "STEP 5: Target health + endpoints"
echo "  /healthz"
curl -fsS "${BASE_URL}/healthz" && echo

echo "  /weather Los Angeles, CA"
curl -fsS "${BASE_URL}/weather?lat=34.05&lon=-118.25" | pp

echo "  /weather New York, NY"
curl -fsS "${BASE_URL}/weather?lat=40.7128&lon=-74.0060" | pp

echo "  /weather Phoenix, AZ"
curl -fsS "${BASE_URL}/weather?lat=33.4484&lon=-112.0740" | pp

echo "STEP 6: Stop the server"
stop_server

# ── Ensure port free at the very end ───────────────────────────────────────────
free_port
echo "Done"
