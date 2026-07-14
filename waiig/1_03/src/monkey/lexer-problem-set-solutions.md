# Lexer Problem Set — Solutions

Full code answers for every question in `lexer-problem-set.md`. Complete final `lexer.go` and `token.go` are at the bottom.

---

## Part 1 — Understanding what you already have

**Q1.1. What do `position` and `readPosition` each point to?**
**A.** `position` = index of the character currently in `l.ch`. `readPosition` = index of the next character to read.

**Q1.2. Why two pointers instead of one?**
**A.** Lexing needs to look ahead (e.g. checking if `=` is followed by `=`) without losing the current position. One index can't hold both "here" and "next" at once.

**Q1.3. What does `readChar()` set `l.ch` to at end of input, and why?**
**A.**
```go
if l.readPosition >= len(l.input) {
    l.ch = 0
}
```
`0` (NUL) is used because it can never appear as real source code content, so it's an unambiguous "no more input" signal.

**Q1.4. Why not use `' '` as the end-of-input sentinel?**
**A.** Because `' '` is already meaningful (whitespace to skip) — reusing it for EOF would make real trailing whitespace indistinguishable from "no more input."

**Q1.5. Why do `readIdentifier()`/`readNumber()` return early instead of falling through to `l.readChar()`?**
**A.** They already advance `l.ch` themselves in an internal loop:
```go
func (l *Lexer) readIdentifier() string {
	position := l.position
	for isLetter(l.ch) {
		l.readChar()
	}
	return l.input[position:l.position]
}
```
By the time they return, `l.ch` already holds the character *after* the identifier/number. An extra `readChar()` would skip it.

**Q1.6. What breaks if that early return is removed?**
**A.** Lexing `five;` would drop the `;` token — the character right after the identifier gets consumed by the stray `readChar()` call.

**Q1.7. Why does `LookupIdent` live in `token`, not `lexer`?**
**A.** Keyword-vs-identifier is a fact about the language's vocabulary (token types), which other consumers (parser, tooling) may need independent of lexing mechanics.

**Q1.8. What does allowing `_` in `isLetter` enable?**
**A.**
```go
func isLetter(ch byte) bool {
	return 'a' <= ch && ch <= 'z' || 'A' <= ch && ch <= 'Z' || ch == '_'
}
```
Lets `foo_bar` lex as one identifier instead of splitting at `_`.

**Q1.9. If `!`/`?` were added to `isLetter`, what would need re-checking?**
**A.** The `BANG` case and any `!=` lookahead logic — `!` would get swallowed into `readIdentifier()` and never reach `case '!':` at all. You'd need to lex `!` *before* falling into the identifier path, or exclude `!`/`?` from leading position.

---

## Part 2 — New single-character tokens

**Q2.1. New `token.go` constants:**
```go
MINUS    = "-"
SLASH    = "/"
ASTERISK = "*"
LT       = "<"
GT       = ">"
BANG     = "!"
```

**Q2.2–2.3. New `case` branches in `NextToken()`:**
```go
case '-':
    tok = newToken(token.MINUS, l.ch)
case '/':
    tok = newToken(token.SLASH, l.ch)
case '*':
    tok = newToken(token.ASTERISK, l.ch)
case '<':
    tok = newToken(token.LT, l.ch)
case '>':
    tok = newToken(token.GT, l.ch)
case '!':
    tok = newToken(token.BANG, l.ch)
```
`newToken` needs no changes — all six are single bytes.

**Q2.4. Test table additions:**
```go
input := `=+-/*<>!(){},;`

tests := []struct {
	expectedType    token.TokenType
	expectedLiteral string
}{
	{token.ASSIGN, "="},
	{token.PLUS, "+"},
	{token.MINUS, "-"},
	{token.SLASH, "/"},
	{token.ASTERISK, "*"},
	{token.LT, "<"},
	{token.GT, ">"},
	{token.BANG, "!"},
	{token.LPAREN, "("},
	{token.RPAREN, ")"},
	{token.LBRACE, "{"},
	{token.RBRACE, "}"},
	{token.COMMA, ","},
	{token.SEMICOLON, ";"},
	{token.EOF, ""},
}
```

**Q2.5. Any other functions touched?**
**A.** No — `skipWhitespace`, `readIdentifier`, `readNumber` are untouched.

---

## Part 3 — Two-character tokens: `==` and `!=`

**Q3.1–3.2. `peekChar()`:**
```go
func (l *Lexer) peekChar() byte {
	if l.readPosition >= len(l.input) {
		return 0
	}
	return l.input[l.readPosition]
}
```
Read-only: never touches `l.position`, `l.readPosition`, or `l.ch`.

**Q3.4. New constants:**
```go
EQ     = "=="
NOT_EQ = "!="
```

**Q3.5–3.7. Updated `case '='` and `case '!'`:**
```go
case '=':
    if l.peekChar() == '=' {
        ch := l.ch
        l.readChar()
        literal := string(ch) + string(l.ch)
        tok = token.Token{Type: token.EQ, Literal: literal}
    } else {
        tok = newToken(token.ASSIGN, l.ch)
    }
case '!':
    if l.peekChar() == '=' {
        ch := l.ch
        l.readChar()
        literal := string(ch) + string(l.ch)
        tok = token.Token{Type: token.NOT_EQ, Literal: literal}
    } else {
        tok = newToken(token.BANG, l.ch)
    }
```
`newToken` isn't reused for the two-char branch because it only wraps a single `byte` — the two-character token is built as a `token.Token{}` literal directly instead.

**Q3.8. Why `peekChar()` and not `readChar()` to check ahead?**
**A.** `readChar()` commits — it advances the lexer. If you used it to "peek" and the next char wasn't `=`, you'd have already consumed and lost the character you still needed to lex on the next call.

**Q3.9. Test additions:**
```go
input := `5 == 5;
5 != 10;
x = 5;
!true`

tests := []struct {
	expectedType    token.TokenType
	expectedLiteral string
}{
	{token.INT, "5"},
	{token.EQ, "=="},
	{token.INT, "5"},
	{token.SEMICOLON, ";"},
	{token.INT, "5"},
	{token.NOT_EQ, "!="},
	{token.INT, "10"},
	{token.SEMICOLON, ";"},
	{token.IDENT, "x"},
	{token.ASSIGN, "="},
	{token.INT, "5"},
	{token.SEMICOLON, ";"},
	{token.BANG, "!"},
	{token.TRUE, "true"},
	{token.EOF, ""},
}
```

---

## Part 4 — New keywords

**Q4.1. New constants:**
```go
TRUE   = "TRUE"
FALSE  = "FALSE"
IF     = "IF"
ELSE   = "ELSE"
RETURN = "RETURN"
```

**Q4.2. Updated `keywords` map:**
```go
var keywords = map[string]TokenType{
	"fn":     FUNCTION,
	"let":    LET,
	"true":   TRUE,
	"false":  FALSE,
	"if":     IF,
	"else":   ELSE,
	"return": RETURN,
}
```

**Q4.3–4.4. Changes needed in `lexer.go`:**
**A.** None. Keywords flow entirely through `readIdentifier()` → `token.LookupIdent()`. Needing a lexer-side edit would mean keyword handling had leaked out of the table-driven design.

**Q4.5. Test additions:**
```go
input := `if (5 < 10) {
	return true;
} else {
	return false;
}`

tests := []struct {
	expectedType    token.TokenType
	expectedLiteral string
}{
	{token.IF, "if"},
	{token.LPAREN, "("},
	{token.INT, "5"},
	{token.LT, "<"},
	{token.INT, "10"},
	{token.RPAREN, ")"},
	{token.LBRACE, "{"},
	{token.RETURN, "return"},
	{token.TRUE, "true"},
	{token.SEMICOLON, ";"},
	{token.RBRACE, "}"},
	{token.ELSE, "else"},
	{token.LBRACE, "{"},
	{token.RETURN, "return"},
	{token.FALSE, "false"},
	{token.SEMICOLON, ";"},
	{token.RBRACE, "}"},
	{token.EOF, ""},
}
```

---

## Part 5 — Debugging challenge

**Q5.1. Missing line:**
```go
func (l *Lexer) skipWhitespace() {
	for l.ch == ' ' || l.ch == '\t' || l.ch == '\n' || l.ch == '\r' {
		l.readChar()   // <-- this was missing
	}
}
```

**Q5.2. Why only whitespace inputs hang:**
**A.** If `l.ch` never matches one of the four whitespace bytes, the loop condition is false on the first check and the body never runs — no hang. Only once `l.ch` genuinely equals whitespace does the (empty) body execute, and without `readChar()` advancing it, the condition stays true forever.

**Q5.3. Effect of calling `skipWhitespace()` after the `switch` instead of before:**
**A.** The `switch` would evaluate `l.ch` while it might still be a space/tab/newline. None of the `case` branches match whitespace, so it falls into `default`, where `isLetter`/`isDigit` are both false — producing a spurious `ILLEGAL` token for every whitespace character instead of silently skipping it.

---

## Part 6 — Stretch goals

**Q6.1. Current behavior of `-5`:**
**A.** Two tokens: `{MINUS, "-"}` then `{INT, "5"}`.

**Q6.2. Argument for keeping it that way:**
**A.** The lexer has no surrounding context to know whether `-` is unary negation or subtraction (`a - 5` looks identical token-by-token up to this point). Leaving disambiguation to the parser, which sees the full expression, is simpler and more correct than guessing in the lexer.

**Q6.3. Sketch of float support in `readNumber()`:**
```go
func (l *Lexer) readNumber() (string, token.TokenType) {
	position := l.position
	tokType := token.INT
	for isDigit(l.ch) {
		l.readChar()
	}
	if l.ch == '.' && isDigit(l.peekChar()) {
		tokType = token.FLOAT
		l.readChar() // consume '.'
		for isDigit(l.ch) {
			l.readChar()
		}
	}
	return l.input[position:l.position], tokType
}
```
A second `.` wouldn't satisfy `isDigit(l.peekChar())` after the first float is fully consumed, so `3.14.15` would lex as `FLOAT("3.14")` followed by `ILLEGAL(".")`  then `INT("15")` — worth flagging explicitly as a known edge case rather than silently accepting it.

**Q6.4. New token type needed:**
```go
FLOAT = "FLOAT"
```
Yes — `INT` shouldn't double as the float type, since later stages (evaluator) need to know which Go numeric type to parse the literal into.

**Q6.5. Unterminated string handling:**
```go
func (l *Lexer) readString() string {
	position := l.position + 1
	for {
		l.readChar()
		if l.ch == '"' || l.ch == 0 {
			break
		}
	}
	return l.input[position:l.position]
}
```
Breaking on `l.ch == 0` (EOF) as well as `'"'` prevents an infinite loop; the caller can check whether the loop ended because of EOF and emit `ILLEGAL` instead of `STRING` in that case.

**Q6.6. Change needed for Unicode identifiers:**
**A.** `l.ch` would need to become `rune` instead of `byte`, and reading would need to decode UTF-8 (e.g. via `utf8.DecodeRuneInString`) rather than single-byte indexing, since multi-byte characters can't be read correctly one byte at a time.

---

## Complete final `token.go`

```go
package token

type TokenType string

type Token struct {
	Type    TokenType
	Literal string
}

const (
	// Special types
	ILLEGAL = "ILLEGAL"
	EOF     = "EOF"

	// Identifiers + literals
	IDENT = "IDENT" // add, foobar, x, y, ...
	INT   = "INT"   // 1343456

	// Operators
	ASSIGN   = "="
	PLUS     = "+"
	MINUS    = "-"
	SLASH    = "/"
	ASTERISK = "*"
	BANG     = "!"
	LT       = "<"
	GT       = ">"
	EQ       = "=="
	NOT_EQ   = "!="

	// Delimiters
	COMMA     = ","
	SEMICOLON = ";"
	LPAREN    = "("
	RPAREN    = ")"
	LBRACE    = "{"
	RBRACE    = "}"

	// Keywords
	FUNCTION = "FUNCTION"
	LET      = "LET"
	TRUE     = "TRUE"
	FALSE    = "FALSE"
	IF       = "IF"
	ELSE     = "ELSE"
	RETURN   = "RETURN"
)

var keywords = map[string]TokenType{
	"fn":     FUNCTION,
	"let":    LET,
	"true":   TRUE,
	"false":  FALSE,
	"if":     IF,
	"else":   ELSE,
	"return": RETURN,
}

func LookupIdent(ident string) TokenType {
	if tok, ok := keywords[ident]; ok {
		return tok
	}
	return IDENT
}
```

## Complete final `lexer.go`

```go
package lexer

import "monkey/token"

type Lexer struct {
	input        string
	position     int  // current position in input (points to current char)
	readPosition int  // current reading position in input (after current char)
	ch           byte // current char under examination
}

func New(input string) *Lexer {
	l := &Lexer{input: input}
	l.readChar()
	return l
}

func (l *Lexer) NextToken() token.Token {
	var tok token.Token

	l.skipWhitespace()

	switch l.ch {
	case '=':
		if l.peekChar() == '=' {
			ch := l.ch
			l.readChar()
			literal := string(ch) + string(l.ch)
			tok = token.Token{Type: token.EQ, Literal: literal}
		} else {
			tok = newToken(token.ASSIGN, l.ch)
		}
	case '+':
		tok = newToken(token.PLUS, l.ch)
	case '-':
		tok = newToken(token.MINUS, l.ch)
	case '!':
		if l.peekChar() == '=' {
			ch := l.ch
			l.readChar()
			literal := string(ch) + string(l.ch)
			tok = token.Token{Type: token.NOT_EQ, Literal: literal}
		} else {
			tok = newToken(token.BANG, l.ch)
		}
	case '/':
		tok = newToken(token.SLASH, l.ch)
	case '*':
		tok = newToken(token.ASTERISK, l.ch)
	case '<':
		tok = newToken(token.LT, l.ch)
	case '>':
		tok = newToken(token.GT, l.ch)
	case ';':
		tok = newToken(token.SEMICOLON, l.ch)
	case '(':
		tok = newToken(token.LPAREN, l.ch)
	case ')':
		tok = newToken(token.RPAREN, l.ch)
	case ',':
		tok = newToken(token.COMMA, l.ch)
	case '{':
		tok = newToken(token.LBRACE, l.ch)
	case '}':
		tok = newToken(token.RBRACE, l.ch)
	case 0:
		tok.Literal = ""
		tok.Type = token.EOF
	default:
		if isLetter(l.ch) {
			tok.Literal = l.readIdentifier()
			tok.Type = token.LookupIdent(tok.Literal)
			return tok
		} else if isDigit(l.ch) {
			tok.Type = token.INT
			tok.Literal = l.readNumber()
			return tok
		} else {
			tok = newToken(token.ILLEGAL, l.ch)
		}
	}

	l.readChar()
	return tok
}

func (l *Lexer) readChar() {
	if l.readPosition >= len(l.input) {
		l.ch = 0
	} else {
		l.ch = l.input[l.readPosition]
	}
	l.position = l.readPosition
	l.readPosition += 1
}

func (l *Lexer) peekChar() byte {
	if l.readPosition >= len(l.input) {
		return 0
	}
	return l.input[l.readPosition]
}

func (l *Lexer) skipWhitespace() {
	for l.ch == ' ' || l.ch == '\t' || l.ch == '\n' || l.ch == '\r' {
		l.readChar()
	}
}

func (l *Lexer) readIdentifier() string {
	position := l.position
	for isLetter(l.ch) {
		l.readChar()
	}
	return l.input[position:l.position]
}

func (l *Lexer) readNumber() string {
	position := l.position
	for isDigit(l.ch) {
		l.readChar()
	}
	return l.input[position:l.position]
}

func isLetter(ch byte) bool {
	return 'a' <= ch && ch <= 'z' || 'A' <= ch && ch <= 'Z' || ch == '_'
}

func isDigit(ch byte) bool {
	return '0' <= ch && ch <= '9'
}

func newToken(tokenType token.TokenType, ch byte) token.Token {
	return token.Token{Type: tokenType, Literal: string(ch)}
}
```

*Note: `readNumber`/`FLOAT` and `readString`/`STRING` from Part 6 are sketches only and are not included in the final listing above — they'd change `readNumber`'s signature and add a new `case '"'` branch if you choose to implement them.*
