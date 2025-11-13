#!/bin/bash

APP_NAME="app"
LOG_FILE="app.log"

# Function to get the PID of the app
get_pid() {
  pgrep -f "./$APP_NAME"
}

# Function to stop the app
stop_app() {
  PID=$(get_pid)
  if [ -n "$PID" ]; then
    echo "🔴 Stopping $APP_NAME (PID: $PID)..."
    kill "$PID"
    sleep 2
    if pgrep -f "./$APP_NAME" > /dev/null; then
      echo "❗ Failed to stop $APP_NAME, trying kill -9..."
      kill -9 "$PID"
    else
      echo "✅ $APP_NAME stopped successfully."
    fi
  else
    echo "ℹ️ No running $APP_NAME process found."
  fi
}

# Function to start the app
start_app() {
  echo "🚀 Starting $APP_NAME..."
  nohup ./"$APP_NAME" > "$LOG_FILE" 2>&1 &
  echo "✅ $APP_NAME started with PID $(get_pid)"
}

# Function to restart the app
restart_app() {
  stop_app
  sleep 2
  start_app
}

# Function to build the app
build_app() {
  echo "🔨 Building $APP_NAME..."
  go build -o "$APP_NAME"
  if [ $? -eq 0 ]; then
    echo "✅ Build successful."
  else
    echo "❌ Build failed."
    exit 1
  fi
}

# Handle arguments
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
      echo "🟢 $APP_NAME is running (PID: $PID)"
    else
      echo "🔴 $APP_NAME is not running."
    fi
    ;;
  build)
    build_app
    ;;
  build-start)
    build_app
    start_app
    ;;
  build-restart)
    build_app
    restart_app
    ;;
  *)
    echo "Usage: $0 {start|stop|restart|status|build|build-start|build-restart}"
    exit 1
    ;;
esac
