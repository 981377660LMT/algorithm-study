zhangxiaokai 博客笔记
https://jasonkayzk.github.io/archive/

- Rust 从 panic 中恢复
  Rust 中可以通过在 Cargo.toml 中的 profile 中增加 panic 相应的配置来修改 panic 的行为；

  例如：

  默认情况下 panic 时，进程会打印当前出错的位置，然后退出；
  panic = “unwind”，允许抓取异常；
  panic = “abort”，出错 panic 时，直接 SigAbort 退出进程；

  ```toml
  [profile.dev]
  panic = "unwind"

  [profile.release]
  panic = "abort"
  ```

- Rust 中的默认初始化和初始化重载

```rust
#[derive(Debug, Default)]
pub struct Foo {
    bar: String,
    baz: i32,
    abc: bool,
}

fn main() {
    let x = Foo::default();

    let y = Foo { baz: 2, ..Default::default() };

    println!("{:?}", x);
    println!("{:?}", y);
}
```

- Rust 中创建全局变量
  全局变量的生命周期肯定是'static，即全局变量会一直存活至程序结束！
  但是不代表它需要用 static 来声明，例如常量、字符串字面值等无需使用 static 进行声明，原因是它们已经被打包到二进制可执行文件中
  创建全局变量主要包括下面几种情况：

  1. 编译期初始化
     const 创建静态常量，static 创建静态变量，Atomic 创建原子类型；
  2. 运行期初始化
     lazy_static 用于懒初始化，Box::leak 利用内存泄漏将一个变量的生命周期变为'static
     如果你使用的是较新的 Rust 版本，那么建议：
     优先直接使用 Rust 标准库中自带的 OnceCell 即可，而不再需要使用 lazy_static、once_cell 等 Crate！

- 在 Rust 中处理整数溢出
  默认情况下，当出现整型溢出时，Debug 模式会发生 panic，Release 模式下会在溢出后取舍归零； #[allow(arithmetic_overflow)] 来关闭该检查
  对于所有的有符号和无符号整数，Rust 提供了四组不同的运算函数，这提供了`显式`处理整数溢出的方式；
  1. `wrapping_` 系列函数处理整数溢出的方法是回绕，即从整数类型的最大值回绕到最小值
  2. `overflowing_`系列函数返回值会多一个 bool 以指明是否有溢出产生
  3. `checked_`系列函数返回值是 Option，当溢出时返回 None
  4. `saturating_`系列函数处理整数溢出的方法是 Clamp，即将溢出的值设为整数类型的最大值或最小值
