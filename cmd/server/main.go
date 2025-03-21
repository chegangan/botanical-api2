package main

func main() {
	server := NewServer()
	if err := server.Run(); err != nil {
		panic(err)
	}
}
