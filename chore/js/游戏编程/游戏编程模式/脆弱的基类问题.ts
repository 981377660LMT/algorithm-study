// 脆弱的基类问题
// https://en.wikipedia.org/wiki/Fragile_base_class
// 对基类看似安全的修改如何通过进入无限递归而导致继承子类发生故障
// !超类最好避免更改对动态绑定方法的调用。

class A {
  private _counter = 0

  inc1() {
    this._counter++
    // this.inc2()
  }

  inc2() {
    this._counter++
  }
}

class B extends A {
  override inc2() {
    super.inc1()
  }
}

const b = new B()
b.inc1()

export {}
