# Gonitor
Poll websites over HTTP, analyze return codes and emit alerts (e.g. e-mails).

It's a toy project to try out Go :)

## Todo next :
- [ ] Use config to spawn go-routines (use config.json in current path)

## Ideas
- [ ] Add parameter to load config from arbibtrary file
- [ ] Add parameter to load config from inline
- [ ] Emit alerts via e-mail
- [ ] Save return codes in database
- [ ] Save alerts in database
- [ ] Hot reload config
- [ ] HTML dashboard with config
- [ ] HTML dashboard showing datapoints with d3js / websockets
- [ ] Try out React + RxJS

# Done !
- [x] Poll websites
- [x] Load config from JSON file