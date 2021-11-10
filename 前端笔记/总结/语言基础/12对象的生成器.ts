export {}
const myObj = {
  *myGeneratorMethod() {
    yield 1
    yield 2
    yield 3
  },
}
const genObj1 = myObj.myGeneratorMethod()
for (const item of genObj1) {
}
/////////////////////////////////////////////////

class MyClass {
  *myGeneratorMethod() {
    yield 1
    yield 2
    yield 3
  }
}
const myObject = new MyClass()
const genObj2 = myObject.myGeneratorMethod()
for (const item of genObj2) {
}
/////////////////////////////////////////////////////
const SomeObj = {
  *[Symbol.iterator]() {
    yield 1
    yield 2
    yield 3
  },
}

console.log(Array.from(SomeObj)) // [ 1, 2, 3 ]
