.ONESHELL:
dev:
	cd ./server
	gow run . &
	cd ..
	cd app
	yarn dev
