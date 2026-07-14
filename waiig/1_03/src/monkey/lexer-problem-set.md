# Lexer Problem Set

Parts 1–5 walk through the code you already have, in the same order the text introduces it. Parts 6–9 extend the lexer beyond what the text covers.

---

## Part 1 — The `Lexer` struct and `readChar()`

**Q1.1.** What do `position` and `readPosition` each point to?
**A.** `position` points to the index of the character currently stored in `l.ch`. `readPosition` points to the index of the *next* character to be read.

**Q1.2.** Why does the lexer need two separate pointers instead of just one?
**A.** Lexing sometimes needs to look ahead at what comes after the current character without losing track of where "current" actually is. One index can't represent both "where I am" and "where I'm about to look" simultaneously.

**Q1.3.** What does `readChar()` set `l.ch` to once it runs out of input, and why that value specifically?
**A.** It sets `l.ch = 0` (the NUL byte) — chosen because `0` can never appear as real source code content, so it works as an unambiguous "no more input" signal.

**Q1.4.** Why would using `' '` (space) as the end-of-input sentinel instead of `0` be a bad idea?
**A.** Space is already a meaningful character (whitespace to skip). Reusing it for "end of input" would make genuine trailing whitespace indistinguishable from EOF.

---

## Part 2 — The first version of `NextToken()` and `newToken()`

**Q2.1.** Why does `NextToken()` switch on `l.ch` rather than on, say, `l.position` or the whole remaining input string?
**A.** Because the decision of "which token is this" only ever depends on the single character currently under examination (plus, later, a one-character lookahead) — `l.ch` is exactly the piece of state that determines it.

**Q2.2.** Why is `newToken(tokenType, ch)` written as a separate helper instead of writing `token.Token{Type: ..., Literal: string(ch)}` inline in every `case`?
**A.** It removes repetition — every single-character case needs the exact same two-step conversion (wrap the tokenType, stringify the one byte), so factoring it into one function avoids repeating that boilerplate seven-plus times.

**Q2.3.** Why is `l.readChar()` called once, after the `switch` statement finishes, rather than at the top of `NextToken()`?
**A.** Because `l.ch` needs to still hold the *current* character while the `switch` inspects it and builds `tok`. Only once the token has been built is it safe to advance — calling `readChar()` before the switch would mean building the token from the wrong character.

---

## Part 3 — Identifiers and keywords

**Q3.1.** What does `isLetter` check, and why does it also return `true` for `_`?
**A.** It checks whether a byte is an uppercase or lowercase ASCII letter. It also allows `_` so that identifiers like `foo_bar` are read as a single identifier instead of being split at the underscore.

**Q3.2.** If you also allowed `!` and `?` in `isLetter` (to permit names like `valid?`), what's the tradeoff?
**A.** It would let more expressive identifiers through, but it risks colliding with other lexer logic that treats `!` as its own token (e.g. `BANG`, or `!=`) — that logic would need to run *before* falling into the identifier path, or `!`/`?` would need to be excluded from at least the leading character check.

**Q3.3.** How does `readIdentifier()` know where the identifier starts and ends?
**A.** It captures `l.position` as the start before advancing, then calls `l.readChar()` in a loop for as long as `isLetter(l.ch)` is true, and finally slices `l.input[position:l.position]` — `l.position` has moved to just past the last letter by the time the loop exits.

**Q3.4.** Why does `LookupIdent` live in the `token` package rather than the `lexer` package?
**A.** Deciding whether a word is a keyword or a user-defined identifier is a fact about the language's vocabulary (token types), not about the mechanics of scanning characters — other consumers (like a future parser) may need that vocabulary independent of anything lexing-specific.

**Q3.5.** Why does the `default` branch's identifier case end with an early `return tok` instead of falling through to the shared `l.readChar()` at the bottom of `NextToken()`?
**A.** Because `readIdentifier()` already advances `l.ch` past the whole identifier internally, in its own loop. By the time it returns, `l.ch` already holds whatever character comes *after* the identifier. A second `readChar()` call would skip that next character entirely.

**Q3.6.** What token would go missing if that early return were deleted?
**A.** Whatever character immediately follows the identifier — e.g. lexing `five;` would silently drop the `;` token.

---

## Part 4 — Skipping whitespace

**Q4.1.** Why does the lexer skip over whitespace instead of turning it into its own token type?
**A.** In Monkey, whitespace only serves to separate tokens from each other — it carries no meaning of its own, so there's nothing useful for a parser to do with a "whitespace token." Skipping it keeps the token stream simpler for the parsing stage that comes next.

**Q4.2.** A teammate wrote this and now `TestNextToken` hangs on some inputs:
```go
func (l *Lexer) skipWhitespace() {
	for l.ch == ' ' || l.ch == '\t' || l.ch == '\n' || l.ch == '\r' {
	}
}
```
What's missing, and why does it only hang on inputs that *contain* whitespace rather than all inputs?
**A.** The loop body is missing `l.readChar()`. If `l.ch` never equals one of the four whitespace bytes, the condition is false immediately and the (empty) body never runs — no hang. Only once `l.ch` genuinely is whitespace does the body execute, and without `readChar()` advancing `l.ch`, the condition never becomes false.

**Q4.3.** What would go wrong if `skipWhitespace()` were called *after* the `switch` statement instead of before it?
**A.** The `switch` would evaluate `l.ch` while it might still be whitespace. None of the `case` branches match a space/tab/newline, so it would fall into `default`, where neither `isLetter` nor `isDigit` match either — producing a spurious `ILLEGAL` token for every whitespace character instead of silently passing over it.

---

## Part 5 — Numbers

**Q5.1.** `readNumber()` looks almost identical to `readIdentifier()`. What's the one meaningful difference?
**A.** `readNumber()` loops on `isDigit(l.ch)` where `readIdentifier()` loops on `isLetter(l.ch)` — otherwise the start/end bookkeeping (`position := l.position`, loop, slice) is the same pattern.

**Q5.2.** Given how similar they are, why not generalize `readIdentifier`/`readNumber` into one function that takes a character-classifying function as a parameter?
**A.** It's a deliberate simplicity tradeoff — two small, obviously-named functions are easier to read and reason about than one generic function passed a predicate, especially for an educational codebase where clarity outweighs avoiding a small amount of duplication.

**Q5.3.** What numeric formats does `readNumber()` intentionally *not* support?
**A.** Floats, hexadecimal notation, and octal notation — it only reads plain base-10 integers.

---

## Part 6 — Extending the lexer: new single-character operators

New symbols to support: `-` `/` `*` `<` `>` `!`

**Q6.1.** What should the new `token.go` constants be, and what string value should each hold?
**A.**
```go
MINUS    = "-"
SLASH    = "/"
ASTERISK = "*"
LT       = "<"
GT       = ">"
BANG     = "!"
```

**Q6.2.** What's the shape of the `case` you'd add to `NextToken()` for `-`?
**A.**
```go
case '-':
    tok = newToken(token.MINUS, l.ch)
```

**Q6.3.** Why can `newToken()` be reused unchanged for all six of these?
**A.** Because all six are single-byte tokens, and `newToken(tokenType, ch byte)` was already built to wrap exactly one byte.

**Q6.4.** What's a minimal test input that exercises all six new tokens plus at least one existing one?
**A.** Something like `` `-/*<>!5;` `` — a short string mixing the new operators with a digit and a delimiter you already lex correctly.

**Q6.5.** Does adding these six require touching `skipWhitespace`, `readIdentifier`, or `readNumber`?
**A.** No — single-character symbol tokens are handled entirely inside the `switch`.

---

## Part 7 — Extending the lexer: two-character tokens (`==`, `!=`)

**Q7.1.** Why can't a single-character `switch` on `l.ch` correctly lex `==` on its own?
**A.** By the time execution reaches `case '='`, the lexer only knows the *current* character — it has no information yet about whether the next character is also `=`, so it can't decide between `ASSIGN` and `EQ`.

**Q7.2.** What new method should you add to `*Lexer`, and how must it differ from `readChar()`?
**A.** A `peekChar()` method. Unlike `readChar()`, it must look at `l.input[l.readPosition]` **without** modifying `l.position`, `l.readPosition`, or `l.ch` — pure read-only lookahead.

**Q7.3.** What should `peekChar()` return if `l.readPosition` is past the end of input?
**A.** `0`, the same NUL sentinel used everywhere else for end-of-input.

**Q7.4.** What are the two new token constants, and their string values?
**A.**
```go
EQ     = "=="
NOT_EQ = "!="
```

**Q7.5.** Inside `case '='`, what condition decides between `=` and `==`?
**A.** `if l.peekChar() == '='` — true means it's `==`, false means it's plain `ASSIGN`.

**Q7.6.** Once you know it's `==`, what three things need to happen before building the token?
**A.**
1. Save the current character (`ch := l.ch`) before it's overwritten.
2. Call `l.readChar()` to advance onto and consume the second `=`.
3. Build the two-character literal, e.g. `literal := string(ch) + string(l.ch)`.

**Q7.7.** Why can't you just call `newToken(token.EQ, l.ch)` for the two-character case?
**A.** `newToken` only takes a single `byte` and stringifies it. `==` is two characters, so you build `token.Token{Type: token.EQ, Literal: literal}` directly instead.

**Q7.8.** Why is `peekChar()` necessary instead of just calling `readChar()` to check the next character?
**A.** If you used `readChar()` to peek and the next character turned out *not* to be `=`, you'd have already consumed and lost the character you still needed to lex normally on the very next call.

**Q7.9.** What test input would confirm both `EQ` and `NOT_EQ` work, alongside `=`/`!` appearing alone?
**A.** Something like `5 == 5; 5 != 10; x = 5; !true` covers both the two-character and single-character forms.

---

## Part 8 — Extending the lexer: new keywords

New keywords: `true` `false` `if` `else` `return`

**Q8.1.** What are the new `token.go` constants?
**A.**
```go
TRUE   = "TRUE"
FALSE  = "FALSE"
IF     = "IF"
ELSE   = "ELSE"
RETURN = "RETURN"
```

**Q8.2.** Where do these five entries get added?
**A.** Into the `keywords` map that `LookupIdent` checks, e.g. `"true": TRUE, "false": FALSE, "if": IF, "else": ELSE, "return": RETURN`.

**Q8.3.** Does `NextToken()` need any new `case` branches for these keywords?
**A.** No — they're words made of letters, so they flow through `readIdentifier()` → `LookupIdent()` exactly like `let` and `fn` already do.

**Q8.4.** If you find yourself editing `NextToken()` to make a new keyword work, what does that suggest?
**A.** That keyword recognition has leaked out of the table-driven design — the lexer isn't supposed to know the specific list of keywords at all; that's entirely `token.LookupIdent`'s job.

**Q8.5.** What's a test input exercising all five new keywords together?
**A.**
```
if (5 < 10) {
    return true;
} else {
    return false;
}
```

---

## Part 9 — Stretch goals: floats, strings, unicode

The original text explicitly calls out integers-only, ASCII-only support as simplifications "for the educational aim and limited scope of this book" — these pick up exactly where it leaves off.

**Q9.1.** Right now, how does the lexer tokenize `-5`?
**A.** As two separate tokens: `MINUS` followed by `INT("5")`.

**Q9.2.** What's an argument for leaving that behavior as-is, rather than making the lexer emit a single `INT("-5")` token?
**A.** The lexer has no surrounding context to know whether `-` is unary negation or subtraction (`a - 5` looks identical up to this point). Leaving disambiguation to the parser, which sees the full expression, is simpler and more correct.

**Q9.3.** To support floats like `3.14`, what would `readNumber()` need to handle differently?
**A.** After consuming digits, check for a single `.` followed by more digits and keep consuming — while making sure a second `.` (as in `3.14.15`) isn't silently accepted as part of the same number.

**Q9.4.** Would `INT` still be an appropriate token type name for float literals?
**A.** No — a separate type (e.g. `FLOAT`) better reflects that the literal isn't an integer, which matters later when the evaluator decides which Go numeric type to parse it into.

**Q9.5.** For string literals like `"hello world"`, what should the lexer do if it hits end-of-input before finding the closing `"`?
**A.** Stop and signal an error/illegal token rather than reading past the end of input or looping forever — an unterminated string is a lexing error, not something to pass through silently.

**Q9.6.** What would have to change about `Lexer.ch`'s type to support Unicode identifiers like `café`?
**A.** `ch` would need to change from `byte` to `rune`, since non-ASCII characters can span multiple bytes in UTF-8 — single-byte indexing (`l.input[l.readPosition]`) would no longer correctly identify one character.

---

## Self-check

- Can you answer every question in Parts 1–5 without looking anything up? Those are about code you already have.
- Can you add one brand-new single-character token end-to-end (constant → case → test) from memory (Part 6)?
- Can you explain, out loud, why `==` needs `peekChar()` but `-` doesn't (Part 7)?
