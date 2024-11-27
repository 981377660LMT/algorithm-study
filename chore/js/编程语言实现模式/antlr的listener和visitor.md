### ANTLR 中 Listener 和 Visitor 的区别

**Listener（监听器）** 和 **Visitor（访问者）** 是 ANTLR 提供的两种遍历语法树的方式，主要用于处理和操作解析后的语法树。它们的主要区别如下：

#### 1. 工作机制

- **Listener（监听器）**：

  - 基于事件驱动。
  - ANTLR 自动生成的监听器会在进入和退出每个语法规则时触发相应的方法，如 `enterRuleName` 和 `exitRuleName`。
  - 适用于无需返回值或只需执行操作的场景。

- **Visitor（访问者）**：
  - 基于递归遍历。
  - ANTLR 自动生成的访问者包含针对每个语法规则的 `visitRuleName` 方法，允许自定义遍历逻辑。
  - 适用于需要返回值或进行复杂计算的场景。

#### 2. 实现方式

- **Listener**：

  - 需要继承并实现 `ParseTreeListener` 接口。
  - 通过重写 `enter` 和 `exit` 方法来处理特定语法节点。
  - 更适合于简单的语法树遍历和操作。

- **Visitor**：
  - 需要继承并实现 `ParseTreeVisitor` 接口。
  - 通过重写 `visit` 方法来处理特定语法节点，并可以返回结果。
  - 更适合于需要从语法树中计算和收集信息的任务。

#### 3. 使用场景

- **Listener**：

  - 适用于代码生成、简单的语法检查等无需返回值的任务。
  - 由于是事件驱动，处理流程更为简单和直接。

- **Visitor**：
  - 适用于需要从语法树中提取信息、计算表达式值或构建复杂数据结构的任务。
  - 提供更大的灵活性，允许在遍历过程中返回和传递数据。

#### 4. 示例对比

假设有以下简单的表达式语法：

```antlr
expr: expr '+' term
    | term
    ;

term: INT;
```

- **Listener 使用示例**：

  ```java
  public class MyListener extends ExprBaseListener {
      @Override
      public void exitExpr(ExprParser.ExprContext ctx) {
          if (ctx.getChildCount() == 3) { // expr '+' term
              System.out.println("Addition operation");
          }
      }
  }
  ```

- **Visitor 使用示例**：

  ```java
  public class MyVisitor extends ExprBaseVisitor<Integer> {
      @Override
      public Integer visitExpr(ExprParser.ExprContext ctx) {
          if (ctx.getChildCount() == 3) { // expr '+' term
              int left = visit(ctx.expr(0));
              int right = visit(ctx.term());
              return left + right;
          } else {
              return visit(ctx.term());
          }
      }

      @Override
      public Integer visitTerm(ExprParser.TermContext ctx) {
          return Integer.parseInt(ctx.INT().getText());
      }
  }
  ```

#### 5. 总结

- **Listener** 更适合于简单的、事件驱动的语法树处理，不需要返回值的操作。
- **Visitor** 提供了更大的灵活性，适用于需要返回值或对语法树进行复杂操作的场景。

选择使用 Listener 还是 Visitor，取决于具体的应用需求和处理复杂度。
