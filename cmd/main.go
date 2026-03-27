package cmd

import (
	"fmt"
	"net/http"
)

func main() {
	fmt.Println("Server starting")
	http.ListenAndServe(":8080", nil)
}
