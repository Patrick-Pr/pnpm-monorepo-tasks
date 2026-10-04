 CGO_ENABLED=0 GOOS=darwin GOARCH=arm64 \
  go build -trimpath -ldflags="-s -w" \
  -o dist/golist ./cmd/monorepo-tasks

cp dist/golist ~/monorepo-tasks/golist
