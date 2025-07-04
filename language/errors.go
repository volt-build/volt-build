package language

import "fmt"

type CustomError struct {
	// Name of the file where the error was caused.
	Filename string
	// Message for the error
	Message string
	// Line of the error (For formatting properly)
	Ln int
	// Column of the error (Place ^ or ^~~ at here)
	Col int
}

func NewCustomError(file, msg string, ln, col int) *CustomError {
	customError := &CustomError{
		Filename: file,
		Message:  msg,
		Ln:       ln,
		Col:      col,
	}
	return customError
}

func (c *CustomError) Human() {
}

func (c *CustomError) Error() string {
	// Efficient error for machine readability
	return fmt.Sprintf(
		"%s:%d:%d: %s\n",
		c.Filename, c.Ln, c.Col,
		c.Message,
	)
}
