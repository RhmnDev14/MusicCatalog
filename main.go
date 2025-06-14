package main

import "music_catalog/internal"

func main() {
	server := internal.NewServer()
	if server == nil {
		println("gagal running")
	}

	println("Server running on", server)
}
