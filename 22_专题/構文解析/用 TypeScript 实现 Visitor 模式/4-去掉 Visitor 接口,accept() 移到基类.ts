// !Visitor 模式又多加了一层调用，所以理解起来可能不是很直观：

export {}

abstract class AstTree {
  accept(v: any) {
    return v[`visit${this.constructor.name}`](this)
  }
}

class NumberLiteral extends AstTree {
  constructor(public readonly num: number) {
    super()
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
}

class EvalVisitor {
  visitNumberLiteral(node: NumberLiteral): number {
    return node.num
  }

  visitBinaryExpr(node: BinaryExpr): number {
    // eslint-disable-next-line no-eval
    return eval(`${node.lhs.accept(this)} ${node.op} ${node.rhs.accept(this)}`)
  }
}

if (require.main === module) {
  const expr = new BinaryExpr(new NumberLiteral(1), '+', new NumberLiteral(2))
  console.log(expr.accept(new EvalVisitor()))
}

// !少写些代码，对修改的代价没有影响。
