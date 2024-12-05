FROM alpine:3.18

WORKDIR /pb

RUN apk add --no-cache bash curl

COPY init-pocketbase.sh /pb/init-pocketbase.sh
COPY pb_hooks /pb_hooks

RUN chmod +x /pb/init-pocketbase.sh

EXPOSE 8090

CMD ["sh", "/pb/init-pocketbase.sh"]