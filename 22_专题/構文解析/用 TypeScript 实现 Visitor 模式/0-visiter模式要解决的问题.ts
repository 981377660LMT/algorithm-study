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

// 类型\方法     |  evaluate()
// -----------   |----------------
// NumberLiteral |  tree.num
// BinaryExpr    |  eval(`${evaluate(tree.lhs)} ${tree.op} ${evaluate(tree.rhs)}`)
// !未来对于这张表格的扩展有两种方式，一种是增加新的行，也就是增加新的AST节点：
// StringLiteral |  tree.str
// 另一种是增加新的列，也就是增加新的方法：
// 显然这时对于已存在的行，都需要增加对新列的处理。当然还有可能是需要同时增加新的行和新的列：
// !各种不同的写法，包括 Visitor 模式想要解决的问题，都是希望以最小的改动已有代码的代价，完成对这张表格的扩展。
