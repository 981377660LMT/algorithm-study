// !由于上面第一种写法不够面向对象，所以在面向对象语言里通常会使用 Interpreter 模式：

export {}

abstract class AstTree {}

class NumberLiteral extends AstTree {
  constructor(public readonly num: number) {
    super()
  }

  eval(): number {
    return this.num
  }
}

class BinaryExpr extends AstTree {
  constructor(
    public readonly lhs: NumberLiteral,
    public readonly op: string,
    public readonly rhs: NumberLiteral
  ) {
    super()
  }

  eval(): number {
    // eslint-disable-next-line no-eval
    return eval(`${this.lhs.eval()} ${this.op} ${this.rhs.eval()}`)
  }
}

if (require.main === module) {
  const expr = new BinaryExpr(new NumberLiteral(1), '+', new NumberLiteral(2))
  console.log(expr.eval())
}

// !跟第一种写法相比的话，相当于把独立的遍历方法中对不同类型节点的处理拆到了节点各自的实现中
// !然而当需要新增一种遍历方法(新增列)时，则需要对已有的节点类型都增加新方法的实现，所以代价就是修改N种节点类型。
