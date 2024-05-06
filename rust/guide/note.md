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
  创建方式：[type; size]
