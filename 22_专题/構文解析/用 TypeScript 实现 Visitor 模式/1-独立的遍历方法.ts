// https://qszhu.github.io/2021/10/31/visitor-revisited.html
// 在编译器的实现中经常会看到Visitor模式

// 假设现在有两种AST节点，数字NumberLiteral和二元表达式BinaryExpr：
// 当对表达式求值时，就需要遍历AST节点：
export {}

abstract class AstTree {}

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

function evaluate(tree: AstTree): number {
  if (tree instanceof NumberLiteral) {
    return tree.num
  }
  if (tree instanceof BinaryExpr) {
    // eslint-disable-next-line no-eval
    return eval(`${evaluate(tree.lhs)} ${tree.op} ${evaluate(tree.rhs)}`)
  }
  throw new Error('unknown ast node')
}

if (require.main === module) {
  const expr = new BinaryExpr(new NumberLiteral(1), '+', new NumberLiteral(2))
  console.log(evaluate(expr))
}

// 在一个独立的方法中完成对所有不同类型节点的遍历。
// 由于需要显式判断节点的类型，就不那么面向对象，在用面向对象语言的实现中很少见。
// 当需要新增节点类型(新增行)时，在所有已有的遍历方法中都需要增加对节点类型的处理，
// !所以代价就是修改M个遍历方法。
