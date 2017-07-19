# Gonitor
[![Build Status](https://travis-ci.org/Kehrlann/gonitor.svg?branch=master)](https://travis-ci.org/Kehrlann/gonitor)

Poll websites over HTTP, analyze return codes and emit alerts (e.g. e-mails).

It's a toy project to try out Go :)

## Todo next

## Ideas
### Testing 
- [ ] Integration testing of main.go ?

### Features :
- [ ] Add a "help" flag
- [ ] Add a "example" flag
- [ ] Configure log level
- [ ] Daily summary e-mail (saying i'm alive !) 
- [ ] Per-resource recipient-list
- [ ] Hot reload config
- [ ] Print datetime in e-mail

### Features : visualisation tool
- [ ] Save alerts in database
- [ ] Save return codes in database
- [ ] HTML dashboard showing datapoints with d3js / websockets
- [ ] Try out React + RxJS ?
- [ ] HTML dashboard with config

### Misc ideas
- [ ] Try gotrace (https://github.com/divan/gotrace)

## Done !
- [x] Update the default config page for config with Commands
- [x] Add a "hooks" folder
- [x] Add basic HTTP server to display an HTML welcome page
- [x] Add Travis CI
- [x] Refactor main to be more of glue code (e.g. Resource slice -> worker)
- [x] Package-ify all the things (make things private where they need to be)
- [x] Mock SMTP for tests
- [x] HttpMock -> Done
- [x] Test exec command (works on Linux)
- [x] Allow running external scripts
- [x] Add parameter to load config from arbibtrary file
- [x] Don't send an e-mail if SMTP is not configured properly
- [x] Emit alerts via e-mail
- [x] Poll websites
- [x] Load config from JSON file
- [x] Use config to spawn go-routines (use config.json in current path)
- [x] Add logging : messages, errors ...
- [x] Add a welcome message on startup
- [x] Generate an example config file on startup : 
    - [x] make special error in config.go 
    - [x] throw it when gonitor.config.json not found  
    - [x] print warning message, saying 'oh we created this file for you'
    - [x] what about when file creation is disabled ? should just print the file structure ?
- [x] Add a websockets server, shooting fake StateChangeMessages with a ticker
- [x] Display those messages in the default index page
- [x] Hook up an emitter to the websockets server
- [x] Write a wrapper for gorilla websocket connection
    - Make it thread-safe
- [x] Try expvarmon (https://github.com/divan/expvarmon)


## Abandonned ? 
- [ ] Add a "wizard" to configure a polling on first launch if no config.json found

## Expvarmon :
This is extremely cool : https://github.com/divan/expvarmon . See more at : http://developers-club.com/posts/257593/
First install expvarmon :
`go get github.com/divan/expvarmon`

For counting Goroutines, add to main.go : 
```
    import _ "expvar"
    
    [...]
    
    // in main() 
    Publish("Goroutines", expvar.Func(func() interface{} { return runtime.NumGoroutine() } ))
```

Then try : 
```
    ./expvarmon -ports="3000" -vars="mem:memstats.Alloc,Goroutines"
```


## Profiling :
28/03/2017 : there seems to be a very small memory leak ... Hard to pin down.

In main.go, do :
```go
    import _ "net/http/pprof"
    [...]
	go http.ListenAndServe(":8080", http.DefaultServeMux)
    [...]
```

And then, from the command line, in the folder :
```sh
    go tool pprof gonitor http://localhost:8080/debug/pprof/heap
```
