
build: clean
	npm run build:all

publish:
	npm run publish

install:
	npm install

run:
	go run .

clean:
	rm -rf .parcel-cache
	rm -rf ./dist/*
