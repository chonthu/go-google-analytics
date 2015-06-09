package sql

// Token represents a lexical token.
type Token int

const (
	// Special tokens
	ILLEGAL Token = iota
	EOF
	WS

	// Literals
	IDENT // main

	// Misc characters
	ASTERISK // *
	COMMA    // ,
	LTHAN    // <
	GTHAN    // >
	EQUAL    // =
	NOTEQUAL // <>, !=

	// Keywords
	SELECT  = 1
	FROM    = 2
	WHERE   = 3
	AND     = 3
	ORDERBY = 4
	ASC     = 5
	DESC    = 6
)
