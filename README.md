# Call Me Maybe
A simple compiler implemented in Go. 
It compiles strings of [this grammar](documentation/grammar.md) to x86 (Intel) assembly code. 
The assembly code is assembled with [NASM](https://www.nasm.us/) and linked with [GCC](https://gcc.gnu.org/).

## Features

- Type safety
- Pure functions (except for IO)
- Higher order functions
- Characters, ints, booleans, structs, strings and arrays
- Basic arithmetic and logic
- Loop and if

## Installation
- Use ubuntu (other linux distributions will probably work as well)
- Install nasm, gcc, git and go
- Clone this repository and run `go build -o cmm`
- Install the [vscode plugin](https://marketplace.visualstudio.com/items?itemName=petterdaae.callmemaybe)
- `./cmm build <source>` will output an executable named `out` for the code in the `<source>` file

## Examples

```
println "Hello world!"
```

```
fruit = <string, 3>["banana", "apple", "orange"]
apple = ?fruit[1]
```

```
numbers = <int, 5>[1, 2, 3, 4, 5]
i = 0
sum = 0
loop i < 5 {
    sum = sum + ?numbers[i]
    i = i + 1
}
println sum
```

```
double = | x int | int {
    return x * 2
}

apply = | f func<int, int>, x int | int {
    return #f(x)
}

println #apply(double, 2)
```

