.PHONY: build
build:
	go build -o /tmp/bin/app main.go

.PHONY: run
run: build
	/tmp/bin/app $(bin)

.PHONY: test
test:
	go test -race -coverprofile="coverage.out" ./...

.PHONY: generate
generate:
	go generate ./...

.PHONY: swag
swag:
	swag init -g cmd/api.go -o docs

.PHONY: run/live
run/live:
	go run github.com/cosmtrek/air@v1.43.0 \
		--build.cmd "make build" --build.bin "/tmp/bin/app $(bin)" --build.delay "100" \
		--build.exclude_dir "" \
		--build.include_ext "go, tpl, tmpl, html, css, scss, js, ts, sql, jpeg, jpg, gif, png, bmp, svg, webp, ico" \
		--misc.clean_on_exit "true"