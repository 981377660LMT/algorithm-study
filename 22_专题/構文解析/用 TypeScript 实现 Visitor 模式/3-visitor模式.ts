// !Visitor 模式又多加了一层调用，所以理解起来可能不是很直观：

export {}

interface AstVisitor<T> {
  visitNumberLiteral: (node: NumberLiteral) => T
  visitBinaryExpr: (node: BinaryExpr) => T
}

class EvalVisitor implements AstVisitor<number> {
  visitNumberLiteral(node: NumberLiteral): number {
    return node.num
  }

  visitBinaryExpr(node: BinaryExpr): number {
    // eslint-disable-next-line no-eval
    return eval(`${node.lhs.accept(this)} ${node.op} ${node.rhs.accept(this)}`)
  }
}

abstract class AstTree {}

class NumberLiteral extends AstTree {
  constructor(public readonly num: number) {
    super()
  }

  accept(visitor: AstVisitor<unknown>) {
    return visitor.visitNumberLiteral(this)
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

  accept(visitor: AstVisitor<unknown>) {
    return visitor.visitBinaryExpr(this)
  }
}

if (require.main === module) {
  const expr = new BinaryExpr(new NumberLiteral(1), '+', new NumberLiteral(2))
  console.log(expr.accept(new EvalVisitor()))
}

// !一个方法对应一种节点类型。从代码组织上来看，其实和第一种方法一样，又把对不同节点类型的处理逻辑放在了一起。
