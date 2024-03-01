CC=cargo
FMT=fmt

OPTIONS=

default:
	$(CC) build
	@make clippy

fmt:
	$(CC) fmt --all

check:
	$(CC) test --all

example:
	@echo "No example for the moment"

clean:
	$(CC) clean

clippy:
	$(CC) clippy --all --tests

coffee:
	$(CC) build --release
