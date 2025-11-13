#!/bin/bash

APP_NAME="app"
LOG_FILE="app.log"

# Fungsi untuk mencari PID dari app
get_pid() {
  pgrep -f "./$APP_NAME"
}

# Fungsi untuk menghentikan app
stop_app() {
  PID=$(get_pid)
  if [ -n "$PID" ]; then
    echo "🔴 Menghentikan $APP_NAME (PID: $PID)..."
    kill "$PID"
    sleep 2
    if pgrep -f "./$APP_NAME" > /dev/null; then
      echo "❗ Gagal menghentikan $APP_NAME, mencoba kill -9..."
      kill -9 "$PID"
    else
      echo "✅ $APP_NAME berhasil dihentikan."
    fi
  else
    echo "ℹ️ Tidak ada proses $APP_NAME yang berjalan."
  fi
}

# Fungsi untuk menjalankan app
start_app() {
  echo "🚀 Menjalankan $APP_NAME..."
  nohup ./"$APP_NAME" > "$LOG_FILE" 2>&1 &
  echo "✅ $APP_NAME berhasil dijalankan dengan PID $(get_pid)"
}

# Fungsi untuk restart
restart_app() {
  stop_app
  sleep 2
  start_app
}

# Menangani argumen
case "$1" in
  start)
    start_app
    ;;
  stop)
    stop_app
    ;;
  restart)
    restart_app
    ;;
  status)
    PID=$(get_pid)
    if [ -n "$PID" ]; then
      echo "🟢 $APP_NAME sedang berjalan (PID: $PID)"
    else
      echo "🔴 $APP_NAME tidak berjalan."
    fi
    ;;
  *)
    echo "Gunakan: $0 {start|stop|restart|status}"
    exit 1
    ;;
esac
