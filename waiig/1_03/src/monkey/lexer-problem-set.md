# Lexer Problem Set — Extending the Monkey Lexer (Q&A Format)

Based on your `lexer.go` / `token.go`. Each part is broken into small question/answer pairs — good for working through in order, or dropping straight into Anki.

---

## Part 1 — Understanding what you already have

**Q1.1.** What do `position` and `readPosition` each point to?
**A.** `position` points to the index of the character currently stored in `l.ch`. `readPosition` points to the index of the *next* character to be read.

**Q1.2.** Why does the lexer need two separate pointers instead of just one?
**A.** Because lexing sometimes requires looking ahead at what comes after the current character (e.g. checking if `=` is followed by another `=`) without losing track of where the "current" character actually is. One pointer can't represent both "where I am" and "where I'm about to look" at the same time.

**Q1.3.** What does `readChar()` set `l.ch` to when it runs out of input, and why that value specifically?
**A.** It sets `l.ch = 0` (the NUL byte). Zero is used as a sentinel because it's not a valid/meaningful character in the source code itself, so it can safely mean "no more input" without being confused with real content.

**Q1.4.** Why would using `' '` (space) as the end-of-input sentinel instead of `0` be a bad idea?
**A.** Because space is a real character the lexer already treats meaningfully (as whitespace to skip). Reusing it for "end of input" would make genuine trailing whitespace indistinguishable from EOF.

**Q1.5.** In `NextToken()`, why do `readIdentifier()` and `readNumber()` end with an early `return tok` instead of falling through to the `l.readChar()` at the bottom of the function?
**A.** Because `readIdentifier`/`readNumber` already call `readChar()` internally, in a loop, to consume the whole word or number. By the time they return, `l.ch` is already sitting on the *next* character after the identifier/number. An extra `readChar()` call would skip that character entirely.

**Q1.6.** What token would go missing from the output if that early `return` were deleted?
**A.** Whatever character immediately follows the identifier/number — e.g. lexing `five;` would silently drop the `;` token.

**Q1.7.** Why does `LookupIdent` live in the `token` package rather than the `lexer` package?
**A.** Because deciding whether a word is a keyword or a plain identifier is a fact about the *language's vocabulary* (token types), not about the mechanics of scanning characters. A future parser or other tool may need `LookupIdent`/`keywords` without needing anything else from `lexer`.

**Q1.8.** `isLetter` currently returns true for `_`. What does that one line actually enable?
**A.** It lets identifiers contain underscores, e.g. `foo_bar` is lexed as one identifier instead of being broken at the `_`.

**Q1.9.** If you added `!` and `?` to `isLetter` so names like `valid?` become legal identifiers, what existing token type would you need to double check for conflicts?
**A.** `BANG` (`!`) — if `!` is treated as a letter, the `case '!'` branch (or the `!=`/`BANG` lexing logic) would never fire, since `!` would now always be swallowed into `readIdentifier()` instead.

---

## Part 2 — Adding new single-character tokens

New symbols to support: `-` `/` `*` `<` `>` `!`

**Q2.1.** What should the new `token.go` constants be called, and what string value should each hold?
**A.**
```go
MINUS    = "-"
SLASH    = "/"
ASTERISK = "*"
LT       = "<"
GT       = ">"
BANG     = "!"
```

**Q2.2.** What's the general shape of the `case` you'd add to `NextToken()` for `-`?
**A.**
```go
case '-':
    tok = newToken(token.MINUS, l.ch)
```

**Q2.3.** Why can `newToken()` be reused unchanged for all six of these, without modification?
**A.** Because all six are single-byte tokens, and `newToken(tokenType, ch byte)` was already built to wrap exactly one byte into a `token.Token`.

**Q2.4.** What's a minimal test input string that would exercise all six new tokens plus at least one existing one (e.g. `;`)?
**A.** Something like `` `-/*5;` `` or `` `!-/*5;` `` — any short string mixing the new operators with a digit and a delimiter you already lex correctly.

**Q2.5.** After adding the six new cases, do you need to touch `skipWhitespace`, `readIdentifier`, or `readNumber`?
**A.** No — those are unaffected. Single-character symbol tokens are handled entirely inside the `switch`.

---

## Part 3 — Two-character tokens: `==` and `!=`

**Q3.1.** Why can't a single-character `switch` on `l.ch` correctly lex `==` on its own?
**A.** Because when execution reaches `case '='`, the lexer only knows about the *current* character. It has no information yet about whether the next character is also `=`, so it can't decide between emitting `ASSIGN` or `EQ`.

**Q3.2.** What new method do you need to add to `*Lexer` to solve this, and what should it do differently from `readChar()`?
**A.** A `peekChar()` method. Unlike `readChar()`, it must look at `l.input[l.readPosition]` **without** modifying `l.position`, `l.readPosition`, or `l.ch` — it's read-only lookahead.

**Q3.3.** What should `peekChar()` return if `l.readPosition` is past the end of input?
**A.** `0` (the same NUL sentinel used elsewhere), for consistency with how end-of-input is signaled.

**Q3.4.** What are the two new token constants needed, and their string values?
**A.**
```go
EQ     = "=="
NOT_EQ = "!="
```

**Q3.5.** Inside `case '='`, what condition determines whether you're looking at `=` or `==`?
**A.** `if l.peekChar() == '='` — if true, it's the two-character `==`; otherwise it's plain `ASSIGN`.

**Q3.6.** Once you know it's `==`, what three things need to happen before you can build the token?
**A.**
1. Save the current character (`ch := l.ch`) before it's overwritten.
2. Call `l.readChar()` to advance onto and consume the second `=`.
3. Build the two-character literal, e.g. `literal := string(ch) + string(l.ch)`.

**Q3.7.** Why can't you just call `newToken(token.EQ, l.ch)` for the two-character case?
**A.** Because `newToken` takes a single `byte` and wraps it as a one-character string. `==` is two characters, so you need to construct the `token.Token{Type: token.EQ, Literal: literal}` struct directly instead.

**Q3.8.** Why is `peekChar()` necessary here instead of just calling `readChar()` to check the next character?
**A.** Because if you call `readChar()` to "peek" and the next character turns out *not* to be `=`, you've already consumed and lost the character you needed to lex normally on the very next call. `peekChar()` lets you look without committing.

**Q3.9.** What test input would confirm both `EQ` and `NOT_EQ` work, alongside a case where `=` and `!` appear *without* a following `=`?
**A.** Something like `5 == 5; 5 != 10; x = 5; !true` covers both the two-character and single-character forms.

---

## Part 4 — Adding new keywords

New keywords: `true` `false` `if` `else` `return`

**Q4.1.** What are the new `token.go` constants?
**A.**
```go
TRUE   = "TRUE"
FALSE  = "FALSE"
IF     = "IF"
ELSE   = "ELSE"
RETURN = "RETURN"
```

**Q4.2.** Where do these five new entries get added?
**A.** Into the `keywords` map that `LookupIdent` checks — e.g. `"true": TRUE, "false": FALSE, "if": IF, "else": ELSE, "return": RETURN`.

**Q4.3.** Does `NextToken()` in `lexer.go` need any new `case` branches for these keywords?
**A.** No. They're just words made of letters, so they already flow through `readIdentifier()` → `LookupIdent()`, exactly like `let` and `fn` do.

**Q4.4.** If you find yourself editing `NextToken()` to make keywords work, what does that suggest?
**A.** That something about how identifiers/keywords are being distinguished is off — keyword recognition is supposed to be entirely table-driven through `LookupIdent`, not hardcoded into the lexer's switch statement. The lexer shouldn't need to know the specific list of keywords at all.

**Q4.5.** What's a test input that exercises all five new keywords together?
**A.**
```
if (5 < 10) {
    return true;
} else {
    return false;
}
```

---

## Part 5 — Debugging challenge

Buggy code:
```go
func (l *Lexer) skipWhitespace() {
	for l.ch == ' ' || l.ch == '\t' || l.ch == '\n' || l.ch == '\r' {
	}
}
```

**Q5.1.** What line is missing from the loop body?
**A.** `l.readChar()` — without it, `l.ch` never changes, so the loop condition never becomes false.

**Q5.2.** Why does this only hang on inputs that *contain* whitespace, rather than every input?
**A.** If `l.ch` is never a whitespace character to begin with, the `for` loop's condition is false immediately and the (empty) body never executes — so no infinite loop occurs. It's only once `l.ch` actually equals one of those four characters that the loop body runs and, lacking `readChar()`, never exits.

**Q5.3.** Suppose `skipWhitespace()` were correct, but called *after* the `switch` statement in `NextToken()` instead of before it. What would go wrong?
**A.** The `switch` would run against whatever character is currently under examination *before* whitespace is skipped — so whitespace itself would fall into the `default` branch and get treated as an identifier/illegal token attempt, instead of being silently skipped before token recognition even starts.

---

## Part 6 — Stretch goals (Q&A)

**Q6.1.** Right now, how does the lexer tokenize `-5`?
**A.** As two separate tokens: `MINUS` followed by `INT("5")`.

**Q6.2.** What's an argument for leaving that behavior as-is, rather than making the lexer emit a single `INT("-5")` token?
**A.** It keeps the lexer simple and pushes the decision of "is this a negation operator or a negative literal" to the parser, which has the surrounding context (e.g. `a - 5` vs `-5`) needed to disambiguate correctly — the lexer alone can't always tell.

**Q6.3.** To support floats like `3.14`, what would `readNumber()` need to handle differently?
**A.** After consuming digits, it would need to check for a single `.` followed by more digits, and continue consuming — while making sure a second `.` (as in `3.14.15`) is *not* silently accepted as part of the same number.

**Q6.4.** Would `INT` still be an appropriate token type name for float literals?
**A.** No — a separate type (e.g. `FLOAT`) would better reflect that the literal isn't an integer, which matters later for how the parser/evaluator interprets the value.

**Q6.5.** For string literals like `"hello world"`, what should the lexer do if it hits end-of-input before finding the closing `"`?
**A.** It should stop and produce some kind of error/illegal token rather than reading past the end of input or looping forever — an unterminated string is a lexing error, not valid input to just pass through silently.

**Q6.6.** What would have to change about `Lexer.ch`'s type to support Unicode identifiers like `café`?
**A.** `ch` would need to change from `byte` to `rune`, since non-ASCII characters can span multiple bytes in UTF-8 — `l.input[l.readPosition]` (single-byte indexing) would no longer correctly identify one character.

---

## Self-check

- Can you answer Q1.1–Q1.9 without looking anything up?
- Can you add one brand-new single-character token end-to-end (constant → case → test) from memory?
- Can you explain, out loud, why `==` needs `peekChar()` but `-` doesn't?
