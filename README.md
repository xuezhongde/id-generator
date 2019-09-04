## Install
- Install Go (1.11+) and set your GOPATH
- go get github.com/xuezhongde/id-generator
- cd $GOPATH/src/github.com/xuezhongde/id-generator
- make

## Startup
```
$ ./bin/id-gen -c=./etc/id.toml 
```

## Usage
```
Usage: id-gen [-hv] [-c config file] [-d data center id] [-w worker id] [-p port] [-r router]
Options:
  -c string
    	id generator config file (default "./etc/id.toml")
  -d int
    	data center id (default -1)
  -h	this help
  -p int
    	listen on port (default 8000)
  -r string
    	router path (default "/id")
  -v	show version and exit
  -w int
    	worker id (default -1)
```