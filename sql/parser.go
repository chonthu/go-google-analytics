package sql

import (
	"fmt"
	"io"
	"log"
)

type WhereCondition struct {
	Key    string
	Filter string
	Value  string
}

// SelectStatement represents a SQL SELECT statement.
type SelectStatement struct {
	Fields      []string
	TableName   string
	WhereClause []WhereCondition
	OrderBy     string
	Sort        string
}

// Parser represents a parser.
type Parser struct {
	s   *Scanner
	buf struct {
		tok Token  // last read token
		lit string // last read literal
		n   int    // buffer size (max=1)
	}
}

// NewParser returns a new instance of Parser.
func NewParser(r io.Reader) *Parser {
	return &Parser{s: NewScanner(r)}
}

// Parse parses a SQL SELECT statement.
func (p *Parser) Parse() (*SelectStatement, error) {
	stmt := &SelectStatement{}

	// First token should be a "SELECT" keyword.
	if tok, lit := p.scanIgnoreWhitespace(); tok != SELECT {
		return nil, fmt.Errorf("found %q, expected SELECT", lit)
	}

	// Next we should loop over all our comma-delimited fields.
	for {
		// Read a field.
		tok, lit := p.scanIgnoreWhitespace()
		if tok != IDENT && tok != ASTERISK {
			return nil, fmt.Errorf("found %q, expected field", lit)
		}
		stmt.Fields = append(stmt.Fields, lit)

		// If the next token is not a comma then break the loop.
		if tok, _ := p.scanIgnoreWhitespace(); tok != COMMA {
			p.unscan()
			break
		}
	}

	// Next we should see the "FROM" keyword.
	if tok, lit := p.scanIgnoreWhitespace(); tok != FROM {
		return nil, fmt.Errorf("found %q, expected FROM", lit)
	}

	// Finally we should read the table name.
	tok, lit := p.scanIgnoreWhitespace()
	if tok != IDENT {
		return nil, fmt.Errorf("found %q, expected table name", lit)
	}
	stmt.TableName = lit

	for {
		if tok, _ := p.scanIgnoreWhitespace(); tok != WHERE && tok != AND {
			fmt.Println("h:", lit)
			break
		}

		for {

			err := scanWhere(p, stmt)
			if err != nil {
				return nil, err
			}

			// If the next token is not a comma then break the loop.
			if tok, _ := p.scanIgnoreWhitespace(); tok != COMMA {
				p.unscan()
				break
			}
		}
	}

	// check for by
	tok, lit = p.scanIgnoreWhitespace()
	if tok != IDENT {
		return nil, fmt.Errorf("found %q, expected BY", lit)
	}

	// check for by
	tok, lit = p.scanIgnoreWhitespace()
	if tok != IDENT {
		return nil, fmt.Errorf("found %q, expected ORDER BY value", lit)
	}

	stmt.OrderBy = lit

	// check for by
	tok, lit = p.scanIgnoreWhitespace()
	if tok != IDENT {
		stmt.Sort = "desc"
	} else {
		stmt.Sort = lit
	}

	return stmt, nil
}

func scanWhere(p *Parser, stmt *SelectStatement) error {

	where := WhereCondition{}

	// Read a field.
	tok, lit := p.scanIgnoreWhitespace()

	if tok != IDENT && tok != ASTERISK {
		return fmt.Errorf("found %q, expected where key", lit)
	}

	where.Key = lit

	// Read a field.
	tok, lit = p.scanIgnoreWhitespace()

	if tok != EQUAL && tok != LTHAN && tok != GTHAN && tok != ASTERISK {
		return fmt.Errorf("found %q, expected where filter", lit)
	}

	where.Filter = lit

	// Read a field.
	tok, lit = p.scanIgnoreWhitespace()

	if tok != IDENT && tok != ASTERISK {
		return fmt.Errorf("found %q, expected where value", lit)
	}

	where.Value = lit

	stmt.WhereClause = append(stmt.WhereClause, where)

	return nil
}

// scan returns the next token from the underlying scanner.
// If a token has been unscanned then read that instead.
func (p *Parser) scan() (tok Token, lit string) {
	// If we have a token on the buffer, then return it.
	if p.buf.n != 0 {
		p.buf.n = 0
		return p.buf.tok, p.buf.lit
	}

	// Otherwise read the next token from the scanner.
	tok, lit = p.s.Scan()
	log.Println(lit)

	// Save it to the buffer in case we unscan later.
	p.buf.tok, p.buf.lit = tok, lit

	return
}

// scanIgnoreWhitespace scans the next non-whitespace token.
func (p *Parser) scanIgnoreWhitespace() (tok Token, lit string) {
	tok, lit = p.scan()
	if tok == WS {
		tok, lit = p.scan()
	}
	return
}

// unscan pushes the previously read token back onto the buffer.
func (p *Parser) unscan() { p.buf.n = 1 }
