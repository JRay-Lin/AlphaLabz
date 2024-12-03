#!/bin/bash

set -e

# 設置變數
PB_BINARY="./pocketbase"
PB_VERSION="0.23.4"
PB_URL="https://github.com/pocketbase/pocketbase/releases/download/v${PB_VERSION}/pocketbase_${PB_VERSION}_linux_amd64.zip"

echo "Starting PocketBase initialization..."

# 檢查 PocketBase 是否存在，若不存在則下載
if [ ! -f "${PB_BINARY}" ]; then
    echo "PocketBase binary not found. Downloading..."
    wget -q "${PB_URL}" -O pocketbase.zip
    unzip pocketbase.zip
    rm pocketbase.zip
    chmod +x "${PB_BINARY}"
    echo "PocketBase downloaded and prepared."
else
    echo "PocketBase binary already exists. Skipping download."
fi

# 啟動 PocketBase 服務，並後台執行
"${PB_BINARY}" serve --http=0.0.0.0:8090 &
PB_PID=$!

# 等待 PocketBase 完全啟動
echo "Waiting for PocketBase to start..."
sleep 5

# 使用 `superuser upsert` 指令來建立或更新 superuser 帳戶
echo "Creating or updating superuser..."
"${PB_BINARY}" superuser upsert "${ADMIN_EMAIL}" "${ADMIN_PASSWORD}"

# 檢查是否成功
if [ $? -eq 0 ]; then
    echo "Superuser created or updated successfully."
else
    echo "Failed to create or update superuser."
    kill $PB_PID
    exit 1
fi

# 等待 PocketBase 停止，防止腳本提前退出
echo "PocketBase is running. Waiting for process to terminate..."
wait $PB_PID