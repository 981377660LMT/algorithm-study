1.  为什么函数名在定义时可以立即绑定并安全使用，而其他变量则必须在初始化后才能使用？
    `关键在于，函数名在函数体内是立即可用的。相比之下，变量在声明时通常不会立即绑定到一个值，必须在初始化之后才能安全使用。`

    Consider:

    ```lox
    fun foo() {
      if (itsTuesday) foo();
    }
    ```

    The function does call itself inside it's definition. But it relies on some initial outer call to kick off the **recursion**. Some outside code must refer to "foo" by name first. That can't happen until the function declaration statement itself has finished executing. By then, "foo" is fully defined and is safe to use.

2.  你所知道的其他语言是如何处理在初始化器中引用相同名称的局部变量的，例如

    ```lox
    var a = "outer";
    {
      var a = a;
    }
    ```

    ```c
    #include <stdio.h>

    int main() {
        int a = 10;
        {
            int a = a; // 内部的a引用自身，未初始化
            printf("%d\n", a); // 输出未定义的值
        }
        return 0;
    }
    ```

    ```java
    class Main {
        public static void main(String[] args) {
            String a = "outer";
            {
                String a = a;  // error: variable a might not have been initialized
            }
        }
    }
    ```

    ```py
    def main():
        a = "outer"
        def inner():
            a = a  # UnboundLocalError: local variable 'a' referenced before assignment
        inner()

    main()
    ```

    是运行时错误？编译错误？允许吗？他们对全局变量的处理方式不同吗？你同意他们的选择吗？请说明你的答案。

    - C：允许但危险，可能导致未定义行为。
    - Java：禁止，编译时错误，安全可靠。
    - Dart 和 Python：运行时错误，防止未定义行为。
    - 推荐：像 Java 一样，在编译时阻止变量在初始化器中引用自身，以提高代码安全性和可维护性。

    In C, the variable is put in scope before its initializer, which means that
    the initializer refers to the variable being initialized. Since C does not
    require any clearing of uninitialized memory, it means you could potentially
    access garbage data.

    Java does not allow one local variable to shadow another so it's an error
    because of that if the outer variable is also local. The outer variable
    could be a field on the surrounding class. In that case, like C, the local
    variable is in scope in its own initializer. However, Java makes it an error
    to refer to a variable that may not have been initialized, so this falls
    under that case and is an error.

    Obviously, C's approach is crazy talk. Java is fine and takes advantage of
    definite assignment analysis, which is useful for other things (like
    ensuring final fields are initialized before the constructor body
    completes). I like when languages get a lot of mileage out of a single
    concept.

3.  扩展解析器，使其在局部变量从未使用的情况下报错。

    `基本思想是，不再仅仅为每个局部变量存储一个布尔状态，而是在解析代码时，允许变量处于以下三种状态之一：`

    - 已声明(declared/声明)但尚未定义(defined/初始化)。
    - 已定义但尚未读取。
    - 已读取。

    The basic idea is that instead of storing just a boolean state for each
    local variable as we resolve the code, we'll allow a variable to be in one
    of three states:

    1. It has been declared but not yet defined.
    2. It has been defined but not yet read.
    3. It has been read.

    Any variable that goes out of scope when in the defined-but-not-yet-read
    state is an error. The annoying part is that we can't detect the error until
    the variable goes out of scope, but we want to report it on the line that
    the variable was declared. So we also need to keep track of the token from
    the variable declaration. We'll bundle that and the three-state enum into
    a little class inside the Resolver class:
    任何在处于“已定义但尚未读取”状态下离开作用域的变量都会报错。令人烦恼的是，`我们无法在变量离开作用域之前检测到错误，但我们希望在变量声明的那一行报告错误`。因此，我们还需要跟踪变量声明时的标记（token）。我们将把这个标记和三状态枚举封装到 Resolver 类中的一个小类中：

    ```java
      private static class Variable {
        final Token name;
        VariableState state;

        private Variable(Token name, VariableState state) {
          this.name = name;
          this.state = state;
        }
      }

      private enum VariableState {
        DECLARED,
        DEFINED,
        READ
      }
    ```

    Then we change the scope stack to use that instead of Boolean:

    ```java
      private final Stack<Map<String, Variable>> scopes = new Stack<>();
    ```

    When we resolve a local variable, we mark it used. However, we don't want
    to consider assigning to a local variable to be a "use". Writing to a
    variable that's never read is still pointless. So we change resolveLocal()
    to:

    ```java
      private void resolveLocal(Expr expr, Token name, boolean isRead) {
        for (int i = scopes.size() - 1; i >= 0; i--) {
          if (scopes.get(i).containsKey(name.lexeme)) {
            interpreter.resolve(expr, scopes.size() - 1 - i);

            // 标记为已使用
            if (isRead) {
              scopes.get(i).get(name.lexeme).state = VariableState.READ;
            }
            return;
          }
        }

        // 未找到。当成全局变量。
      }
    ```

    每次调用 resolveLocal() 都需要传入该标志。在 visitVariableExpr() 中，这个标志为 true：
    在 visitAssignExpr() 中，这个标志为 false：

    Every call to resolveLocal() needs to pass in that flag. In
    visitVariableExpr(), it's true:

    ```java
        resolveLocal(expr, expr.name, true);
    ```

    In visitAssignExpr(), it's false:

    ```java
        resolveLocal(expr, expr.name, false);
    ```

    Next, we update the existing code that touches scopes to use the new
    Variable class:
    接下来，我们更新处理作用域的现有代码，以使用新的 Variable 类：

    ```java
      public Void visitVariableExpr(Expr.Variable expr) {
        if (!scopes.isEmpty() &&
            scopes.peek().containsKey(expr.name.lexeme) &&
            scopes.peek().get(expr.name.lexeme).state == VariableState.DECLARED) {
          Lox.error(expr.name,
              "Can't read local variable in its own initializer.");
        }

        resolveLocal(expr, expr.name, true);
        return null;
      }

      private void beginScope() {
        scopes.push(new HashMap<String, Variable>());
      }

      // declare 方法在当前作用域中声明一个变量，并将其状态设置为 DECLARED。如果同一作用域中已经存在同名变量，则报错。
      private void declare(Token name) {
        if (scopes.isEmpty()) return;

        Map<String, Variable> scope = scopes.peek();
        if (scope.containsKey(name.lexeme)) {
          Lox.error(name,
              "Already variable with this name in this scope.");
        }

        scope.put(name.lexeme, new Variable(name, VariableState.DECLARED));
      }

      // define 方法将变量的状态从 DECLARED 更新为 DEFINED。
      private void define(Token name) {
        if (scopes.isEmpty()) return;
        scopes.peek().get(name.lexeme).state = VariableState.DEFINED;
      }
    ```

    Finally, when a scope is popped, we check its variables to see if any were
    not read:
    当一个作用域被弹出时，我们检查其变量是否有未被读取的：

    ```java
      private void endScope() {
        Map<String, Variable> scope = scopes.pop();

        for (Map.Entry<String, Variable> entry : scope.entrySet()) {
          if (entry.getValue().state == VariableState.DEFINED) {
            Lox.error(entry.getValue().name, "Local variable is not used.");
          }
        }
      }
    ```

4.  我们的Resolver会计算变量是在哪个环境中找到的，但它仍然是按名称在Map中查找的。更有效的环境表示方法是将局部变量存储在数组中，然后按索引查找。`扩展Resolver，为作用域中声明的每个局部变量关联一个唯一的索引。在解析变量访问时，查找变量所在的作用域及其索引并将其存储起来`。在解释器中，使用该索引而不是 map 来快速访问变量。

    ```java
    // resolver
    private class Variable {
        boolean isDefined = false;
        final int slot;

        private Variable(int slot) {
          this.slot = slot;
        }
    }

    private void declare(Token name) {
        if (scopes.isEmpty()) return;
        Map<String, Variable> scope = scopes.peek();
        if (scope.containsKey(name.lexeme)) {
          Lox.error(name,
              "Already variable with this name in this scope.");
        }
        scope.put(name.lexeme, new Variable(scope.size()));  // 为变量分配一个唯一的索引
    }

    @Override
    public Void visitVariableExpr(Expr.Variable expr) {
        if (!scopes.isEmpty() &&
            scopes.peek().containsKey(expr.name.lexeme) &&
            !scopes.peek().get(expr.name.lexeme).isDefined) {
          Lox.error(expr.name,
              "Can't read local variable in its own initializer.");
        }

        resolveLocal(expr, expr.name);
        return null;
    }

    @Override
    public Void visitAssignExpr(Expr.Assign expr) {
        resolve(expr.value);
        resolveLocal(expr, expr.name);
        return null;
    }

    private void resolveLocal(Expr expr, Token name) {
        for (int i = scopes.size() - 1; i >= 0; i--) {
          Map<String, Variable> scope = scopes.get(i);
          if (scope.containsKey(name.lexeme)) {
            interpreter.resolve(expr, scopes.size() - 1 - i,
                scope.get(name.lexeme).slot);
            return;
          }
        }

        // Not found. Assume it is global.
    }


    // interpreter
    void resolve(Expr expr, int depth, int slot) {
        locals.put(expr, depth);
        slots.put(expr, slot);
    }

    @Override
    public Object visitVariableExpr(Expr.Variable expr) {
        return lookUpVariable(expr.name, expr);
    }

    private Object lookUpVariable(Token name, Expr expr) {
        Integer distance = locals.get(expr);
        if (distance != null) {
          return environment.getAt(distance, slots.get(expr));
        } else {
          if (globals.containsKey(name.lexeme)) {
            return globals.get(name.lexeme);
          } else {
            throw new RuntimeError(name,
                "Undefined variable '" + name.lexeme + "'.");
          }
        }
    }

    @Override
    public Object visitAssignExpr(Expr.Assign expr) {
        Object value = evaluate(expr.value);
        Integer distance = locals.get(expr);
        if (distance != null) {
          environment.assignAt(distance, slots.get(expr), value);
        } else {
          if (globals.containsKey(expr.name.lexeme)) {
            globals.put(expr.name.lexeme, value);
          } else {
            throw new RuntimeError(expr.name,
                "Undefined variable '" + expr.name.lexeme + "'.");
          }
        }
        return value;
    }

    // environment
    class Environment {
      final Environment enclosing;
      private final List<Object> values = new ArrayList<>();

      // ...

      Object getAt(int distance, int slot) {
        Environment environment = this;
        for (int i = 0; i < distance; i++) {
          environment = environment.enclosing;
        }
        return environment.values.get(slot);
      }

      void assignAt(int distance, int slot, Object value) {
        Environment environment = this;
        for (int i = 0; i < distance; i++) {
          environment = environment.enclosing;
        }
        environment.values.set(slot, value);
      }

      // ...
    }
    ```
