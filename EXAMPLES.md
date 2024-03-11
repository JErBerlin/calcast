# Worked out examples

The following simple examples illustrate how the parsing algorithm in `parse.go` operates. We will focus particularly on `parseBinary()`. Despite being less than 20 lines of code, this function accomplishes most of the parsing work.

An Abstract Syntax Tree (AST) is a binary tree that represents syntactical relationships between different parts of an expression. Even when an explicit tree-like data structure is not defined within our parser's code, the conceptual model of an AST still underpins the parsing process.

This process, especially when parsing expressions as seen in our examples, inherently builds a hierarchical structure that mirrors the properties of an AST. This is because the parser recognizes the precedence and associativity rules of the language's grammar. It nests expressions within each other based on these rules. 

For instance, in an expression like `1 + 2 * 3`, the parser inherently understands that `2 * 3` must be evaluated before adding `1` due to the higher precedence of multiplication over addition. This understanding leads to a hierarchical parsing process that, while not explicitly creating nodes and edges in a tree data structure, conceptually builds an AST where `+` is the root with `1` and the result of `2 * 3` as its children.

The function calls and recursion in the parsing algorithm trace the paths and nodes of this implicit AST. Each call to parse a component of the expression —whether it's a number, a unary operation, or a binary operation— acts as a traversal down the tree to a child node. The return from these calls, moving back up the call stack, represents moving back up the tree, assembling the nodes into a complete structure along the way.

## Example 1: 1 + 2

Final Tree:

```  
  +
 / \
1   2
```

### Initial Setup
- The parser starts by reading the first part of the expression, which is 1. 
- This is the initial state, with the lexer pointing to 1.

#### parseUnary is called:
1. This call is to handle any unary operators (like - or + before a number), but since 1 is a number (and not preceded by any unary operator), parseUnary essentially delegates to parsePrimary to process 1.
2. Inside parsePrimary, 1 is identified as a numeric literal. The function returns an expression representing the number 1, and this expression is assigned to left.
3. The lexer's token is now updated to + because lex.next() is called after parsing 1.

#### Encountering +:
1. The parser checks the priority of the current token (+), which is 1. 
2. Since the loop in parseBinary is prepared to handle operators with at least the priority of prio0 (which is also 1 at the start), it proceeds.

#### Processing the operator +:
1. The operator (+) is saved in a temporary variable, op.
2. lex.next() is invoked again, advancing the lexer and updating its token to point to 2, the next part of the expression.

#### Right side parsing (2):
1. The parser, through a recursive call to parseBinary with prio+1 (making it 2 to ensure higher precedence operations are evaluated first), attempts to parse the right side of the +. However, it encounters 2, which, like 1, is processed through parseUnary and parsePrimary, effectively identifying it as a numeric literal.
2. This numeric literal (2) becomes the right operand for the + operation.

#### Constructing the binary expression for 1 + 2:
1. With both left (the number 1) and right (the number 2) operands parsed and the operator (+) identified, the parser constructs a binary expression representing 1 + 2.
2. This binary expression is the final result of parsing the input 1 + 2, with the left operand being 1, the operator being +, and the right operand being 2.

## Construction Steps:

### Step 1: Parse 1

```
1
```

### Step 2: Encounter +, Parse 2

```
  +
 / \
1   2
```

## Example 2: 2 + 1*2

Final Tree:
```
    +
   / \
  2   *
     / \
    1   2
```

### Initial Setup
- Input: The expression 2 + 1*2.
- The lexer tokenizes the input. Initially, the lexer's token is set to the first character/token in the expression, which is 2.

#### First parseBinary call
1. parseUnary: Called to parse the left operand before encountering any operators. Since 2 is a primary (numeric) value, parseUnary essentially delegates to parsePrimary, setting left to the numeric expression representing 2.
2. lexer.token: Now points to + after parseUnary consumes the 2.
3. priority of '+': Checked and found to be 1.
4. for loop: Continues because the priority of + is equal to the initial prio0 (1 in this case).

#### Handling +
1. op = lexer.token: The operator is set to +.
2. lex.next(): Consumes the +, moving the lexer to the next token, which is 1.
3. Right-hand side parsing: Calls parseBinary recursively with prio+1 (2 in this case), to ensure that any operations on the right with equal or higher precedence are evaluated first.

#### Inside Right-hand Side parseBinary for 1*2
1. parseUnary for 1: Similar to the initial parseUnary call, 1 is parsed as a primary numeric value, setting a temporary left to 1.
2. lexer.token: Now points to *.
3. priority of '*': Checked and found to be 2, which is higher than the current prio0 for this context, allowing the loop to continue.
4. op = lexer.token: The operator is set to *.
5. lex.next(): Consumes the *, and the lexer moves to 2.
6. Right-hand side parsing for *: A recursive call to parseBinary is made with prio+1 (3 in this case), but since there are no more operators with higher precedence, this call will essentially end up parsing 2 as a primary numeric value and return it as the right operand for *.

#### Finalizing 1*2
- Construct binary expression: A binary expression object is created with * as the operator, 1 as the left operand, and 2 as the right operand. This binary expression represents the 1*2 sub-expression.

#### Returning to the Top Level +
- Construct binary expression for +: The left operand is the 2 parsed initially, the operator is +, and the right operand is the binary expression representing 1*2.

## Construction Steps:

### Step 1: Parse 2

```
2
```
(waiting for right operand)

### Step 3: Parse 1, then *, then 2 for the right operand

```
    +
   / \
  2   *
     / \
    1   2
```

## Example 3: 2*1 + 2

Final Tree:

```
    +
   / \
  *   2
 / \
2   1
```

#### Beginning with 2*1
1. parseUnary: Called first to parse the left part of the expression, which is 2. Since it's a number, it directly translates into a numeric expression (num type, for example) and assigns this to left. This is because the first token 2 is immediately recognized as a primary expression.

2. lexer.token: After consuming 2, the lexer's token is updated to *.

3. priority of '*': Determined to be 2 because multiplication has higher precedence. The parser is now ready to process the binary operation.

4. Entering the first for loop with *: The current token is *, so the loop proceeds since its priority matches the condition.

5. op = lexer.token: Assigns the operator * to op.

6. lex.next(): Moves the lexer forward, consuming *. The next token is 1, which is the right operand for the multiplication operation.

7. Right-hand side parsing for *: Another call to parseBinary is made, this time to parse the expression starting with 1. Given the precedence of the current operator *, this recursive call aims to resolve any higher-priority operations to the right, but it finds only 1, a primary expression, which becomes the right operand.

8. Construct binary expression for * : The left operand 2, the operator * , and the right operand 1 are combined into a binary expression representing 2*1. This binary expression is now the left side of the top-level expression being parsed.

#### Transition to + 2

1. lexer.token after consuming 1: Now points to +, since the parser has moved past the 2*1 expression.

2. priority of '+': It's 1, indicating a lower precedence compared to multiplication. This shift signifies moving to a broader scope in the expression hierarchy.

3. Processing +: The loop continues because + matches the outer scope's expected precedence. The parser is effectively at the top-level expression again, with 2*1 as the accumulated left side.

4. op = lexer.token: The operator + is assigned to op.

5. lex.next(): Advances the lexer, consuming +. The next token is 2, which will serve as the right operand for the addition operation.

6. Right-hand side parsing for +: Invokes parseBinary again, but this time, since it's parsing the right operand of + and there are no further operators of equal or higher precedence, it effectively parses the 2 as a primary numeric value, setting it as the right operand.

7. Construct binary expression for +: Finally, the parser combines the previously constructed left side (the binary expression representing 2 * 1 ), the operator +, and the right side (2) into a new binary expression. This final expression represents the entire 2 * 1 + 2 input.


## Construction Steps:

### Step 1: Parse 2*1

```
  *
 / \
2   1
```

### Step 2: Encounter +, then parse 2

```
    +
   / \
  *   2
 / \
2   1
```

## Example 4: 1 + 2 + 3 + 4

Final Tree:

```
        +
       / \
      +   4
     / \
    +   3
   / \
  1   2
```

#### Parse 1
1. parseUnary is called to process potential unary operations but directly handles 1 as a numeric literal.
2. The numeric value 1 is assigned to left.

#### Encounter +, Prepare to Parse 2
1. lexer.token updates to + after consuming 1.
2. parseBinary recognizes + and prepares to parse the next part of the expression.

#### Parse 2, Building 1 + 2
1. parseUnary processes 2 similarly as a numeric literal.
2. A binary expression for 1 + 2 is formed with left being 1, op being +, and right being 2.

#### Encounter second +, Prepare to Parse 3
1. lexer.token updates to the second +.
2. The existing 1 + 2 expression is now the left operand of a new binary expression.

#### Parse 3, Extending to 1 + 2 + 3
1. parseUnary processes 3 as a numeric literal.
2. A new binary expression is formed extending the tree to include + and 3.

#### Encounter third +, Prepare to Parse 4
1. lexer.token updates to the third +.
2. The 1 + 2 + 3 expression forms the left operand of yet another binary expression.

#### Parse 4, Finalizing 1 + 2 + 3 + 4
1. parseUnary treats 4 as a numeric literal.
2. The final binary expression 1 + 2 + 3 + 4 is constructed

Construction Steps:

### Step 1: Parse 1 + 2

```
  +
 / \
1   2
```

### Step 2: Add 3 to the tree

```
    +
   / \
  +   3
 / \
1   2
```

### Step 3: Add 4 to the tree

```
      +
     / \
    +   4
   / \
  +   3
 / \
1   2
```