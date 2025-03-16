下面提供一份手写版代码，分别实现了 ES5 的 Object.create、new、extends（继承）以及 Object.assign 的简化版。每个方法都附有详细注释，便于记忆和面试书写。

---

### 1. 手写 Object.create

```js
function myObjectCreate(proto) {
  // 创建一个空的构造函数
  function F() {}
  // 将该构造函数的原型指向传入的对象
  F.prototype = proto
  // 返回新创建的对象，其内部 [[Prototype]] 指向 proto
  return new F()
}

// 使用示例：
const parent = { hello: 'world' }
const obj = myObjectCreate(parent)
console.log(obj.hello) // 输出：world
```

> **要点**：利用一个空函数，将其 prototype 指向目标对象，然后 new 出来的对象就继承了该原型。

---

### 2. 手写 new 操作符

```js
function myNew(constructor, ...args) {
  // 1. 创建一个空对象，并将其原型链接到 constructor.prototype
  const obj = Object.create(constructor.prototype)
  // 2. 执行构造函数，将 this 指向新对象，并传入参数
  const result = constructor.apply(obj, args)
  // 3. 如果构造函数返回一个对象，则返回它，否则返回新创建的对象
  return result !== null && (typeof result === 'object' || typeof result === 'function')
    ? result
    : obj
}

// 使用示例：
function Person(name) {
  this.name = name
}
Person.prototype.sayHello = function () {
  console.log('Hello, ' + this.name)
}

const person = myNew(Person, 'Alice')
person.sayHello() // 输出：Hello, Alice
```

> **要点**：手写 new 操作符时，要注意设置原型链、调用构造函数并处理构造函数返回值的情况。

---

### 3. 手写 extends（继承）

这里采用经典的原型链继承方式，模拟 ES5 中常用的继承方法。

```js
function inherits(Child, Parent) {
  // 1. 让子类原型继承父类原型，避免直接赋值造成引用问题
  Child.prototype = Object.create(Parent.prototype)
  // 2. 修正 constructor 指向子类构造函数
  Child.prototype.constructor = Child
}

// 使用示例：
function Animal(name) {
  this.name = name
}
Animal.prototype.speak = function () {
  console.log(this.name + ' makes a noise.')
}

function Dog(name, breed) {
  // 继承父类属性
  Animal.call(this, name)
  this.breed = breed
}
// 继承父类方法
inherits(Dog, Animal)

Dog.prototype.bark = function () {
  console.log(this.name + ' barks.')
}

const dog = new Dog('Buddy', 'Labrador')
dog.speak() // 输出：Buddy makes a noise.
dog.bark() // 输出：Buddy barks.
```

> **要点**：子类构造函数中用 Parent.call(this, ...) 继承父类属性，再用 Object.create 设置原型链，最后修正 constructor。

---

### 4. 手写 Object.assign

```js
function myAssign(target, ...sources) {
  if (target == null) {
    throw new TypeError('Cannot convert undefined or null to object')
  }
  // 强制转换为对象
  const to = Object(target)
  // 遍历所有源对象
  for (let source of sources) {
    if (source != null) {
      // 跳过 null 和 undefined
      // 只复制对象自身的可枚举属性
      for (let key in source) {
        if (Object.prototype.hasOwnProperty.call(source, key)) {
          to[key] = source[key]
        }
      }
    }
  }
  return to
}

// 使用示例：
const targetObj = { a: 1 }
const source1 = { b: 2 }
const source2 = { c: 3 }
const resultObj = myAssign(targetObj, source1, source2)
console.log(resultObj) // 输出：{ a: 1, b: 2, c: 3 }
```

> **要点**：遍历每个源对象的自有属性，复制到目标对象中；注意对 null/undefined 的判断和 TypeError 抛出。

---

### 总结

- **Object.create**：通过空构造函数实现原型链的链接。
- **new 操作符**：新建对象、绑定原型、调用构造函数、返回结果。
- **extends（继承）**：利用构造函数调用和原型链继承实现父子类关系。
- **Object.assign**：复制源对象的自有属性到目标对象中，注意细节判断。

这几段代码都是常见的面试题目，记住原理和代码结构，可以帮助你在面试中自信地写出答案。
