HOST := "0.0.0.0"
PORT := "9000"
USERNAME := $(shell whoami)
SECRET_STORE_PATH := "./temp"
WS_ENDPOINT := "/ws"

build:
	@echo "Building Websocket Server image..."
	- podman build -t codebox-wss:latest .

start: build
	@echo "Starting Websocket Server container... $(USERNAME)"
	- podman run --name codebox-wss --replace -d \
		-p $(PORT):$(PORT) \
		-e PORT=$(PORT) \
		-e HOST=$(HOST) \
		-e USERNAME=$(USERNAME) \
		-e SECRET_STORE_PATH=$(SECRET_STORE_PATH) \
		-e WS_ENDPOINT=$(WS_ENDPOINT) \
		-v temp:/app/temp:Z \
		--net codebox-network \
		localhost/codebox-wss:latest


