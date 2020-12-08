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
<val>         := "(" <calulation> ")" | <num> | <identifier> | <bool>
<bop>         := "+" | "*" | "<" | ">" | "=="
<num>         := sequence of digits
<identifier>  := simple words, only letters

<if>          := "if" <expr> "{" <seq> "}" // Typecheck this later
<bool>        := "true" | "false"

<function>    := <argList> "=>" [ <type> ] "{" <seq> "}"
<argList>     := "<" ">" | "<" { <arg> "," } <arg> ">"
<arg>         := <identifier> <type>

<call>        := "call" <identifier> [ "with" { <expr> "," } <expr> ]

<type>        := "int" | "empty" | "bool"

--------------------------------------------------------------

Current terminals: letters, digits, "{", "}", "(", ")", "<", ">", "+", "*", ",", "=", 
    "call", "with", "int", "empty", "println", "=>"
```

- `{ _ }` means zero or more
- `[ _ ]` means zero or one
- `|` means or  
- Operators are left-associative
- Should not be possible to ignore function return types

### TODO
- Make it possible to return functions and pass functions as arguments
- Booleans
- If
- Loop
- Refactor  
- Strings & IO
- Typechecking  
- Structs

### Resources
- The structure of the tokenizer and parser is inspired by this blog post: https://blog.gopheracademy.com/advent-2014/parsers-lexers/.
- This site was very helpful when constructing the grammar: http://www.engr.mun.ca/~theo/Misc/exp_parsing.htm.
- Nice x86 assembly reference: http://www.cs.virginia.edu/~evans/cs216/guides/x86.html.
