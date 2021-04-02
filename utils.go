package main

func Notify(
	boolChan chan bool,
	boolVal bool,
) {
	select {
	case boolChan <- boolVal:
	default:
	}
}
