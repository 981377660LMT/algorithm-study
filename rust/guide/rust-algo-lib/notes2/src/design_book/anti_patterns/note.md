# Clone to satisfy the borrow checker

# #[deny(warnings)]

# Deref Polymorphism

不要滥用 deref 特性，仅为了调用方法方便而继承.
Deref 特性是被设计用来实现自定义指针类型的。它的用处是将 T 的引用转变为 T 的值，而不是在类型间转换。
