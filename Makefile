REPO=github.com/sean9999/good-graph
SEMVER := $$(git tag --sort=-version:refname | head -n 1)

build: clean
	npm run build:all

publish:
	GOPROXY=https://goproxy.io,direct go list -m ${REPO}@${SEMVER}

install:
	npm install

run: build
	go run .

clean:
	rm -rf .parcel-cache
	rm -rf ./dist/*
