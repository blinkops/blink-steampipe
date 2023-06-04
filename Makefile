install:
	GOOS=linux GOARCH=amd64 go build -o /bin/generate  ./scripts/*.go

local:
	docker pull blinkops/blink-steampipe-$(plugin):$(image)
	docker run -it --rm --entrypoint bash blinkops/blink-steampipe-$(plugin):$(image)

