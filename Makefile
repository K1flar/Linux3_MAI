all: clean
	cd task1 && go build -o ../bin/t1 ./main.go

clean:
	rm -rf ./bin/*
