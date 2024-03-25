.PHONY: help
## help: print this help message
help:
	@echo "Usage: \n"
	@sed -n 's/^##//p' ${MAKEFILE_LIST} | column -t -s ':' |  sed -e 's/^/ /'

.PHONY: hash
## hash: update static files hashes in index.html
hash:
	sed -i '' -E "s/styles\.css\?crc=[0-9a-z]+/styles.css?crc=$$(crc32 ./pages/styles.css)/" pages/index.html
	sed -i '' -E "s/index\.js\?crc=[0-9a-z]+/index.js?crc=$$(crc32 ./pages/index.js)/" pages/index.html

.PHONY: dev
## dev: run the wrangler pages dev
dev: hash
	wrangler pages dev pages/
