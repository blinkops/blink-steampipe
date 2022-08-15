install:
	CGO_ENABLED=0 GOOS=linux go build -o ./build/scripts/bin/generate  ./build/scripts/*.go


