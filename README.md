# Stage

`Stage` is a command-line tool to organise development scripts and run sequences of commands on file change.

`Stage` can be installed as a Golang package. 

```bash
go get github.com/AIRTucha/stage
```

It executes commands from the `stage.yaml` file located on a current directory.

The file can contain the `_config` property with an array of globes to `watch`. Debounce time to avoid circular rebuilds and command to run steps. 

Each action can be defined as a property with an array of steps. Steps are executed sequentially.

```yaml
_config: 
  watch: 
    - '*.go'
    - 'src/*'
  debounce: 1000
  engine: '/bin/bash'
test:
  - echo "ok"
  - echo "ok2"
  - ping -c 4 google.com
  - apt-get
  - echo "ok3"
run:
  - go run *.go test -w
build-run:
  - go build -o bin/stage
  - ./bin/stage -watch test
install:
  - go install .
export:
  - PATH=$PATH:/Users/alextukalo/go/bin/
```

## Known limitations

`Stage` does not support `Windows`. Implementation of `Windows` support is very welcome.
