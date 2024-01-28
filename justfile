
# Build
build:
    cd src/c8y-tedge && go build -ldflags "-s -w" -o ../../c8y-tedge main.go
