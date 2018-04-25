
websocket project

using go version `1.7.5`

set up a `GOPATH` variable, which will be the home of all your Go projects. For example `~/workspace/go`. You can export it for the local session, export `GOPATH=...` or set it permanently in you `~/.bash_profile`. 

create the directories `$GOPATH/bin` `$GOPATH/src` `$GOPATH/pkg`.

add the gobin to your path, `export PATH=$PATH:/$GOPATH/bin`. again do it for the local session or permanentyly.

create the directory `$GOPATH/src/github.com/packetzoom/` and clone this project there.

build by going into the subproject `example`and running `go build`. this will create a binary name after the subproject

we use the tool `govendor` to manage libraries. dependencies will be in the vendor directory which is checked into git, as you add more dependencies make sure to run `govendor add +external` 

if you do not have govendor install it by running `go get  github.com/kardianos/govendor`

*Start the example game server*
```
cd example
go build && ./example
```


