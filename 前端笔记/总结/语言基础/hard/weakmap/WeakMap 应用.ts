// 4.1 通过 WeakMap 缓存计算结果
// 4.2 在 WeakMap 中保留私有数据

// 实际上 JavaScript 的 WeakMap 并不是真正意义上的弱引用：其实只要键仍然存活，它就强引用其内容。WeakMap 仅在键被垃圾回收之后，才弱引用它的内容。为了提供真正的弱引用，TC39 提出了 WeakRefs 提案。
// WeakRef 是一个更高级的 API，它提供了真正的弱引用，并在对象的生命周期中插入了一个窗口。同时它也可以解决 WeakMap 仅支持 object 类型作为 Key 的问题。
