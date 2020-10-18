# lang
Project to learn concepts from *INF225 (Program Translation)* that will hopefully result in some fun programming language.

### The current grammar that I am trying to implement
```
<exp> := <val> { <bop> <val> }
<val> := "(" <exp> ")"
<val> := <num>
<bop> ::= "+" | "*"
<num> ::= simple integers
```

- `{ _ }` means zero or more
- Operators are left-associative

### Resources
- The structure of the tokenizer and parser is inspired by this blog post: https://blog.gopheracademy.com/advent-2014/parsers-lexers/.
- This site was very helpful when constructing the grammar: http://www.engr.mun.ca/~theo/Misc/exp_parsing.htm.
