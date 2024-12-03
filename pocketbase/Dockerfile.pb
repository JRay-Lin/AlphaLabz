FROM alpine:3.18

WORKDIR /pb

# 安裝所需工具
RUN apk add --no-cache bash curl

# 複製初始化腳本與相關檔案
COPY init-pocketbase.sh /pb/init-pocketbase.sh
COPY pb_hooks /pb_hooks

# 確保腳本可執行
RUN chmod +x /pb/init-pocketbase.sh

EXPOSE 8090

CMD ["sh", "/pb/init-pocketbase.sh"]