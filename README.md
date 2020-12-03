# lang
A simple compiler.

### The current grammar that I am trying to implement
```
<seq>         := { <stmt> }
<stmt>        := <assign> | <prinln> | <call>
<assign>      := <identifier> "=" <exp>
<println>     := "println" <exp>

<expr>        := <calculation> | <function> | <call>
<exprList>    := { <expr> "," } <expr>
<calculation> := <val> { <bop> <val> }
<val>         := "(" <calulation> ")" | <num> | <identifier>
<bop>         := "+" | "*"
<num>         := sequence of digits
<identifier>  := simple words, only letters

<function>    := <argList> => [ <type> ] <block>
<argList>     := "<" ">" | <arg> | "<" { <arg> "," } <arg> ">"
<arg>         := <identifier> <type>
<block>       := "{" <seq> [ <exp> ] "}"

<call>        := "call" <identifier> [ "with" <exprList> ]

<type>        := "int" | "empty"

--------------------------------------------------------------

Current terminals: letters, digits, "{", "}", "(", ")", "<", ">", "+", "*", ",", "call", "with"
```

- `{ _ }` means zero or more
- `[ _ ]` means zero or one
- `|` means or  
- Operators are left-associative

### Resources
- The structure of the tokenizer and parser is inspired by this blog post: https://blog.gopheracademy.com/advent-2014/parsers-lexers/.
- This site was very helpful when constructing the grammar: http://www.engr.mun.ca/~theo/Misc/exp_parsing.htm.
- Nice x86 assembly reference: http://www.cs.virginia.edu/~evans/cs216/guides/x86.html.
