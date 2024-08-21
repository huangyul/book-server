.PHONY: all server web

all: server web

server:
	cd server && air

web:
	cd web && pnpm dev
