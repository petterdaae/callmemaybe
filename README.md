# lang
Project to learn concepts from *INF225 (Program Translation)* that will hopefully result in some fun programming language.

### The current grammar that I am trying to implement
```
<exp> ::= <exp> + <exp>
<exp> ::= <exp> * <exp>
<exp> ::= "(" <exp> ")">
<exp> ::= <num>
<exp> ::= "let" <identifier> "=" <exp> "in" <exp>
<num> ::= simple integers
<exp> ::= <identifier>
<identifier> ::= sequence of letters
```

Operators are currently right associative.

### Resources
The structure of the tokenizer and parser is inspired by this blog post: https://blog.gopheracademy.com/advent-2014/parsers-lexers/.
