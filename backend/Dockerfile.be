FROM golang:1.23

WORKDIR /app

# 複製 Go modules 文件並安裝依賴
COPY go.mod go.sum ./
RUN go mod download

# 複製設定檔
COPY settings.yml ./settings.yml

# 複製程式碼並編譯
COPY . .
RUN go build -o backend main.go

# 曝露服務埠
EXPOSE 8080

# 啟動應用
CMD ["./backend"]