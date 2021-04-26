package main

func main() {
	// fatal error: all goroutines are asleep - deadlock!
	select {}
}
