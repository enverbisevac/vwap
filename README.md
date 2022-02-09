# VWAP calculator
A realtime VWAP calculator of crypto currencies. It uses coinbase as its default provider for real time data over websocket.

## Run it
First, make sure that you have go version 1.17 installed on your machine. Then ...
```
make run
```
or 
```
make build
``` 
then ...
```
./vwap
```

## Tests
* Runs all the tests.
```
make test
``` 
* Runs the unit tests.
```
make test-unit
```
* Runs the integration test.
```
make test-intergration
``` 
