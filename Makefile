.PHONY: all server web wire

all: server web

server:
	cd server && air

web:
	cd web && pnpm dev

wire:
	cd server && wire
