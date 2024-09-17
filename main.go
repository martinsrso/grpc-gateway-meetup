package main

func main() {
	go RunGRPCServer()
	RunGRPCGWServer()
}
