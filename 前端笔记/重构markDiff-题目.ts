interface IBase {
  foo(arg: any): void
}

abstract class Base implements IBase {
  foo(arg: any): void {
    console.log('Base foo', arg, this)
  }
}

abstract class People extends Base {}

// !有100个继承People的类，有的实现了foo方法有的没有，且所有重写的foo方法只接受字符串类型参数
class Boy extends People {
  override foo(arg: string): void {
    if (typeof arg !== 'string') {
      throw new Error(`arg must be string, but got ${arg}`)
    }
    console.log('Boy foo', arg, this)
  }
}

class Girl extends People {}

if (require.main === module) {
  const obj1: IBase = new Girl()
  obj1.foo('1')
  obj1.foo(1)
  const obj2: IBase = new Boy()
  obj2.foo('1')
  obj2.foo(1) // 要求：如果发现foo传入数字，需要转成字符串；请对上述代码重构，让代码运行时不报错
}
