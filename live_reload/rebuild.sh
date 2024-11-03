#!/usr/bin/env bash

# watch code changes, trigger re-build, and kill process
while true; do
    go build -o _build/launcher && pkill -f '_build/launcher'
    # Use -1 to make fswatch exit after a single change is detected
    fswatch -1 --event Updated --event Created --event Renamed --event Removed $(find . -name '*.go')
done
