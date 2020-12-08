# lang
A simple compiler.

### The current grammar that I am trying to implement
```
<seq>         := { <stmt> }
<stmt>        := <assign> | <prinln> | <call> | <return> | <if>
<assign>      := <identifier> "=" <exp> | "_" "=" <exp>
<println>     := "println" <exp>
<return>      := "return" <exp>

<expr>        := <calculation> | <function> | <call>
<calculation> := <val> { <bop> <val> }
<val>         := "(" <calulation> ")" | <num> | <identifier> | <bool> | <uop> <num>
<bop>         := "+" | "*" | "<" | ">" | "==" | "-" | "/"
<uop>         := "-"
<num>         := sequence of digits
<identifier>  := words consisting og letters, digits and underscores, starting with a letter

<if>          := "if" <expr> "{" <seq> "}" 
<loop>        := "loop" <expr> "{" <seq> "}"
<bool>        := "true" | "false"

<function>    := <argList> "=>" [ <type> ] "{" <seq> "}"
<argList>     := "<" <recurse> ">" | "<" <recurse>, { <arg> "," } <arg> ">"
<arg>         := <identifier> <type>
<recurse>     := <identifier> | "_"

<call>        := "call" <identifier> [ "with" { <expr> "," } <expr> ]

<type>        := "int" | "empty" | "bool"
```

- `{ _ }` means zero or more
- `[ _ ]` means zero or one
- `|` means or  
- Operators are left-associative
- Should not be possible to ignore function return types

### TODO
- Negative numbers, subtraction and division
- Add an extra parameter to functions to make recursion possible
- Loop
- Refactor
- Characters, lists and strings
- IO
- Make it possible to return functions and pass functions as arguments
- Structs

### Resources
- The structure of the tokenizer and parser is inspired by this blog post: https://blog.gopheracademy.com/advent-2014/parsers-lexers/.
- This site was very helpful when constructing the grammar: http://www.engr.mun.ca/~theo/Misc/exp_parsing.htm.
- Nice x86 assembly reference: http://www.cs.virginia.edu/~evans/cs216/guides/x86.html.
