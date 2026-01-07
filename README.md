# Mooshy Interpreter

A tree-walking interpreter for a dynamically-typed programming language, implemented in Go. The interpreter follows a standard architecture: lexical analysis → parsing → AST construction → evaluation.

## Implementation

### Lexer
- Hand-written lexer that performs tokenization via character stream processing
- Supports single-character tokens (`+`, `-`, `*`, `/`, etc.) and multi-character tokens (identifiers, integers, strings)
- Implements lookahead for disambiguation (e.g., `=` vs `==`)

### Parser
- Recursive descent parser using Pratt parsing (precedence climbing) for expression handling
- Produces an Abstract Syntax Tree (AST) with distinct node types for statements and expressions
- Operator precedence hierarchy: `LOWEST < EQUALS < LESSGREATER < SUM < PRODUCT < PREFIX < POSTFIX < CALL`

### AST
- Node interface hierarchy: `Node` ← `Statement`, `Expression`
- Statement types: `LetStatement`, `ReturnStatement`, `ExpressionStatement`, `BlockStatement`
- Expression types: `IntegerLiteral`, `StringLiteral`, `ArrayLiteral`, `InfixExpression`, `PrefixExpression`, `PostfixExpression`, `IfExpression`, `FunctionLiteral`, `CallExpression`

### Evaluator (Runtime)
- Tree-walking interpreter with environment-based scoping
- Object system implements value representation at runtime
- Object types: `Integer`, `String`, `Array`, `Function`, `ReturnValue`, `Null`, `Error`
- Function calls create new enclosed environments for lexical scoping
- Supports closures through environment chaining

## Language Specification

**Primitive Types:**
- `INTEGER`: 64-bit signed integers
- `STRING`: UTF-8 encoded strings
- `BOOL`: Boolean type
- `NULL`: null value singleton

**Composite Types:**
- `ARRAY`: heterogeneous, dynamically-sized lists
- `FUNCTION`: first-class function objects with closure support

**Syntax:**

Variable binding:
```mooshy
let x = 10;
let name = "string";
```

Functions:
```mooshy
let add = func(a, b) { return a + b; };
let fibonacci = func(n) { 
    if (n < 2) { return n; } 
    return fibonacci(n - 1) + fibonacci(n - 2); 
};
```

Control flow:
```mooshy
if (x > 10) { 
    return x; 
} else { 
    return 0; 
}

for (let i = 0; i < 10; i = i + 1) {
    // loop body
}
```

Arrays:
```mooshy
let arr = [1, 2, 3, 4];
let nested = [[1, 2], [3, 4]];
```

Higher-order functions:
```mooshy
let map = func(arr, f) {
    // applies f to each element
};
let result = map([1, 2, 3], func(x) { return x * 2; });
```

**Operators:**
- Arithmetic: `+`, `-`, `*`, `/`
- Comparison: `==`, `!=`, `<`, `>`
- Prefix: `-`, `!`

**Operator Precedence:** Implements standard mathematical precedence (PEMDAS/BODMAS)

## Project Structure

```
mooshy/
├── token/          # Token type definitions and lookup table
├── lexer/          # Lexical analysis (string → tokens)
├── parser/         # Syntax analysis (tokens → AST)
├── ast/            # AST node definitions
├── object/         # Runtime object system
├── evaluator/      # Tree-walking evaluation logic
├── repl/           # Read-Eval-Print Loop
└── main.go         # Entry point
```

## Testing

Test coverage for each pipeline stage:

**Lexer tests:** Token generation correctness
**Parser tests:** AST structure validation, precedence verification, error detection
**Evaluator tests:** Expression evaluation, control flow, function application, recursion, error propagation

Run tests:
```bash
go test ./...
```

## Usage

Start REPL:
```bash
go run main.go
```

Execute file:
```bash
go run main.go codeFile.mooshy
```

## Technical Notes

- **Scoping:** Lexical scoping with environment chains
- **Evaluation strategy:** Eager evaluation (call-by-value)
- **Memory management:** Relies on Go's garbage collector
- **Error handling:** Runtime errors propagate through `Error` object type
- **Immutability:** All values are immutable (functional semantics)

## Limitations

- No type system or type checking
- Single-threaded execution model
- No standard library
- Limited error diagnostics (no line/column tracking)
- No optimization pass

## Future Work

- Hash map data structure
- Boolean type (currently uses truthy/falsy integer semantics)
- Bytecode compilation target
- Better error messages with source location
- Standard library (I/O, string manipulation, etc.)

## References

Implementation follows patterns from "Writing An Interpreter In Go" by Thorsten Ball.
