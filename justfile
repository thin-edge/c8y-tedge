
# Build
build:
    mkdir -p bin
    cd src/c8y-tedge && go build -ldflags "-s -w" -o ../../bin/c8y-tedge main.go

run *args:
    cd src/c8y-tedge && go run main.go {{args}}
