# Gonitor
Poll websites over HTTP, analyze return codes and emit alerts (e.g. e-mails).

It's a toy project to try out Go :)

## Todo next
- [ ] Generate an example config file on startup : make special error in config.go 
    - throw it when gonitor.config.json not found 
    - handle file creation with sane defaults 
    - print warning message, saying 'oh we created this file for you'
    - what about when file creation is disabled ? should just print the file structure ?

## Ideas
- [ ] Mock HTTP and SMTP for tests ?
- [ ] Test main.go -> emit, config file ...
- [ ] Configure log level
- [ ] Allow running external scripts
- [ ] Per-resource recipient-list
- [ ] Add a "wizard" to configure a polling on first launch if no config.json found
- [ ] Save return codes in database
- [ ] Save alerts in database
- [ ] Hot reload config
- [ ] HTML dashboard with config
- [ ] HTML dashboard showing datapoints with d3js / websockets
- [ ] Try out React + RxJS

## Done !
- [x] Add parameter to load config from arbibtrary file
- [x] Don't send an e-mail if SMTP is not configured properly
- [x] Emit alerts via e-mail
- [x] Poll websites
- [x] Load config from JSON file
- [x] Use config to spawn go-routines (use config.json in current path)
- [x] Add logging : messages, errors ...
- [x] Add a welcome message on startup

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
