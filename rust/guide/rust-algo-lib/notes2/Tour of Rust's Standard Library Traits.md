<!-- Tour of Rust's Standard Library Traits -->

1. 衍生宏 Derive Macros
   标准库导出了一系列实用的衍生宏，我们可以利用它们方便快捷地为特定类型实现某种特性，前提是该类型的成员亦实现了相应的特性。衍生宏与它们各自所实现的特性同名：

   Clone
   Copy
   Debug
   Default
   Eq
   Hash
   Ord
   PartialEq
   PartialOrd

2. Generic Blanket Impls 泛型覆盖实现
   泛型覆盖实现是一种在`泛型类型`而不是具体类型上的实现

```rust
trait Even {
    fn is_even(self) -> bool;
}

impl Even for i8 {
    fn is_even(self) -> bool {
        self % 2_i8 == 0_i8
    }
}

impl Even for u8 {
    fn is_even(self) -> bool {
        self % 2_u8 == 0_u8
    }
}

// generic blanket impl
// 一揽子泛型实现
impl<T> Even for T
where
    T: Rem<Output = T> + PartialEq<T> + Sized,
    u8: TryInto<T>,
    <u8 as TryInto<T>>::Error: Debug,
{
    fn is_even(self) -> bool {
        // these unwraps will never panic
        // 以下 unwrap 永远不会 panic
        self % 2.try_into().unwrap() == 0.try_into().unwrap()
    }
}
```

3. 特性对象 Trait Objects
   特性对象就给了我们运行时的多态性，可以允许函数在运行时动态地返回不同的类型

```rust
fn example(condition: bool, vec: Vec<i32>) -> Box<dyn Iterator<Item = i32>> {
    let iter = vec.into_iter();
    if condition {
        // Has type:
        // Box<Map<IntoIter<i32>, Fn(i32) -> i32>>
        // But is cast to:
        // Box<dyn Iterator<Item = i32>>
        Box::new(iter.map(|n| n * 2))
    } else {
        // Has type:
        // Box<Filter<IntoIter<i32>, Fn(&i32) -> bool>>
        // But is cast to:
        // Box<dyn Iterator<Item = i32>>
        Box::new(iter.filter(|&n| n >= 2))
    }
}
        // 以上代码中，两种不同的指针类型转换成相同的指针类型
```
