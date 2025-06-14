package examples

import (
	"fmt"
	"io"
	"net/http"

	a "github.com/randomdude16671/mini-build/arena"
)

// This DAG is still rough around the corners. Will improve it a lot.. this is a current example:

// This function is an example, this panics on an error, NEVER do that on an actual use of the DAG. Always handle errors gracefully.
func DoWork() {
	// function for inital future
	ffunc := func(args []any) string {
		resp, err := http.Get("https://randomdude16671homepage.netlify.app/")
		if err != nil {
			panic(err)
		}
		defer func() {
			if respError := resp.Body.Close(); respError != nil {
				fmt.Printf("Error closing response: %v\n", respError)
			}
		}()
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			panic(err)
		}
		return string(body)
	}
	// Second future
	f2func := func(args []any) string {
		fmt.Printf("Previous future result: %s\n", args[0])
		resp, err := http.Get("https://cowsay.morecode.org/say?message=Hello%3F+%0D%0A&format=text")
		if err != nil {
			panic(err)
		}
		defer func() {
			if respError := resp.Body.Close(); respError != nil {
				fmt.Printf("error: %v\n", respError)
			}
		}()
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			panic(err)
		}
		return string(body)
	}
	f := a.NewFuture[string](ffunc, nil)
	f2 := a.NewFuture[string](f2func, f)

	ar := a.NewArena(f)
	ar.AddFuture(f2)
	things := ar.AwaitAll()
	for i, str := range things {
		fmt.Printf("Result from future #%d:\n %s\n", i, str)
	}
}

// func main() {
// 	DoWork()
// }
