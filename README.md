# lang
Project to learn concepts from *INF225 (Program Translation)* that will hopefully result in some fun programming language.

### The current grammar that I am trying to implement
```
<seq>        := { <stmt> }
<stmt>         := <assign> | <prinln>
<assign>     := <identifier> "=" <exp>
<println>    := "println" <exp>

<exp>        := <val> { <bop> <val> }
<val>        := "(" <exp> ")" | <num> | <let> | <identifier>
<bop>        := "+" | "*"
<num>        := simple integers
<let>        := "let" <identifier> "=" <exp> "in" <exp>
<identifier> := simple words, only letters
```

- `{ _ }` means zero or more
- Operators are left-associative

### Resources
- The structure of the tokenizer and parser is inspired by this blog post: https://blog.gopheracademy.com/advent-2014/parsers-lexers/.
- This site was very helpful when constructing the grammar: http://www.engr.mun.ca/~theo/Misc/exp_parsing.htm.
