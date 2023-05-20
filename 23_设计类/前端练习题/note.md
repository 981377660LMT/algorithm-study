fft 加速预测

https://leetcode.cn/circle/discuss/TbWS5j/
https://github.com/cheran-senthil/TLE/pull/40
https://github.com/cheran-senthil/TLE/blob/master/tle/util/ranklist/rating_calculator.py
https://github.com/XYShaoKang/refined-leetcode/blob/master/src/content/utils/predict.ts

---

## proxy

以下是一些常见的 Proxy 方法和技巧：

- 拦截器方法（Handler methods）：Proxy 对象通过使用一组特定的拦截器方法来拦截目标对象的操作，例如 get、set、apply、construct 等。这些方法可以在拦截时执行自定义的操作。
- 属性代理（Property proxy）：通过在 get 和 set 拦截器方法中操作目标对象的属性，可以实现对属性的访问和修改的自定义逻辑。这使得你可以对属性的读取和写入进行拦截和控制。
- 函数代理（Function proxy）：通过在 apply 拦截器方法中拦截函数的调用，你可以对`函数的调用`进行自定义操作，例如参数验证、性能监控、缓存等。
- 构造函数代理（Constructor proxy）：通过在 construct 拦截器方法中拦截对构造函数的 new 操作，你可以`自定义构造函数的行为`，例如创建单例、限制实例化次数等。
- 代理嵌套（Nested proxy）：`Proxy 对象可以嵌套/递归使用`，即在目标对象的代理上再创建一个代理。这允许你对更复杂的操作进行拦截和处理。
- 验证和过滤（Validation and filtering）：通过在拦截器方法中验证和过滤操作，可以实现对目标对象的访问和修改的限制。例如，你可以验证属性值的类型或长度，并阻止对无效数据的修改。
- 观察者模式（Observer pattern）：通过在拦截器方法中添加触发逻辑，可以实现对目标对象的观察。你可以在对象属性发生变化时触发事件或执行其他操作。
  这些只是 Proxy 的一些常见方法和技巧，Proxy 提供了强大的能力来自定义和增强 JavaScript 对象的行为。具体的用法和技巧取决于你的需求和创造力。

### 拦截器方法（Handler methods）

拦截器方法（Handler methods）是 Proxy 对象提供的一组特殊方法，用于拦截并自定义对目标对象的操作。通过在这些方法中添加自定义逻辑，可以实现对目标对象操作的拦截、修改和控制。

以下是一些常见的拦截器方法：

- get(target, property, receiver)：拦截对目标对象的属性读取操作。当通过 target[property] 或 target.property 读取属性时，此方法将被调用。你可以在该方法中返回自定义的属性值或执行其他操作。

- set(target, property, value, receiver)：拦截对目标对象的属性设置操作。当通过 target[property] = value 或 target.property = value 设置属性值时，此方法将被调用。你可以在该方法中进行自定义的属性设置逻辑。

- apply(target, thisArg, argumentsList)：拦截对目标对象的函数调用操作。当通过 target(...arguments) 调用目标对象的函数时，此方法将被调用。你可以在该方法中执行自定义的函数调用逻辑。

- construct(target, argumentsList, newTarget)：拦截对目标对象的构造函数调用操作。当通过 new target(...arguments) 创建目标对象的实例时，此方法将被调用。你可以在该方法中自定义实例的创建逻辑。

- has(target, property)：拦截对目标对象的 in 操作符。当通过 property in target 判断属性是否存在时，此方法将被调用。你可以在该方法中返回自定义的判断结果。

- deleteProperty(target, property)：拦截对目标对象的属性删除操作。当通过 delete target[property] 删除属性时，此方法将被调用。你可以在该方法中执行自定义的属性删除逻辑。

通过重写这些拦截器方法，你可以在对目标对象进行读取、设置、调用等操作时插入自定义的行为。这样可以实现属性访问的限制、函数调用的验证、属性的动态计算等功能。注意，拦截器方法应返回相应的值或执行相应的操作，以确保正确的行为。

### 属性代理（Property proxy）

属性代理（Property proxy）是指使用 Proxy 对象拦截并自定义目标对象的属性访问操作。通过在拦截器方法中定义自定义逻辑，可以对属性的读取、设置、删除等操作进行拦截和修改。

属性代理可以用于实现各种功能和行为，如属性访问的限制、属性值的动态计算、属性的缓存和延迟加载等。下面介绍一些常见的属性代理的应用场景和技巧：

- 属性访问控制：通过拦截 get 和 set 方法，可以限制对某些属性的访问和修改。例如，可以在 get 方法中判断当前用户的权限，只允许有权限的用户访问敏感属性。

- 属性值计算：通过拦截 get 方法，可以在访问属性时动态计算属性的值。这对于一些需要复杂计算或从其他属性派生的属性特别有用。例如，可以在 get 方法中根据其他属性的值动态计算某个属性的结果。

- 属性缓存和延迟加载：通过拦截 get 方法，可以在首次访问属性时进行计算并将结果缓存起来，以便后续的访问可以直接返回缓存的结果。这样可以减少重复计算的开销，提高性能。类似地，可以通过拦截 get 方法实现属性的延迟加载，只在需要时才真正获取属性的值。

- 属性验证和转换：通过拦截 set 方法，可以对属性进行验证和转换。可以在 set 方法中检查传入的值是否符合预期的类型、范围或其他条件，并进行相应的处理。这样可以确保属性值的有效性和一致性。

- 属性监听和触发事件：通过拦截 get 和 set 方法，可以在属性被访问或修改时触发相应的事件或回调函数。这样可以实现属性的监听和响应机制，用于处理属性变化时的逻辑。

属性代理通过拦截器方法提供了对属性访问操作的灵活控制，使我们能够根据需求自定义属性的行为。但需要注意的是，在拦截器方法中确保返回正确的值或执行相应的操作，以确保代理行为的正确性和一致性。

### 函数代理（Function proxy）

函数代理可以用于实现各种功能和行为，如函数调用的拦截、参数验证和转换、函数执行的记录和追踪等。下面介绍一些常见的函数代理的应用场景和技巧：

- 函数调用拦截：通过拦截 apply 或 call 方法，可以在函数调用之前和之后执行自定义的逻辑。可以用于记录函数调用的次数、执行时间等统计信息，或在函数调用前后进行一些准备或清理工作。

- 参数验证和转换：通过拦截 apply 或 call 方法，可以对函数的参数进行验证和转换。可以在函数调用之前检查参数的类型、范围或其他条件，并根据需要进行相应的处理。这样可以确保函数参数的有效性和一致性。

- 函数执行记录和追踪：通过拦截 apply 或 call 方法，可以记录函数的执行过程和结果，用于调试和追踪函数的执行流程。可以在函数调用前后记录日志、打印输出或触发相应的事件，以便分析和排查问题。

- 函数返回值修改：通过拦截 apply 或 call 方法，可以修改函数的返回值。可以根据实际需求对函数的返回值进行修改、转换或处理。例如，可以对函数的返回结果进行缓存、修饰或加工，以满足特定的业务需求。

- 函数的动态调用：通过拦截 get 方法，可以实现函数的动态调用。可以通过属性访问的方式获取函数对象，并对其进行调用操作。这对于实现一些动态函数调用的场景特别有用，例如根据条件选择不同的函数进行调用。

函数代理通过拦截器方法提供了对函数调用操作的灵活控制，使我们能够根据需求自定义函数的行为。然而，需要注意在拦截器方法中确保返回正确的结果、参数和执行上下文，以确保代理行为的正确性和一致性。

### 构造函数代理（Constructor proxy）

构造函数代理可以用于实现各种功能和行为，如构造函数的参数验证和转换、实例化过程的拦截和修改、实例化的缓存和单例模式等。下面介绍一些常见的构造函数代理的应用场景和技巧：

- 参数验证和转换：通过拦截 construct 方法，可以对构造函数的参数进行验证和转换。可以在实例化过程中检查参数的类型、范围或其他条件，并根据需要进行相应的处理。这样可以确保实例的参数的有效性和一致性。

- 实例化过程的拦截和修改：通过拦截 construct 方法，可以拦截并修改构造函数的实例化过程。可以在实例化之前和之后执行自定义的逻辑，例如在实例化之前执行一些准备工作，或在实例化之后对实例进行一些处理。

- 实例化的**缓存和单例模式**：通过拦截 construct 方法，可以实现实例化的缓存和单例模式。可以在拦截器方法中维护一个缓存，检查是否已经实例化过相同参数的对象，如果是则直接返回缓存的实例，避免重复实例化。这对于需要复用相同参数的实例或实现单例模式非常有用。

- 实例属性和方法的动态添加：通过拦截 get 方法，可以动态添加实例的属性和方法。可以在实例化之后，通过属性访问的方式动态添加实例的属性和方法，以扩展实例的功能。这对于需要动态给实例添加额外功能或行为的场景特别有用。

构造函数代理通过拦截器方法提供了对构造函数实例化过程的灵活控制，使我们能够根据需求自定义构造函数的行为。然而，需要注意在拦截器方法中确保返回正确的实例对象，并保持构造函数原有的行为和约束，以确保代理行为的正确性和一致性。

### 代理嵌套(Nested proxy)

递归地使用 proxy
https://leetcode.cn/problems/make-object-immutable/

代理嵌套（Nested proxy）是指在 Proxy 对象的拦截器方法中创建另一个 Proxy 对象，形成多层代理的结构。通过代理嵌套，我们可以对对象的多个层级进行拦截和自定义操作，实现更复杂的行为和功能。

代理嵌套可以应用于多种场景，以下是一些常见的应用示例和技巧：

- 深层属性拦截：通过代理嵌套，可以对对象的深层属性进行拦截和操作。例如，可以创建一个代理对象来拦截对象的顶层属性访问，并在拦截器方法中创建嵌套的代理对象来拦截深层属性的访问。这样可以实现对整个对象的属性访问进行拦截和修改，而不仅限于顶层属性。

- 联级操作：通过代理嵌套，可以在对象的各个层级上进行联级操作。例如，在代理的 get 方法中创建嵌套的代理对象，用于拦截和操作属性的访问。这样可以实现对对象的属性访问进行级联操作，例如自动转换、计算衍生属性等。

- 多层校验和过滤：通过代理嵌套，可以在多个层级上进行校验和过滤。例如，在代理的 set 方法中创建嵌套的代理对象，用于对属性值进行校验和过滤。这样可以实现对属性值的多层级校验和过滤，确保属性值的有效性和一致性。

- 深层方法拦截：通过代理嵌套，可以对对象的深层方法进行拦截和操作。例如，在代理的 get 方法中创建嵌套的代理对象，用于拦截和操作方法的调用。这样可以实现对对象的方法进行深层级的拦截和修改，以增强方法的行为和功能。

### 观察者模式(Observer pattern)

JavaScript 中的 Proxy 对象可以用来创建一个观察者模式。使用 Proxy，我们可以创建一个拦截器在某个对象上进行各种操作。结合 Observer 模式，我们可以创建一个 Proxy，当特定属性发生改变时，可以自动通知所有的 Observer。

```js
class Observable {
  constructor(target) {
    this.listeners = []
    this.target = target
    this.proxy = new Proxy(target, {
      set: (target, property, value) => {
        target[property] = value
        this.notify(property, value)
        return true
      }
    })
  }

  subscribe(listener) {
    this.listeners.push(listener)
  }

  notify(property, value) {
    for (let listener of this.listeners) {
      listener.update(property, value)
    }
  }

  getProxy() {
    return this.proxy
  }
}

class Observer {
  update(property, value) {
    console.log(`Property ${property} has been updated with ${value}`)
  }
}

const obj = { value: 1 }
const observable = new Observable(obj)
const observer = new Observer()

observable.subscribe(observer)

const proxy = observable.getProxy()
proxy.value = 2
```

在这个例子中，我们有一个 Observable 类，它接收一个目标对象并创建一个 Proxy 对象来包装它。当 Proxy 对象的属性被修改时，set 拦截器会被触发，并通过调用 notify 方法来通知所有的监听器（观察者）。

我们还有一个 Observer 类，它定义了 update 方法，这个方法在属性变化时被 Observable 对象调用。

最后，我们创建一个 Observable 对象和一个 Observer 对象，并通过 subscribe 方法订阅了 Observable 对象。然后，我们通过修改 Proxy 对象的 value 属性，这导致所有的 Observer 对象得到通知。

### proxy 里的 receiver 参数

在 JavaScript 中，Proxy 的 handler 对象的方法经常会包含三个参数：target，prop，和 receiver。其中：

- target 是被代理的`原始对象`。
- prop 是要获取或设置的属性名。
- receiver `最初是被调用的对象。通常是代理本身`，但 handler 的方法也可以在不同的上下文中被调用，如 Reflect 或 Object.prototype 方法。

**receiver 在处理继承的时候会特别有用**。如果代理对象作为原型链上的一个对象，那么读取操作（get）可能需要从代理对象的原型链上获取值。在这种情况下，receiver 参数将是初始被调用的对象，这有助于保持原型链的正确性。

```js
const handler = {
  get(target, prop, receiver) {
    console.log(`Get was called with receiver: ${receiver === proxy}`)
    return Reflect.get(target, prop, receiver)
  }
}

const obj = { value: 42 }
const proxy = new Proxy(obj, handler)

console.log(proxy.value) // 输出：Get was called with receiver: true

const obj2 = Object.create(proxy)
console.log(obj2.value) // 输出：Get was called with receiver: false
```

在这个例子中，当通过 proxy 访问 value 属性时，receiver 是 proxy 本身。然后我们创建了 obj2，它的原型是 proxy。当我们尝试访问 `obj2.value 时，JavaScript 运行时将会在 obj2 的原型链上查找 value 属性，这会触发 proxy 的 get handler`。在这种情况下，receiver 是 obj2，而不是 proxy，因为 obj2 是最初尝试访问 value 的对象。

有点类似 event.target 和 event.currentTarget 的关系。
event.target: 是触发事件的元素，也就是事件最初发生的地方。如果你在一个按钮上点击，event.target 就会是那个按钮。
event.currentTarget: 是绑定事件处理函数的元素，即当前正在处理该事件的元素。在事件的捕获或冒泡阶段，event.currentTarget 可能会发生改变，因为事件可能会沿着 DOM 树向上或向下传播。

```js
document.querySelector('#parent').addEventListener('click', function (event) {
  console.log(event.target) // 输出：如果你点击了子元素，这将会输出子元素
  console.log(event.currentTarget) // 输出：无论你点击了父元素还是子元素，这都会输出父元素
})

document.querySelector('#child').addEventListener('click', function (event) {
  console.log(event.target) // 输出：无论你点击了哪个元素，这都会输出子元素
  console.log(event.currentTarget) // 输出：无论你点击了哪个元素，这都会输出子元素
})
```

**这里 proxy 的 receiver 就相当于 event.target，而 proxy 的 target 就相当于 event.currentTarget。**
正好是反过来的!

---

在 JavaScript 中，receiver 参数的主要用途是在处理 getter 和 setter 函数以及实现继承时保持正确的上下文。然而，你也可以使用 receiver 参数来实现一些更复杂的行为和设计模式，如跟踪属性的访问，创建一个链式 API，或实现一些类型的数据绑定。

在某些情况下，receiver 可以帮助你确定是否通过代理直接访问属性，或者是否通过继承的对象访问属性。这在你需要区分这两种情况时非常有用，例如，当你想要跟踪或控制对某些属性的直接访问时。

另外，**如果你的 getter 和 setter 函数依赖于它们的调用上下文（即 this 值），那么你可能需要使用 receiver 参数来保证它们的行为**。例如，如果你有一个计算属性，它的值是基于其他属性计算出来的，那么你的 getter 函数可能需要知道 this 是什么，以便正确地获取其他属性的值。在这种情况下，你可以使用 receiver 参数来调用 getter 或 setter 函数，确保 this 值是正确的。

这里有一个例子来展示如何使用 receiver 参数：

```js
let obj = {
  baseValue: 1,
  derivedValue: 0
}

let handler = {
  get(target, prop, receiver) {
    if (prop === 'derivedValue') {
      return receiver.baseValue * 2
    }
    return Reflect.get(target, prop, receiver)
  }
}

let proxy = new Proxy(obj, handler)

console.log(proxy.derivedValue) // 输出：2

let child = Object.create(proxy)
child.baseValue = 2

console.log(child.derivedValue) // 输出：4
```

在这个例子中，derivedValue 属性的 getter 函数会返回 baseValue 的两倍。因为我们使用 receiver 参数来获取 baseValue 的值，所以即使我们通过一个继承的对象访问 derivedValue，也可以得到正确的结果。

### Reflect 有什么用，为什么通常和 proxy 结合使用

并不是不用 proxy 就不行，而是 proxy 具有某些优点
https://stackoverflow.com/questions/25421903/what-does-the-reflect-object-do-in-javascript

- More useful return values 更有用的返回值

```js
try {
  Object.defineProperty(obj, name, desc)
  // property defined successfully
} catch (e) {
  // possible failure (and might accidentally catch the wrong exception)
}

if (Reflect.defineProperty(obj, name, desc)) {
  // success
} else {
  // failure
}
```

- First-class operations 函数式操作
  在 ES5 中，检测对象 obj 是否定义或继承某个属性名称的方法是写入 (name in obj) 。同样，要删除属性，请使用 delete obj[name] 。虽然专用语法既好又短，但这也意味着当您想要将这些操作作为一等值传递时，必须将这些操作显式包装在函数中。

  使用 Reflect 时，这些操作很容易被定义为一等函数：

  Reflect.has(obj, name) 是 (name in obj) 的功能等价物， Reflect.deleteProperty(obj, name) 是与 delete obj[name]. 相同的函数

- Control the this-binding of accessors 控制访问器的 this 绑定
  Reflect.get 和 Reflect.set 方法允许你执行相同的操作，但另外接受 receiver 参数作为最后一个可选参数，该参数允许你在获取/设置的属性是访问器时显式设置 this 绑定：

  ```js
  const name = ... // get property name as a string
  Reflect.get(obj, name, wrapper) // if obj[name] is an accessor, it gets run with `this === wrapper`
  Reflect.set(obj, name, value, wrapper)
  ```

  当您包装 obj 并且您希望访问器中的任何自发送重新路由到包装器时，这偶尔很有用，例如，如果 obj 定义为：

  ```js
  const obj = {
    get foo() { return this.bar(); },
    bar: function() { ... }
  }
  ```

  `调用 Reflect.get(obj, "foo", wrapper) 可以将 this.bar() 调用重新路由到 wrapper。`

### 注意 TypeError: Cannot create proxy with a non-object as target or handler

在使用 proxy 前必须保证 target 是一个 object
