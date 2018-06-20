package circuitbreaker

// StatusEnum circuitbreaker status enum
type StatusEnum int

const (
	_ StatusEnum = iota
	// StatusClosed StatusClosed
	StatusClosed
	// StatusOpen StatusOpen
	StatusOpen
	// StatusHalfOpen StatusHalfOpen
	StatusHalfOpen
)
