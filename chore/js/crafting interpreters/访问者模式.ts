// 我们希望能够定义新的糕点操作（烹饪，食用，装饰等），而不必每次都向每个类添加新方法
// 我们为每个类添加了一个accept（）方法，我们可以根据需要将其用于任意数量的访问者，而无需再次修改pastry类。 这是一个聪明的模式。
//
// !思路：不能在里面加，那就只能把自己丢给外面(visitor)，让外面的人来处理自己
// !对应在表格中，每个pastry类都是一行，但如果你看一个visitor的所有方法，它们就会形成一列。
// 在实践中，访问者通常希望定义能够产生值的操作。

interface PastryVisitor {
  visitBeignet(beignet: Beignet): void
  visitCruller(cruller: Cruller): void
}

/** 糕点. */
abstract class Pastry {
  abstract accept(visitor: PastryVisitor): void
}

/** 贝涅饼/甜甜圈. */
class Beignet extends Pastry {
  override accept(visitor: PastryVisitor): void {
    visitor.visitBeignet(this)
  }
}

/** 油炸麻花. */
class Cruller extends Pastry {
  override accept(visitor: PastryVisitor): void {
    visitor.visitCruller(this)
  }
}

export {}

if (typeof require !== 'undefined' && typeof module !== 'undefined' && require.main === module) {
  class PastryCooker implements PastryVisitor {
    visitBeignet(beignet: Beignet): void {
      console.log('烹饪贝涅饼...')
    }

    visitCruller(cruller: Cruller): void {
      console.log('烹饪油炸麻花...')
    }
  }

  class PastryEater implements PastryVisitor {
    visitBeignet(beignet: Beignet): void {
      console.log('吃贝涅饼...')
    }

    visitCruller(cruller: Cruller): void {
      console.log('吃油炸麻花...')
    }
  }

  class PastryDecorator implements PastryVisitor {
    visitBeignet(beignet: Beignet): void {
      console.log('装饰贝涅饼...')
    }

    visitCruller(cruller: Cruller): void {
      console.log('装饰油炸麻花...')
    }
  }

  const beignet = new Beignet()
  const cruller = new Cruller()

  const cooker = new PastryCooker()
  const eater = new PastryEater()
  const decorator = new PastryDecorator()

  beignet.accept(cooker)
  beignet.accept(eater)
  beignet.accept(decorator)

  cruller.accept(cooker)
  cruller.accept(eater)
  cruller.accept(decorator)
}
