// 重构markDiff
//
// 肯定不能在子类拦截，必须在People类拦截。
// 1. People类加foo拦截，内部调用子类doFoo方法。成立的前提是原来子类每个上都有foo方法。
// 2. IBase接口加preFoo方法，外部调用时前置拦截(感觉代码不内聚...，且改动了框架)。
// 3. People类构造函数种改写this.foo方法(this总是指向当前对象，因此调用this.foo自然可以获取到此时调用的foo.)
// !类初始化顺序：先初始化子类，子类构造函数种super再调用父类.
// !第三种方法最好，关键在于：构造函数里的this.foo总是可以获取到当前对象最后被调用的是哪个方法(相当于手动查找原型链)。
//
// !总结两点技巧：
// 1. 父类构造函数中可以对子类方法进行统一拦截。(不仅仅是装饰器可以拦截)
// 2. this.method 在任意方法中可以获取到当前对象`真实调用`到的方法。
//
// 复盘：
// 为什么不能在People使用 foo 拦截 + doFoo 抽象类方法？
// 因为有的子类调用基类的foo，拦截后调用对象不明确.

export {}

interface IBase {
  foo(arg: any): void
}

abstract class Base implements IBase {
  foo(arg: any): void {
    console.log('Base foo', arg, this)
  }
}

abstract class People extends Base {
  // !答案在这里：
  constructor() {
    super()

    // !拦截所有子类的foo方法
    const preFoo = this.foo.bind(this)
    this.foo = (arg: any) => {
      if (typeof arg !== 'string') arg = String(arg)
      preFoo(arg)
    }
  }
}

// 有100个继承People的类，有的实现了foo方法有的没有，且所有重写的foo方法只接受字符串类型参数
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
  obj2.foo(1) // 要求：如果发现foo传入数字，需要转成字符串；请对上述代码重构，让代码不报错
}
