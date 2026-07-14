# Lexer Problem Set — Build Steps (Code Answers + Explanations)

Each question asks you to write the next piece of the lexer, in the order the text actually builds it. Every answer is a code block, followed by a short note on why that piece is being added.

---

## Q1. Define the `Lexer` struct with its four fields, and a first version of `New()` that just wraps the input (no character reading yet).

**A.**
```go
package lexer

type Lexer struct {
	input        string
	position     int  // current position in input (points to current char)
	readPosition int  // current reading position in input (after current char)
	ch           byte // current char under examination
}

func New(input string) *Lexer {
	l := &Lexer{input: input}
	return l
}
```
**Why:** `position` and `readPosition` exist as two separate fields because the lexer will need to "peek" past the current character to see what comes next, and one index can't represent both "current" and "next" at once. `input` is kept as a plain `string` rather than an `io.Reader` to avoid the extra complexity of tracking filenames and line numbers — that's left for a production implementation.

---

## Q2. Write `readChar()` — the method that reads the next character and advances `position`/`readPosition`.

**A.**
```go
func (l *Lexer) readChar() {
	if l.readPosition >= len(l.input) {
		l.ch = 0
	} else {
		l.ch = l.input[l.readPosition]
	}
	l.position = l.readPosition
	l.readPosition += 1
}
```
**Why:** This is the core mechanism that gives the lexer its next character and advances its position in the input string. Setting `l.ch = 0` at the end of input gives the lexer an explicit, unambiguous "end of file" signal (the NUL byte), rather than leaving `l.ch` in some undefined state. Updating `position` to the just-used `readPosition`, then bumping `readPosition`, keeps the two pointers correctly in sync on every call.

---

## Q3. Update `New()` so the lexer is fully initialized (`l.ch`, `l.position`, `l.readPosition` all set) before anyone calls `NextToken()`.

**A.**
```go
func New(input string) *Lexer {
	l := &Lexer{input: input}
	l.readChar()
	return l
}
```
**Why:** Without this call, a freshly constructed `*Lexer` would have `l.ch` at its zero value instead of the actual first character of input. Calling `readChar()` once inside `New()` guarantees the lexer is in a fully working state the moment it's created, so the very first call to `NextToken()` behaves correctly.

---

## Q4. Write the first working version of `NextToken()`, covering `= ; ( ) , + { }` and EOF, plus the `newToken()` helper it uses.

**A.**
```go
package lexer

import "monkey/token"

func (l *Lexer) NextToken() token.Token {
	var tok token.Token

	switch l.ch {
	case '=':
		tok = newToken(token.ASSIGN, l.ch)
	case ';':
		tok = newToken(token.SEMICOLON, l.ch)
	case '(':
		tok = newToken(token.LPAREN, l.ch)
	case ')':
		tok = newToken(token.RPAREN, l.ch)
	case ',':
		tok = newToken(token.COMMA, l.ch)
	case '+':
		tok = newToken(token.PLUS, l.ch)
	case '{':
		tok = newToken(token.LBRACE, l.ch)
	case '}':
		tok = newToken(token.RBRACE, l.ch)
	case 0:
		tok.Literal = ""
		tok.Type = token.EOF
	}

	l.readChar()
	return tok
}

func newToken(tokenType token.TokenType, ch byte) token.Token {
	return token.Token{Type: tokenType, Literal: string(ch)}
}
```
**Why:** The `switch` on `l.ch` is the simplest possible way to turn "what character am I looking at" into "what token is this" — each single-character symbol maps directly to one token type. `newToken()` is factored out because every branch needs the same two-step conversion (wrap the type, stringify the byte), so writing it once avoids repeating it in every case. `l.readChar()` is called *after* the token is built, not before, so `l.ch` still reflects the character actually being turned into a token; calling it any earlier would build the token from the wrong character.

---

## Q5. Add a `default` branch to `NextToken()` that recognizes identifiers, along with the `isLetter()` and `readIdentifier()` functions it depends on. (At this point, don't worry yet about telling keywords apart from user identifiers — just get the raw identifier text out.)

**A.**
```go
func (l *Lexer) NextToken() token.Token {
	var tok token.Token

	switch l.ch {
	// [... existing single-character cases ...]
	default:
		if isLetter(l.ch) {
			tok.Literal = l.readIdentifier()
			return tok
		} else {
			tok = newToken(token.ILLEGAL, l.ch)
		}
	}

	l.readChar()
	return tok
}

func (l *Lexer) readIdentifier() string {
	position := l.position
	for isLetter(l.ch) {
		l.readChar()
	}
	return l.input[position:l.position]
}

func isLetter(ch byte) bool {
	return 'a' <= ch && ch <= 'z' || 'A' <= ch && ch <= 'Z' || ch == '_'
}
```
**Why:** Not every character matches one of the fixed single-character cases, so a `default` branch is needed to catch everything else — starting with identifiers and keywords, which are made up of letters. `isLetter` also accepts `_` so identifiers like `foo_bar` are read as one identifier instead of being split at the underscore. The early `return tok` matters here: `readIdentifier()` already calls `readChar()` in its own loop to consume the whole word, so by the time it returns, `l.ch` is already sitting on the character *after* the identifier — the shared `l.readChar()` at the bottom of the function would skip that character if it ran again. Anything that isn't a letter and isn't already handled falls through to `token.ILLEGAL`, since the lexer genuinely doesn't know what to do with it yet.

---

## Q6. In `token.go`, add the `keywords` map and a `LookupIdent()` function that tells keywords apart from plain identifiers.

**A.**
```go
// token/token.go

var keywords = map[string]TokenType{
	"fn":  FUNCTION,
	"let": LET,
}

func LookupIdent(ident string) TokenType {
	if tok, ok := keywords[ident]; ok {
		return tok
	}
	return IDENT
}
```
**Why:** `readIdentifier()` can only get the raw text of a word — it has no way of knowing whether that word is a language keyword (`let`, `fn`) or a name the programmer made up. `LookupIdent` closes that gap by checking a lookup table. It lives in the `token` package rather than `lexer` because it's really a fact about the language's vocabulary, not about the mechanics of scanning characters.

---

## Q7. Update the `default` branch in `NextToken()` to use `LookupIdent()` so `tok.Type` is set correctly for both keywords and identifiers.

**A.**
```go
default:
	if isLetter(l.ch) {
		tok.Literal = l.readIdentifier()
		tok.Type = token.LookupIdent(tok.Literal)
		return tok
	} else {
		tok = newToken(token.ILLEGAL, l.ch)
	}
```
**Why:** This is the piece that actually completes identifier/keyword lexing — until now, every word (including `let` and `fn`) would have been typed as a plain identifier. Feeding the literal through `LookupIdent` lets the lexer correctly emit `token.LET`/`token.FUNCTION` for keywords while still falling back to `token.IDENT` for everything else.

---

## Q8. A test on real Monkey source (`let five = 5; ...`) failed on the whitespace between `let` and `five`. Write `skipWhitespace()` and wire it into `NextToken()` to fix it.

**A.**
```go
func (l *Lexer) NextToken() token.Token {
	var tok token.Token

	l.skipWhitespace()

	switch l.ch {
	// [...]
	}
	// [...]
}

func (l *Lexer) skipWhitespace() {
	for l.ch == ' ' || l.ch == '\t' || l.ch == '\n' || l.ch == '\r' {
		l.readChar()
	}
}
```
**Why:** The failing test showed the lexer trying to lex the space character itself, since nothing was yet skipping over it. In Monkey, whitespace only separates tokens and carries no meaning of its own, so the fix is to consume it silently before the `switch` ever runs — calling `skipWhitespace()` first means `l.ch` is guaranteed to be a real, meaningful character (or EOF) by the time token recognition begins.

---

## Q9. Add support for integer literals: write `readNumber()` and `isDigit()`, and extend the `default` branch to use them.

**A.**
```go
func (l *Lexer) NextToken() token.Token {
	var tok token.Token

	l.skipWhitespace()

	switch l.ch {
	// [...]
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
	// [...]
}

func (l *Lexer) readNumber() string {
	position := l.position
	for isDigit(l.ch) {
		l.readChar()
	}
	return l.input[position:l.position]
}

func isDigit(ch byte) bool {
	return '0' <= ch && ch <= '9'
}
```
**Why:** With whitespace out of the way, the next thing the test input hits is a plain integer like `5`, which the lexer still can't handle. `readNumber()` mirrors `readIdentifier()` almost exactly — same start/end bookkeeping, just swapping `isDigit` for `isLetter` — since scanning a run of digits is the same kind of problem as scanning a run of letters. It's kept as its own small function rather than merged with `readIdentifier()` for simplicity and readability, even though that means a bit of duplication.
