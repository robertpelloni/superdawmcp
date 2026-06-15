package main
import "fmt"
import "os"
func main() {
	fmt.Fprintf(os.Stdout, "width: 100%%%%\n")
	fmt.Fprintf(os.Stdout, "width: 100%%\n")
}
