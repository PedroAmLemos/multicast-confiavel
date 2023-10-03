BINARY_NAME=multicast-confiavel

build:
	@echo "Compilando o código..."
	go build -o $(BINARY_NAME)

run: build
	@echo "Executando o código..."
	./$(BINARY_NAME)

clean:
	@echo "Limpando arquivos binários..."
	go clean

all: build

