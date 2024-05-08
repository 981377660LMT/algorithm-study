- Stack vs Heap

  - Stack 上的数据必须固定大小，未知则必须放在 Heap 上
  - 所有权解决的问题：

    - 跟踪代码的哪些部分正在使用 heap 的哪些数据
    - 最小化堆上的重复数据量
    - 清理堆上不再使用的数据以避免内存泄漏

    管理 heap 数据是所有权存在的原因

---

- 所有权规则

  - 每个值有且仅有一个变量，这个变量是其所有者
  - 当所有者超出作用域(scope)时，该值被删除

---

- String 类型
  三个部分：
  `ptr`：指向存储在其他地方的 UTF-8 编码的字符串数据
  `len`：长度
  `capacity`：容量
- 变量离开作用域时，Rust 会调用 `drop` 方法来清理堆上的内存

  ```rust
  {
      let s = String::from("hello");
  } // s 离开作用域，调用 drop 方法
  ```

- Move (移交)

  - Rust 不会自动创建引用数据类型的深拷贝
  - 当一个引用数据类型变量赋值给另一个变量时，原变量将失效，避免了 double free 问题
  - 遵循了原则：不能存在两个指向同一内存的指针，一个值只能有一个所有者

    ```rust
    let s1 = String::from("hello");
    let s2 = s1; // s1 失效
    ```

    - `s1` 被移动到 `s2`，`s1` 无效
    - `s2` 会在离开作用域时调用 `drop` 方法

- Clone (克隆)

  - 使用 `clone` 方法创建深拷贝

    ```rust
    let s1 = String::from("hello");
    let s2 = s1.clone();
    ```

    - `s1` 和 `s2` 都有效
    - `s1` 和 `s2` 都会在离开作用域时调用 `drop` 方法

- Stack 上的数据：复制
  - Copy trait：如果一个类型实现了 Copy trait，那么它的值可以在赋值后继续使用，不会失效
  - Drop trait：如果一个类型实现了 Drop trait，那么 Rust 不允许让他再去实现 Copy trait

---

- 所有权与函数
  函数传参与变量赋值类似
- 返回值与作用域
  - 返回值的所有权会被移动

---

一个变量的所有权总是遵循同样的模式：

- 把一个值赋给其他变量时就会发生移动
- 当一个包含 heap 数据的变量离开作用域时，它的值将被丢弃，除非数据被移动为另一个变量的所有权

---

如何让函数使用一个值但不获取其所有权？
shared ptr?
**引用**: `&` 符号

- 不可变引用：`&T`
- 可变引用：`&mut T`
  限制 1：**注意在同一作用域内只能有一个可变引用**

  ```rust
  let mut s = String::from("hello");
  let r1 = &mut s;
  let r2 = &mut s; // error
  ```

  好处：编译时避免数据竞争
  数据竞争的产生例子：

  - 两个或多个指针同时访问同一数据
  - 至少有一个指针在写入数据
  - 没有同步数据访问的机制

  限制 2：**不可同时拥有一个可变引用和一个不变引用**

- 悬垂引用：指向已经释放的内存
  Rust 编译器会在编译时检查悬空引用

总结：
在任意时刻，只能满足下列条件之一：

1. 一个可变引用
2. 任意数量的不可变引用
3. 引用必须一直有效

---

字符串切片：`&str`
切片不持有所有权，它只是引用了数据的一部分
&s[0..2]：从索引 0 开始，到 2 之前的位置（不包括 2）
&s[0..]
&s[..2]
&s[..]

字符串字面值是切片

---

struct

- 如果 struct 的实例是可变的，那么实例所有的字段都是可变的
- struct 更新语法
  ```rust
  let user1 = User {
    name: String::from("name"),
    email: String::from("email"),
    ..user1
  };
  ```
- tuple struct

  ```rust
  struct Point(i32, i32, i32);
  ```

- 空结构体(unit-like struct)

  ```rust
  struct Empty;
  ```

- struct 数据的所有权

---

打印 struct

- std::fmt::Display trait
  `println!("{}", instance);`
- std::fmt::Debug trait
  `println!("{:?}", instance);`
- #[derive(Debug)] 宏
  派生 Debug trait
- {:?} 和 {:#?}
  `println!("{:?}", instance);` -> 表示输出的是一个 Debug 格式的字符串
  `println!("{:#?}", instance);` -> 表示输出的是一个更易读的 Debug 格式的字符串

---

struct 的方法

- 类似 golang，结构体定义数据结构，impl 定义方法
- 调用方法时，**自动引用或解引用**，类似 golang

  ```rust
  p1.distance(&p2);
  (&p1).distance(&p2);
  (&&&&&&&&p1).distance(&p2);
  ```

关联函数

- 可以在 impl 块中定义不以 `self` 作为参数的函数
- 通常用于构造器

  ```rust
  String::from("hello")

  ```

- 每个 struct 都可以有多个 impl 块

---

枚举

```rust
enum IpAddrKind {
    V4,
    V6,
}

let four = IpAddrKind::V4;
let six = IpAddrKind::V6;
```

- 数据附加到枚举的每个变量上

  ```rust
  enum IpAddr {
      V4(String),
      V6(String),
  }

  let home = IpAddr::V4(String::from("127.0.0.1"));
  let loopback = IpAddr::V6(String::from("::1"));
  ```

  优点:

  - 不需要额外的结构体
  - 每个变量可以有不同的类型及其关联的数据

  ```rust
  enum IpAddr {
      V4(u8, u8, u8, u8),
      V6(String),
  }
  ```

- Option<T> 枚举

  rust 没有 null

  > 其他语言里 null 是一个值，表示"没有值"，
  > 一个变量可以处于两种状态：空值(null)或非空值
  > null 引用：Billion Dollar Mistake

  ```rust
  enum Option<T> {
      Some(T),
      None,
  }
  ```

  优点，Optional<T>比 null 好在哪：

  - Optional<T>和 T 是不同的类型，不可以把 Optional<T>当成 T 使用
  - Optional<T> 强制你处理 None 的情况

- match 控制流运算符

  ```rust
  fn test_match(coin: Coin) -> u8 {
      match coin {
          Coin::Penny => {
              println!("Lucky penny!");
              1
          }
          Coin::Nickel => 5,
          Coin::Dime => 10,
          Coin::Quarter => 25,
      }
  }
  ```

  绑定值的模式匹配

  ```rust
  fn test_match_binding(coin: Coin) -> u8 {
      match coin {
          Coin::Penny => {
              println!("Lucky penny!");
              1
          }
          Coin::Nickel => 5,
          Coin::Dime => 10,
          Coin::Quarter(state) => {
              println!("State quarter from {:?}!", state);
              25
          }
      }
  }
  ```

  匹配 Option<T>

  ```rust
  fn test_match_option(option: Option<u8>) -> u8 {
      match option {
          Some(value) => value,
          None => 0,
      }
  }
  ```

  match 匹配必须穷尽所有可能性(exhausive patterns)，否则编译器会报错

  ```rust
  fn test_match_underscore(coin: Coin) -> u8 {
      match coin {
          Coin::Penny => {
              println!("Lucky penny!");
              1
          }
          _ => 0,  // _ 匹配所有其他情况
      }
  }
  ```

- if let

  只关心一种模式，其他情况不处理，match 的语法糖

  ```rust
  if let Some(3) = some_u8_value {
      println!("three");
  }
  ```

  等价于

  ```rust
  match some_u8_value {
      Some(3) => println!("three"),
      _ => (),
  }
  ```

---

- package, crate, module，path

  - package(包): 一个 Cargo 包，可以包含多个 crate
    - 包含一个 Cargo.toml 文件，描述如何构建这些 crate
    - 只能包含 0 个或 1 个库 crate，可以包含任意数量的二进制 crate
  - crate(单元包): 一个库或二进制文件
    - src/main.rs 是 binary crate 的 crate root，crate 名与包名相同
    - src/lib.rs 是 library crate 的 crate root
    - cargo 把 crate root 文件交给 rustc 来构建库或二进制文件
  - module(模块)、use: 一个文件或多个文件，包含多个函数
    - 在一个 crate 内将代码进行分组
    - 控制作用域与私有性
  - path(路径): 模块树的 path

  ```rust
   mod front_of_house {
       mod hosting {
           fn add_to_waitlist() {}
       }
   }

   pub fn eat_at_restaurant() {
       // 绝对路径
       crate::front_of_house::hosting::add_to_waitlist();
       // 相对路径
       front_of_house::hosting::add_to_waitlist();
   }
  ```

- 私有边界(private boundary)

  - Rust 中的所有项（函数、方法、结构体、枚举、模块和常量）`默认是私有的`
  - 可以使用 pub 关键字来使项变为公有(pub 可以理解为**export**, use 可以理解为**import**)

  - 模块内部的项对外部是私有的

  ```rust
  mod front_of_house {
      pub mod hosting {
          pub fn add_to_waitlist() {}
      }
  }

  pub fn eat_at_restaurant() {
      // 绝对路径
      crate::front_of_house::hosting::add_to_waitlist();
      // 相对路径
      front_of_house::hosting::add_to_waitlist();
  }
  ```

  - 父模块无法访问子模块的私有项
  - 子模块可以访问父模块的私有项

- super

  - 使用 super 关键字来访问父模块

  ```rust
  mod front_of_house {
      pub mod hosting {
          pub fn add_to_waitlist() {}
      }
  }

  fn serve_order() {}

  mod back_of_house {
      fn fix_incorrect_order() {
          cook_order();
          super::serve_order();
      }

      fn cook_order() {}
  }
  ```

- use 关键字
  使用 use 关键字将路径引入作用域
  一般是引入父级模块，而不是子模块

  ```rust
  use crate::front_of_house::hosting;
  ```

- as 关键字
  使用 as 关键字重命名引入的路径

```rust
use std::fmt::Result;
use std::io::Result as IoResult;
```

- pub use 关键字
  重导出，类似 es6 的 export xxx from xxx

  ```rust
  mod front_of_house {
      pub mod hosting {
          pub fn add_to_waitlist() {}
      }
  }

  pub use crate::front_of_house::hosting;
  ```

- 嵌套路径批量引入
  ```rust
  use std::{self, cmp::Ordering, io};
  ```
- glob 运算符

  ```rust
  use std::collections::*;
  ```

  谨慎使用，可能引入大量不必要的项

  - 一般用于测试，将所有被测试代码引入
  - prelude 模块，Rust 标准库的 prelude 模块包含了很多常用的项，可以通过 glob 运算符引入

- 模块拆分

---

常用集合

- Vec 保存同一类型的多个值
  索引 vs get 处理访问越界

  ```rs
  // 注意不能同时存在对同一 Vec 的可变引用和不可变引用
  // 这里可能迭代器失效/扩容导致的内存重新分配
  fn main() {
      let mut v: Vec<i32> = Vec::new();
      v.push(2);
      let first: &i32 = &v[0];
      v.push(22); // cannot borrow `v` as mutable because it is also borrowed as immutable
      println!("{:?}", first);
  }
  ```

- enum 保存不同类型的值

- String

  - rust 的**核心语言**中只有一种字符串类型：str(或&str); 借用
  - String 来自**标准库**而不是核心语言，类似 java 的 StringBuilder；拥有
  - 其他类型：OsString, CString, CStr, Cow<str>
    - String vs Str 后缀，前者拥有所有权，后者借用
  - 四种修改字符串的方法
    - push_str
    - push
    - `+` 运算符
    - format! 宏(不会发生所有权转移)

- 解引用强制转换(Deref coercion)

  - Deref trait
  - Rust 可以自动调用 Deref trait 中的 deref 方法，将 &String 转换为 &str

- 字符串不支持下标访问，因为字符串是 UTF-8 编码的，一个字符可能占用多个字节；此外 String 无法保证 O(1)时间复杂度的索引访问，因为需要遍历整个字符串来确定有多少个有效字符

- rust 三种看待字符串的方式
  bytes、scalar values、grapheme clusters

  - bytes：字节
  - scalars：Unicode 标量值
  - grapheme clusters：字形簇

  ```rust
      // bytes
      for b in "你好".bytes() {
          println!("{}", b); // 228 189 160 229 165 189
      }
      // scalar value
      for c in "你好".chars() {
          println!("{}", c); // 你 好
      }
      // 标准库没有提供直接获取 grapheme clusters 的方法，需要使用第三方库
  ```

- 切割字符串

  - 不推荐使用下标切割字符串，因为可能导致 panic

  ```rust
  let s = "你好".to_string();
  let hello = &s[0..3]; // panic
  ```

  - 使用 split 方法
  - 使用 split_whitespace 方法
  - 使用 split_at 方法
  - 使用 split_inclusive 方法

- HashMap
  - 所有权
    实现了 Copy trait 的类型，值会被复制到 HashMap 中
    没有实现 Copy trait 的类型，值会被移动到 HashMap 中

---

错误处理

- 可恢复错误：Result<T, E>

  - Ok(T)
  - Err(E)

- 不可恢复错误：panic!
  默认情况下，panic 发生：

  - 程序展开调用栈（工作量大）
    - rust 沿着调用栈向上回溯，清理每个函数的数据
  - 或立即终止调用栈 (abort)
    - 不进行清理，直接退出
    - 内存需要由操作系统清理

- unwrap 和 expect
  - unwrap：如果 Result 是 Ok，则返回 Ok 中的值，否则 panic
  - expect：与 unwrap 类似，但可以指定 panic 时的错误信息

-unwrap_or_else

- unwrap_or_else：与 unwrap 类似，但可以指定一个闭包，用于生成 panic 时的错误信息

- 传播错误 propagate error
  使用 ? 运算符
  ? 运算符只能用于返回 Result 的函数
  表示如果 Result 是 Ok，则返回 Ok 的值，否则 Err 是`整个函数`的返回值

  ```rs

  fn read_username_from_file() -> Result<String, io::Error> {
    let mut s = String::new();
    File::open("foo.text")?.read_to_string(&mut s)?;
    Ok(s)
  }
  ```

  main 函数的返回类型可以是 Result，这样就可以在 main 函数中使用?运算符

- 什么时候需要 panic
  当代码可能处于损坏状态时，panic 是合适的
  bad state: 无效的值、无效的参数、不一致的状态、缺失的值

---

- 泛型
- 函数定义中的泛型
  结构体中的泛型
  枚举中的泛型
  方法定义中的泛型

- 偏特化：针对具体的类型实现方法（对满足特定条件的泛型类型实现方法）
  偏特化（Partial Specialization）是泛型编程中的一个概念，主要用于模板（在 Rust 中称为泛型）编程中。它允许开发者为模板提供一个特殊的实现，这个实现只适用于模板参数的一个子集。偏特化可以提供更加精确或优化的实现，针对特定类型或类型组合的特殊行为。

  在 Rust 中，由于语言设计的限制和安全性考虑，直接的偏特化是不支持的。**Rust 采用 trait 和 trait bounds 来实现类似偏特化的功能**。通过为泛型类型实现不同的 trait，可以根据类型特征选择不同的行为实现。

  例如，你可以为一个泛型结构体实现不同的 trait，根据类型的不同选择不同的实现：

  ```rust
  trait GenericTrait {
      fn do_something(&self);
  }

  // 对所有的 T 实现 GenericTrait
  impl<T> GenericTrait for T {
      fn do_something(&self) {
          println!("Generic implementation");
      }
  }

  // 对满足特定条件的 T 实现 GenericTrait
  impl<T: std::fmt::Display> GenericTrait for T {
      fn do_something(&self) {
          println!("Specialized implementation for Display types: {}", self);
      }
  }
  ```

  在这个例子中，尽管 Rust 不直接支持偏特化，但通过 trait 和 trait bounds，我们可以为特定的类型或满足特定条件的类型组合提供特殊的行为实现，这在实践中可以达到类似偏特化的效果。然而，需要注意的是，这种方式并不是真正的偏特化，因为 Rust 的 trait 实现是基于全局的，不能基于部分类型参数进行特化。

- 单态化（monomorphization）：为需要的类型都生成类

- trait
  trait：抽象的定义共享`行为`
  trait bounds：类型约束

---

- 生命周期
  rust 的每个引用都有自己的生命周期
  生命周期是引用有效的作用域
  生命周期注解：告诉编译器引用的生命周期

  存在的主要原因：

  - 避免悬垂引用

- 生命周期标注语法
  描述了多个引用的生命周期之间的**关系**，不会影响引用的生命周期
  以'开头，通常是短小的单词
  单个生命周期标注本身没有意义

  ```rust
  // !取生命周期最短的那个
  fn longest<'a>(x: &'a str, y: &'a str) -> &'a str {
      if x.len() > y.len() {
          x
      } else {
          y
      }
  }
  ```

- 生命周期的省略

  - 每个引用都有一个生命周期
  - 需要为使用生命周期的函数或 struct 的每个引用添加生命周期标注

  没有显示标注生命周期的情况下，编译器使用三条规则来判断引用的生命周期：

  1. 每个引用类型的参数都有它自己的生命周期参数
  2. 如果只有一个输入生命周期参数，那么它被赋给所有输出生命周期参数
  3. 如果方法有多个输入生命周期参数，不过其中之一是 &self 或 &mut self，**那么 self 的生命周期被赋予所有输出生命周期参数**

  如果编译器无法根据这些规则推断生命周期，编译器会报错

- 静态生命周期
  - 'static 是一个特殊的生命周期：整个程序的持续时间
  - 谨慎使用

---

- 自动化测试
  测试体的 3a

  - arrange 准备数据
  - act 执行代码
  - assert 验证结果

  测试函数需要使用 #[cfg(test)] 注解 #[cfg(test)] 注解告诉 Rust 只在执行 cargo test 时才编译和运行测试代码，而在运行 cargo build 时不这么做。 #[test] 注解告诉 Rust 这是一个测试函数

  ```rust
  #[cfg(test)]
  mod tests {
      #[test]
      fn it_works() {
          assert_eq!(2 + 2, 4);
      }
  }
  ```

  cargo test
  **每个测试函数都会在一个新线程中运行**
  当主线程发现某个测试函数 panic 时，会终止测试，但其他测试函数会继续运行

- 断言宏
  assert_eq!
  assert_ne!

- 自定义错误消息
  assert_eq!(2 + 2, 4, "2 + 2 不等于 4");
- 检查 panic

  ```rust
  #[should_panic]
  #[should_panic(expected = "Panic message")] // 检查 panic 的消息

  ```

- 测试中使用 Result<T, E> 返回

  ```rust
  #[cfg(test)]
  mod tests {
      use super::*;

      #[test]
      fn it_works() -> Result<(), String> {
          if 2 + 2 == 4 {
              Ok(())
          } else {
              Err(String::from("2 + 2 不等于 4"))
          }
      }
  }
  ```

- 改变测试行为
  添加命令行参数

  - 默认行为：
    1. 并行运行测试
    2. 所有测试都会运行
    3. 捕获（不显示）所有输出，使读取与测试结果相关的输出更容易
  - 命令行参数：

    1. --test-threads=1：顺序运行测试
    2. --ignored：只运行被忽略的测试
    3. --test：只运行包含 test 标记的测试
    4. --nocapture：显示测试输出

    两种：

    - 针对 cargo test 命令行的参数，紧跟在 cargo test 后面
    - 针对测试可执行文件的参数，紧跟在--后面

    ```shell
    cargo test --help
    cargo test -- --test-threads=1
    ```

- 串行、并行运行测试
- 只运行特定测试
  模式匹配

  ```shell
  cargo test one_hundred
  ```

- 忽略某些测试
  使用 ignore 属性

  ```rust
  #[test]
  #[ignore]
  fn expensive_test() {
      // code
  }
  ```

- 测试的分类
  - 单元测试：测试一个模块、一个函数或一个方法；小、专注，可测试私有函数
    `[cfg(test)]`标注
  - 集成测试：测试整个程序的某个功能；和其他外部代码一样使用你的代码，只能测试 public API
    不需要`#[cfg(test)]`标注，因为完全位于被测试库的外部(tests 目录)，集成测试的覆盖率非常重要
- 针对 binary crate 的集成测试
  只能独立运行，可以在 main.rs 中编写胶水代码

---

二进制程序关注点分离的指导性原则

- 将程序拆分为 main.rs 和 lib.rs，业务逻辑放入 lib.rs
- 命令行解析逻辑较少时，放在 main.rs 也行
- 命令行解析逻辑较多时，需要将解析逻辑从 main.rs 移动到 lib.rs

---

exit(0)表示程序正常退出；除了 0 之外，其他参数均代表程序异常退出，如：exit(1),exit(-1)。

---

Box<dyn Trait> 是 trait 对象，可以在运行时指定具体类型

---

设置环境变量

```shell
CASE_INSENSITIVE=1 cargo run to poem.txt
```

---

将错误信息输出到标准错误流

```rust
eprintln!("Application error: {}", e);
```

---

闭包：可以捕获环境的匿名函数(普通函数不能捕获)

闭包的类型推断与类型标注

- Fn Trait

  - FnOnce：取得所有权
  - FnMut：可变借用
  - Fn：不可变借用
    - 创建闭包时，Rust 会根据闭包体如何使用环境来推断闭包是 FnOnce、FnMut 还是 Fn
    - 所有的闭包都实现了 FnOnce，因为它可以被调用一次
    - 如果一个闭包没有移动环境中的任何值，那么它实现 FnMut
    - 如果一个闭包没有对环境中的任何值进行可变借用，那么它实现 Fn

- 缺点：内存开销

- move 关键字：强制闭包获取环境值的所有权
  - 当将闭包**传递给新线程以移动数据使其归新线程所有**时，此技术最为有用
    所有权移动都是为了准确的内存回收

---

- 迭代器
  iterator trait 和 next 方法

  ```rust
  trait Iterator {
      type Item;  // 关联类型

      fn next(&mut self) -> Option<Self::Item>;  // Self::Item 是关联类型
  }
  ```

  next: 返回 Option，Some 包含
  下一个值，None 表示结束

- 迭代方法：
  .iter()：不可变引用上创建迭代器，用于读取集合中的值
  .iter_mut()：迭代可变的引用，用于修改集合中的值
  .into_iter()：创建的的迭代器会获取所有权，用于移动集合中的值

- 消耗、产生迭代器

  - 消耗适配器：调用 next 方法，consumer
    例如：sum、collect
  - 产生迭代器：不调用 next 方法，supplier
    迭代器适配器
    例如：map、filter

- 闭包捕获环境

- 迭代器与循环的性能
  **迭代器性能更好**
  - Rust 的迭代器是**零成本抽象(Zero-Cost Abstraction)**，不会导致性能损失
  - Rust 的迭代器是**惰性的**，只有在需要时才会执行

---

cargo、crates.io

- release profile(发布配置) 自定义构建

  ```toml
  [profile.dev]
  opt-level = 0

  [profile.release]
  opt-level = 3
  ```

- https://crates.io/ 上发布库

  **文档注释：**

  - 用于生成文档，显示公共 API 的注释
  - 使用 ///
  - 支持 markdown

  ````rust
  /// Adds one to the number given.
  /// # Examples
  /// ```
  /// let arg = 5;
  /// let answer = my_crate::add_one(arg);
  /// assert_eq!(6, answer);
  /// ```
  /// # Panics
  /// The function will panic if the argument is 100.
  ///
  /// # Errors
  /// The function will return an error if the argument is 101.
  ///
  /// # Safety
  /// The function is not safe to call with a negative number.
  ///
  /// # Performance
  /// The function is optimized for speed, not size.
  ````

  生成 html 文档

  ```shell
  cargo doc --open // 生成文档并打开浏览器
  ```

  常用章节：
  `# Examples`
  `# Panics`
  `# Errors`
  `# Safety`

  **文档注释作为测试**

  描述包/模块的注释
  符号：//!

  ```rust
  //! # My Crate

  ```

  发包流程：

  - 获取 api tokens
    https://crates.io/settings/tokens
  - 登录
    cargo login [token]
  - 为新的 crate 添加元数据
    name 不能重名
    description
    license
    version
    author
  - 发布
    cargo publish

  注意，crate 一旦发布，就不能删除，且该版本无法覆盖(不可变 hh)
  目的：依赖于该版本的项目可继续正常工作

  - 撤回版本
    防止新项目依赖于该版本
    已经存在的项目不受影响
    yank 一个版本： cargo yank --vers 1.0.1
    取消 yank： cargo yank --vers 1.0.1 --undo

- cargo 工作空间：通过 workspaces 组织大工程

  - 一个 Cargo.toml 文件
  - 多个二进制或库 crate
  - 一个 target 目录
  - 一个 lock 文件
  - 一个输出目录

  ```shell
  cargo new my_workspace --bin
  cd my_workspace
  cargo new my_bin --bin
  cargo new my_lib --lib
  ```

  - 一个 Cargo.toml 文件

    ```toml
    [workspace]
    members = [
        "my_bin",
        "my_lib",
    ]
    ```

  - 一个 target 目录
  - 一个 lock 文件
  - 一个输出目录

- 从https://crates.io/ 安装二进制 crate
  cargo install：安装的 crate 会被放在 ~/.cargo/bin 目录下

---

- 智能指针 smart pointer
  通常用 struct 实现，并实现了
  Deref 和 Drop trait

  - Box<T>：堆上分配 (unique_ptr)
  - Rc<T>：引用计数 (shared_ptr)
  - RefCell<T>：内部可变性
  - Arc<T>：原子引用计数
  - Mutex<T>：互斥锁
  - RwLock<T>：读写锁
