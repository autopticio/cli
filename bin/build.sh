# Linux x86
export GOOS=linux && export GOARCH=amd64 && export CGO_ENABLED=0 && go build -o ./exe/linux-amd64/autopticli ../src

# Mac ARM
export GOOS=darwin && export GOARCH=arm64 && export CGO_ENABLED=0 && go build -o ./exe/darwin-arm64/autopticli ../src

# Mac Intel
export GOOS=darwin && export GOARCH=amd64 && export CGO_ENABLED=0 && go build -o ./exe/darwin-amd64/autopticli ../src
