ES5 和 ES6 子类 this 生成顺序不同。
**es5 是先创建子类实例对象的 this，然后将父类方法赋到这个 this 上。**

```JS

const obj = Object.create(constructor.prototype)
constructor.apply(obj,...args)
```

**es6 是先在子类构造函数中用 super 创建父类实例的 this,再在构造函数中进行修改它。**
因为 this 生成顺序不同，所以需要在 constructor 中，需要使用 super()
ES6 在继承的语法上不仅继承了类的原型对象，**还继承了类的静态属性和静态方法**

也因此，es5 中 array，error 等原生构造函数无法继承而 es6 就可以自己定义这些原生构造函数。
（es5 中子类无法拿到父类的内部属性，就算是 apply 也不行，es5 默认忽略 apply 传入的 this）。
