export APP_BASE_DIR = $(PWD)

vite:
	pnpm vite
air:
	air
dev: 
	make vite & make air
