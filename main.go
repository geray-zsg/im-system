package main

func main() {

	// NewServer("127.0.0.1", 8888).Start()
	server := NewServer("127.0.0.1", 8888)
	server.Start()
}
