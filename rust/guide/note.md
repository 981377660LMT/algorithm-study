https://www.bilibili.com/video/BV15y421h7j7

# 简介与安装更新

- 更新 rustup update
- 卸载 rustup self uninstall
- 添加组件 rustup component add rustfmt
- 查看版本 rustup --version

两种版本:

- stable: 稳定版
- nightly: 每天构建的最新版本
  安装：rustup install stable/nightly
  切换：rustup default stable/nightly

# 编译器、包管理工具、开发环境搭建

- 编译器：rustc
  查看版本：rustc --version
  编译生成二进制文件：rustc -o output_filename source_filename.rs
  编译生成库文件：rustc --crate-type=lib lib source_filename.rs

- 开发环境：
  vscode + rust-analyzer 插件
- 包管理工具：cargo
  隐式地使用 rustc 编译
  查看版本：cargo --version
  创建项目：cargo new project_name
  创建库项目：cargo new --lib project_name
  构建项目：cargo build
  运行项目：cargo run
  检查项目：cargo check
  测试项目：cargo test
  发布项目(生成优化后的二进制文件)：**cargo build --release**

  cargo.toml 文件
  [package] 项目名、版本、作者、描述
  [dependencies]、[dev-dependencies]、[build-dependencies] 依赖

# 第三方库(crate)

- https://crates.io/
- 修改 cargo.toml 文件 (不推荐)
  [dependencies]
  crate_name = "version"
  保存后自动加载依赖
- cargo-edit 工具（推荐）
  安装：cargo install cargo-edit
  添加依赖：cargo add crate_name
  添加指定版本依赖：cargo add crate_name@x.x.x
  添加开发依赖：cargo add --dev crate_name
  添加构建依赖：cargo add --build crate_name
  删除依赖：cargo rm crate_name
- 设置国内源 修改 cargo/config 文件
  rsproxy.cn

# 变量

- 命名：变量 snake_case，枚举和结构体 PascalCase
- 默认不可变，有助于防止数据竞争和并发问题
- shadowing variables：隐藏一个变量，并不是重新赋值
- 命名空间：小括号

# 常量 const 与静态变量 static

- 常量
  - 必须指定类型与值，在编译时已知；
  - UPPER_SNAKE_CASE 命名
  - 与 cpp 宏不同，rust 常量的值被直接嵌入到生成的底层机器代码中，而不是进行简单的字符替换
  - 块级作用域
- 静态变量

  - 运行时分配内存
  - 可以使用 unsafe 修改
    不要过多研究 unsafe code，不符合 rust 风格

    ```rust
    static mut MY_STATIC: i32 = 42;

    fn main() {
        unsafe {
            MY_STATIC = 100;
            println!("MY_STATIC: {}", MY_STATIC);
        }
    }
    ```

  - 静态变量的生命周期是整个程序的生命周期

# 基础数据类型

1. Integer Types
   **默认推断为 i32**
   i8, i16, i32, i64, i128
2. Unsigned Integer Types
   u8, u16, u32, u64, u128
3. Platform-Specific Types (由平台决定)
   isize, usize
4. Floating-Point Types
   - f32, f64
   - 尽量用 f64，除非你清楚边界需要空间
5. Boolean Type
   bool
6. Character Type
   - char，单引号，4 字节，支持 unicode 字符(Unicode Scalar Value)

# 元组与数组

- 元组与数组长度固定是复合类型(compound types)，区别于集合类型(collection types) Vec 和 Map
- Tuples 类型可以不同，数组类型必须相同

---

- 数组
  创建方式：[初始值; size]

  ```rust
  // 1. 直接指定每个元素的值
  let arr1 = [1, 2, 3, 4, 5];

  // 2. 初始化一个有固定大小的数组，所有元素初始值相同
  let arr2 = [0; 10]; // 创建一个有10个元素的数组，每个元素的初始值都是0

  // 3. 利用变量创建数组时，可以使用类型注解来指定数组的类型和长度
  let arr3: [i32; 3] = [6, 7, 8];
  ```

# 内存管理模型

- 模型

  - C/C++
    手动的，写错了是你菜
    new+delete，reference counting
  - python/java/C#/go
    交给 GC
    安全但 `stop the world` 对性能伤害巨大
  - rust
    the rust conpiler，最特殊的那个
    **Ownership(所有权)**, Borrowing(借用), Lifetimes
    RAII，自动内存管理，零成本抽象

- 性能
  理论：C > rust > Cpp(虚表，运行时多态有性能开销)
- c/c++内存错误大全

  1. 内存泄漏 （memory leak）

  ```c
  int* ptr = new int;
  // forget to delete
  ```

  2. 悬空指针（dangling pointer）
     悬空指针（Dangling Pointer）是指向已经释放或无效内存的指针。
     在内存被释放后，原本指向该内存的指针并没有被清除或更新，它仍然保留着对那块内存的地址引用。如果程序继续使用这样的指针，就会导致不可预测的行为，因为那块内存可能已经被重新分配给其他用途或标记为不可访问，这种情况下对悬空指针的任何读写操作都可能导致程序崩溃或数据损坏。

  ```c
  int* ptr = new int;
  delete ptr;
  // ptr is dangling pointer
  ```

  3. 重复释放（double free）

  ```c
  int* ptr = new int;
  delete ptr;
  delete ptr;
  ```

  4. 数组越界（array out of bounds）

  ```c
  int arr[10];
  arr[10] = 0;
  ```

  5. 野指针（wild pointer）
     野指针是指向“垃圾”内存或无效内存区域的指针

  ```c
  int* ptr;  // 指针变量未被初始化，直接使用
  *ptr = 10;
  ```

  6. 使用已经释放的内存（use after free）

  ```c
  int* ptr = new int;
  delete ptr;
  *ptr = 10; // use after free
  ```

  7. stack overflow
  8. 不匹配的 new/delete 和 malloc/free

- rust 内存管理模型
  - 所有权系统（Ownership System）
    - 每个值都有一个所有者
    - 所有者离开作用域，值被丢弃
  - 借用（Borrowing）
    - 不可变借用（&T）
    - 可变借用（&mut T）
  - 生命周期（Lifetimes）
    - 确保引用总是有效
  - 引用计数（Reference Counting）
    - Rc<T>，Arc<T>

---

Rust 的内存管理模型基于以下几个核心概念：

1. **所有权（Ownership）**：Rust 中的每个值都有一个被称为其所有者的变量。值在任一时刻只能有一个所有者。当所有者离开作用域，这个值将被丢弃。

2. **借用（Borrowing）**：Rust 允许通过引用来`借用值，而不取得其所有权`。借用分为不可变借用（`&T`）和可变借用（`&mut T`），不可变借用允许读取值，不允许修改；可变借用允许修改值。

3. **生命周期（Lifetimes）**：Rust 用生命周期来确保引用总是有效的。生命周期是 Rust 编译器用于确保所有的借用都是有效的，防止`悬垂引用（dangling references）`的一种机制。
   > 悬垂引用（Dangling References）是指指向已经释放或无效内存的引用。在 Rust 中，编译器通过所有权（Ownership）、借用（Borrowing）和生命周期（Lifetimes）的规则来防止悬垂引用的产生。这些规则确保在引用一个值的时候，那个值是有效的，从而避免了悬垂引用带来的安全风险。
4. **RAII（Resource Acquisition Is Initialization）**：资源获取即初始化。Rust 通过这个模型来管理资源（如内存），确保资源被正确初始化，并在不再使用时自动释放。

5. **自动内存管理**：Rust 通过所有权系统自动管理内存，编译器在编译时就会根据所有权规则检查代码，以确保安全地管理内存，避免内存泄漏和悬垂指针等问题。

6. **零成本抽象**：Rust 的内存管理是在编译时进行的，不需要运行时的垃圾收集器。这意味着 Rust 能够在不牺牲性能的情况下提供内存安全保证。

通过这些机制，Rust 能够在编译时期就避免许多内存安全问题，如双重释放、空悬指针、内存泄漏等，而无需运行时的垃圾收集器，从而在保证安全的同时，还能提供接近系统编程语言的性能。
