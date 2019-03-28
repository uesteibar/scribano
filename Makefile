.SUFFIXES:

test:
	source .env.test && \
	go test ./... -cover &&\
	source .env.dev

