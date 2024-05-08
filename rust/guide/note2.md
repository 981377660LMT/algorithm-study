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

- trait
  trait：抽象的定义共享行为
  trait bounds：类型约束
- 生命周期
