# 告别旧时代：不再为了优化而优化的 ES6

# V8 release v4.0: ES6 classes and the removal of legacy features

这是 V8 历史上的一个大版本跨越（虽然版本号跳跃主要是为了对齐 Chrome）。

### 1. 语法糖背后的真相：Class

ES6 `class` 正式落地。
虽然语法看起像 Java，但 V8 下层的实现仍然是基于 **Prototype Chain** 的。不同之处在于，`class` 定义的方法默认是 **Non-enumerable** 的。V8 在解析类定义时，会一次性生成所有的构造函数原型链路，减少了像过去那样在 `prototype` 上逐个挂载方法时产生的多次 Hidden Class 过渡。

### 2. 移除原生实现的观察者模式（Object.observe）

曾被寄予厚望的 `Object.observe` 在这个时期被宣布放弃（转而推向 Proxy）。这是 V8 团队的一个清醒决策：与其在内核层维护复杂的追踪逻辑，不如提供一套更通用、可拦截一切操作的 `Proxy` 模型。

### 3. 一针见血的见解

v4.0 体现了 V8 的“平衡之道”。它在拥抱 ES6 现代语法的同时，果断舍弃那些可能增加内核复杂度却收益有限的非标特性。这种对技术债务的警觉，是保持 V8 在移动端时代依然轻量的核心逻辑。
