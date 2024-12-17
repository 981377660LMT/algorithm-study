## 1 每次遇到标识符时，编译器都会将全局变量的名称以字符串形式添加到常量表中。每次都会创建一个新的常量，即使该变量名已经在常量表的前一个槽中。对此进行优化。与运行时相比，你的优化对编译器的性能有什么影响？这样的权衡是否正确？

The optimization is pretty straightforward. When adding a string constant, we
look in the constant table to see if that string is already in there. The
interesting question is how. The simplest implementation is a linear scan over
the existing constants.
这个优化相当直接。当添加一个字符串常量时，`我们会在常量表中查找该字符串是否已经存在`。一个有趣的问题是如何实现这一点。最简单的实现是对现有常量进行线性扫描。

But that means compilation time is quadratic in the number of unique identifiers
in the chunk. While that's fine for relatively small programs, users have a
habit of writing larger programs than we ever anticipated. Virtually every
algorithm in the compiler that isn't linear is potentially a performance
problem.
但这意味着编译时间与块中唯一标识符的数量呈二次方增长。对于相对较小的程序来说，这没问题，但用户往往会编写比我们预期更大的程序。`编译器中几乎所有非线性的算法都可能成为性能问题。`

Fortunately, we have a way of looking up strings in constant time -- a hash
table. So, in the compiler, we add a hash table that keeps track of the
identifier constants that have already been added. Each key is an identifier,
and its value is the index of the identifier in the constant table.
幸运的是，我们有一种以常数时间查找字符串的方法——哈希表。因此，在编译器中，`我们添加了一个哈希表，用于跟踪已经添加的标识符常量`。每个键是一个标识符，其值是该标识符在常量表中的索引。

In compiler.c, add a module variable:

```c
Table stringConstants;
```

In `compile()`, we initialize and tear it down:

```c
bool compile(const char* source, Chunk* chunk) {
  initScanner(source);

  compilingChunk = chunk;
  parser.hadError = false;
  parser.panicMode = false;
  initTable(&stringConstants); // <--

  advance();

  while (!match(TOKEN_EOF)) {
    declaration();
  }

  endCompiler();
  freeTable(&stringConstants); // <--
  return !parser.hadError;
}
```

When adding an identifier constant, we look for it in the hash table first:
在添加标识符常量时，我们首先在哈希表中查找它：

```c
static uint8_t identifierConstant(Token* name) {
  // See if we already have it.
  ObjString* string = copyString(name->start, name->length);
  Value indexValue;
  if (tableGet(&stringConstants, string, &indexValue)) {
    // We do.
    return (uint8_t)AS_NUMBER(indexValue);
  }

  uint8_t index = makeConstant(OBJ_VAL(string));
  tableSet(&stringConstants, string, NUMBER_VAL((double)index));
  return index;
}
```

That's pretty simple. Compiling an identifier is still (amortized) constant
time, though with slightly worse constant factors. In return, we use up fewer
constant table slots. We don't actually save memory from redundant strings
because clox already interns all strings. But the smaller table is nice.

_Note that we leak memory for the identifier string in `identifierConstant()`
if the name is already found. That's because we don't have a GC yet._

## 2 每次使用全局变量时，通过名称在哈希表中查找全局变量的速度都很慢，即使有一个很好的哈希表也是如此。你能想出一种更有效的方法，在不改变语义的情况下存储和访问全局变量吗？

**思路是将变量名预先映射成一个整数，intepreter 只需读取数组的对应下标元素**

There are a few ways to solve this. I'll do one that introduces another layer
of indirection, and a little information sharing between the compiler and VM.
有几种方法可以解决这个问题。我将采用一种引入另一层间接层的方法，并在编译器和虚拟机（VM）之间共享一些信息。

In the VM, we remove the `globals` hash table and replace it with:

```c
Table globalNames;
ValueArray globalValues;
```

**值数组是全局变量值存储的地方。哈希表将全局变量的名称映射到值数组中的索引。**

The value array is where the global variable values live. The hash table maps
the name of a global variable to its index in the value array. So, if the
program is:

```lox
var a = "value";
```

Then `globalNames` will contain a single entry, `"a" -> 0` and `globalValues`
will contain a single element, `"value"`. This association is all wired up at
compile time:

```c
static uint8_t identifierConstant(Token* name) {
  Value index;
  ObjString* identifier = copyString(name->start, name->length);
  if (tableGet(&vm.globalNames, identifier, &index)) {
    return (uint8_t)AS_NUMBER(index);
  }

  uint8_t newIndex = (uint8_t)vm.globalValues.count;
  writeValueArray(&vm.globalValues, UNDEFINED_VAL);

  tableSet(&vm.globalNames, identifier, NUMBER_VAL((double)newIndex));
  return newIndex;
}
```

When compiling a reference to a global variable, we see if we've ever
encountered its name before. If so, we know what index the value will be in in
the `globalValues` array. Otherwise, we add a new empty undefined value in the
array and then store a new hash table entry binding the name to that index.

Even though these two fields live in the VM, the compiler creates them at
compile time. You can think of it sort of like statically allocating memory for
the globals. We actually store the values in the VM so that they persist across
multiple REPL entries. We need to store the name association there too so that
we can find existing global variables.

`UNDEFINED_VAL` is a new, separate singleton value like `nil`. It's used to
mark a global variable slot as not having been defined yet. We can't use `nil`
because `nil` is a valid value to store in a variable.

At runtime, the instructions work like so:

```c
      case OP_GET_GLOBAL: {
        Value value = vm.globalValues.values[READ_BYTE()];
        if (IS_UNDEFINED(value)) {
          runtimeError("Undefined variable.");
          return INTERPRET_RUNTIME_ERROR;
        }
        push(value);
        break;
      }

      case OP_DEFINE_GLOBAL: {
        vm.globalValues.values[READ_BYTE()] = pop();
        break;
      }

      case OP_SET_GLOBAL: {
        uint8_t index = READ_BYTE();
        if (IS_UNDEFINED(vm.globalValues.values[index])) {
          runtimeError("Undefined variable.");
          return INTERPRET_RUNTIME_ERROR;
        }
        vm.globalValues.values[index] = peek(0);
        break;
      }
```

The operand for the instructions is now the direct index of the global variable
in the `globalValues` array. We've looked up the slot at compile time and
bound the result, so at runtime we don't need to worry about the name at all.
This is much faster. The only perf hit we take now is the necessary check at
runtime to ensure the variable has been initialized.
**这些指令的操作数现在是全局变量在 globalValues 数组中的直接索引**。我们在编译时已经进行了查找并绑定了结果，因此在运行时我们不需要再关心名称。这要快得多。现在我们唯一的性能损失是在运行时进行必要的检查，以确保变量已被初始化。

## 3 在 REPL 中运行时，用户可能会编写一个引用未知全局变量的函数。然后，在下一行中，他们声明了该变量。Lox 应该优雅地处理这种情况，在函数首次定义时不报告 "未知变量 "编译错误。

This question is more subtle than it may seem.

The seemingly safe error is to say that obviously using a variable that is
never defined anywhere is clearly wrong code so it should be an error. That's
a reasonable choice.

But when you're in the middle refactoring a large program, you sometimes have
code in a known broken state. As long as the broken code isn't _called_, it
might be nice to let the user run the other parts of the program that are OK.

You could try to have your cake and eat it too by making a reference to an
undeclared variable be a _warning_. That usually means the language reports it
as an error but still allows the program to be run. That works too, but in
practice, having shades of gray in your error reporting tends to cause user
headaches.

Some teams will want things to be black and white by turning all warnings into
errors, which sacrifices the ability you were trying to provide. Meanwhile,
other teams have the bad habit of committing code containing unfixed warnings,
leading to gradually worsening code. You will likely end up in long arguments
about which diagnostics should be considered fatal errors and which mere
warnings. People have strangely strong opinions about this stuff.

Personally, I'm pretty error-prone and like tools and languages to help me catch
my mistakes, so I'd like it to tell me if there's a use of an undeclared
variable name. If I'm in the middle of refactoring a big codebase, I'm OK with
having to comment out large regions of it to temporarily silence errors. But
that's just me.
