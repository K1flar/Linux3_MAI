all: clean
	go build -o ./bin/t1 ./task1/main.go

clean:
	rm -rf ./bin/*
