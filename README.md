# Gonitor
[![Build Status](https://travis-ci.org/Kehrlann/gonitor.svg?branch=master)](https://travis-ci.org/Kehrlann/gonitor)

Poll websites over HTTP, analyze return codes and emit alerts (e.g. e-mails).

It's a toy project to try out Go :)

## Todo next
- [ ] Basic page with websockets, showing messages as they pop up

## Ideas
### Testing 
- [ ] Integration testing of main.go ?

### Features :
- [ ] Update the help/man page for config with Commands
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
- [ ] HTML dashboard with config
- [ ] HTML dashboard showing datapoints with d3js / websockets
- [ ] Try out React + RxJS ?

## Done !
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

## Abandonned ? 
- [ ] Add a "wizard" to configure a polling on first launch if no config.json found

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
