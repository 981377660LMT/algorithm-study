## 1 我们可以进一步减少我们的二元运算符，甚至比这里做的还要多。你还可以消除哪些其他指令，编译器将如何应对它们的缺失？

Having both `OP_NEGATE` and `OP_SUBTRACT` is redundant. We can replace subtraction with negate-then-add:

```c
// Emit the operator instruction.
switch (operatorType) {
  // ...
  case TOKEN_PLUS:          emitByte(OP_ADD); break;
  case TOKEN_MINUS:         emitBytes(OP_NEGATE, OP_ADD); break; // <--
  case TOKEN_STAR:          emitByte(OP_MULTIPLY); break;
  case TOKEN_SLASH:         emitByte(OP_DIVIDE); break;
  default:
    return; // Unreachable.
}
```

Or we can replace negation with:

1. Push zero.
2. Compile the negate operand.
3. Subtract.

It's also possibly to simplify the comparison and equality instructions using some stack juggling and a bitwise operator. Fundamentally, you only need a single operation, an instruction that returns one of three values: "less", "equal", or "greater". Similar to the `compareTo()` methods in many languages or the `<=>` in Ruby. Once you have that, the other operators can be defined in terms of it.

## 2 我们可以通过添加更多与高级操作相对应的特定指令来提高字节码虚拟机的速度(we can improve the speed of our bytecode VM by adding more specific instructions that correspond to higher-level operations)。你会定义什么指令来加速我们在本章中添加支持的用户代码类型？

Many other instruction sets define dedicated instructions for common small integer constants. 0, 1, 2, and -1 are good candidates.

A few arithmetic operations have common constant operands. For those cases, it may be worth adding instructions for them: incrementing and decrementing by one are the main ones. But maybe even doubling comes up enough to warrant it.

Likewise, comparisons to certain numbers are also common and can be encoded directly in a single instruction instead of needing to load the number from a constant and then use the comparison instruction. Many CPU instruction sets can compare a number with zero in a single instruction.

There's been some research into "superinstructions" -- automated or manual techniques for defining instructions that represent a sequence of common simpler instructions. There is a point of diminishing returns because eventually you run out of opcodes. You can use larger opcodes (16 bits, etc.), but then that slows down dispatch overall because now your code is larger.
