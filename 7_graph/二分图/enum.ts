// 如果没有用到，编译时被删除
const enum Foo {
  A = 1,
  B = 3,
}

// 生成对象
enum Bar {
  A = 1,
  B = 3,
}

const FooA = Foo.A
