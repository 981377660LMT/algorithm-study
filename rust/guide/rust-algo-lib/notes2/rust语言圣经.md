[rust 语言圣经](https://course.rs/about-book.html)

1. 包管理工具最重要的意义就是任何用户拿到你的代码，都能运行起来，而不会因为各种包版本依赖焦头烂额。
2. 很多语言中，你并不需要深入了解栈与堆。 但对于 Rust 这样的系统编程语言，值是位于栈上还是堆上非常重要, 因为这会影响程序的行为和性能。

- 栈中的所有数据都**必须占用已知且固定大小的内存空间**，假设数据大小是未知的，那么在取出数据时，你将无法取到你想要的数据
- 与栈不同，对于**大小未知或者可能变化的数据**，我们需要将它存储在堆上。
  当向堆上放入数据时，需要请求一定大小的内存空间。操作系统在堆的某处找到一块足够大的空位，把它标记为已使用，并返回一个表示该位置地址的**指针**, 该过程被称为在**堆上分配内存，有时简称为 “分配”(allocating)**。
  接着，该指针会被推入**栈**中，因为指针的大小是已知且固定的，在后续使用过程中，你将通过栈中的**指针**，来获取数据在堆上的实际内存位置，进而访问该数据。
  想象一下去餐馆就座吃饭: 进入餐馆，告知服务员有几个人，然后服务员找到一个够大的空桌子（堆上分配的内存空间）并领你们过去。如果有人来迟了，他们也可以通过桌号（栈上的指针）来找到你们坐在哪。
  在堆上分配内存则需要更多的工作，这是因为操作系统必须**首先找到一块足够存放数据的内存空间**，接着做一些记录为下一次分配做准备，如果**当前进程分配的内存页不足时(缺页中断)，还需要进行系统调用来申请更多内存**。 因此，处理器在栈上分配数据会比在堆上分配数据更加高效。

  https://github.com/youngyangyang04/TechCPP/blob/master/problems/%E4%BB%80%E4%B9%88%E6%98%AF%E7%BC%BA%E9%A1%B5%E4%B8%AD%E6%96%AD.md
  缺页中断是计算机操作系统中的一个重要概念，发生在程序访问虚拟内存时，需要加载的页面不在主存中，需要从辅存（如硬盘）中读取的情况下。当程序试图访问一个已经被映射到虚拟地址空间但尚未载入物理内存的页面时，就会引发缺页中断。

  具体来说，当程序访问一个虚拟地址时，操作系统会首先检查该地址对应的页面是否已经在主存中。如果页面在主存中，那么程序可以直接访问；如果页面不在主存中，就会触发缺页中断。此时，操作系统会进行以下步骤：

  中断处理：CPU 接收到缺页中断信号后，暂停当前正在执行的程序，将控制权交给操作系统内核。
  处理程序：操作系统会根据页面表或其他映射信息确定页面所在的位置（通常是磁盘），并将页面加载到主存中的空闲页面框中。
  更新页表：操作系统更新页表中有关该页面的信息，包括物理地址等。
  恢复程序：一旦页面加载到内存中，操作系统会重新启动之前暂停的程序，使其继续执行。

3. 变量遮蔽(shadowing)
   变量遮蔽(shadowing)是指在同一作用域中，你可以定义一个与之前变量同名的新变量，这样新变量会遮蔽之前的变量。
   变量遮蔽的用处在于，如果你在某个作用域内无需再使用之前的变量（在被遮蔽后，无法再访问到之前的同名变量），就可以重复的使用变量名字，**而不用绞尽脑汁去想更多的名字**。
   例如，假设有一个程序要统计一个空格字符串的空格数量：

   ```RUST
    // 字符串类型
    let spaces = " ";
    // usize 数值类型
    let spaces = spaces.len();

    // 变量遮蔽可以帮我们节省些脑细胞，不用去想如 str_spaces 和 num_spaces 此类的变量名
   ```

4. 对于 Rust 语言而言，这种**基于语句（statement）和表达式（expression）的方式**是非常重要的，你需要能明确的区分这两个概念, 但是对于很多其它语言而言，这两个往往无需区分。**基于表达式是函数式语言的重要特征，表达式总要返回值。**
   表达式如果不返回任何值，会隐式地返回一个 () 。

## 流程控制

如果不使用引用的话，所有权会被转移（move）到 for 语句块中，后面就无法再使用这个集合了
对于实现了 copy 特征的数组(例如 [i32; 10] )而言， for item in arr 并不会把 arr 的所有权转移，而是直接对其进行了拷贝

```rust
for item in &container {
  // ...
}
```

如果想在循环中，修改该元素，可以使用 mut 关键字：

```rust
for item in &mut collection {
  // ...
}
```

- for item in collection，转移所有权，等价于 for item in collection.into_iter()
- for item in &collection，不可变借用，等价于 for item in collection.iter()
- for item in &mut collection，可变借用，等价于 for item in collection.iter_mut()

- 两种循环方式优劣对比

  ```rust
    // 第一种
  let collection = [1, 2, 3, 4, 5];
  for i in 0..collection.len() {
    let item = collection[i];
    // ...
  }

  // 第二种
  for item in collection {

  }
  ```

  注意 arr[indedx]会因为边界检查(Bounds Checking)导致运行时的性能损耗

5. unicode 字符集与 utf-8 编码
   Unicode 和 UTF-8 是两个相关但不同的概念。

   Unicode 是一种字符集，它定义了每个字符的唯一标识符（码点）。Unicode 码点可以表示各种语言、符号和表情等字符。每个 Unicode 码点都有一个唯一的整数值，通常用 "U+" 后跟十六进制数字表示，例如 U+0041 表示字符 "A"。

   UTF-8 是一种字符编码方案，它定义了如何将 Unicode 码点编码为字节序列。UTF-8 使用变长编码，根据字符的码点范围，使用不同长度的字节序列来表示字符。UTF-8 编码的特点是兼容 ASCII 编码，对于 ASCII 字符，使用一个字节表示，而对于非 ASCII 字符，使用多个字节表示。

   下面是一个示例，展示了字符 "A" 在 Unicode 和 UTF-8 中的表示：

   Unicode 码点：U+0041
   UTF-8 编码：0x41
   在这个示例中，字符 "A" 的 Unicode 码点是 U+0041，它是一个表示 "A" 的唯一标识符。而在 UTF-8 编码中，字符 "A" 使用一个字节 0x41 来表示，因为它与 ASCII 字符集中的对应字符是相同的。

   总结来说，Unicode 是一个字符集，定义了字符的唯一标识符，而 UTF-8 是一种字符编码方案，定义了如何将 Unicode 码点编码为字节序列。UTF-8 是一种变长编码，可以有效地表示各种字符，并且兼容 ASCII 编码。

6. String 与 &str
   String: StringBuilder，具有所有权
   str: 真正的字符串，字符串字面值
   &str: 字符串切片，对 str 的不可变引用，不具有所有权

7. 操作 UTF-8 字符串
   字符（rune）
   如果你想要以 `Unicode 字符`的方式遍历字符串，最好的办法是使用 chars 方法，例如：

   ```RUST
    for c in "中国人".chars() {
        println!("{}", c);
    }
   ```

   字节（byte）
   这种方式是返回字符串的底层字节数组表现形式：

   ```RUST
    for b in "中国人".bytes() {
        println!("{}", b);
    }
   ```

   获取子串
   想要准确的从 UTF-8 字符串中获取子串是较为复杂的事情
   可以考虑尝试下这个库：utf8_slice。

8. RAII
   Rust 语言中的 RAII（Resource Acquisition Is Initialization）是一种**资源获取即初始化**的编程范式，它是 Rust 语言的一个重要特性。
   RAII 的核心思想是，资源的生命周期应该与其所在作用域的生命周期相对应，资源的释放应该在资源所在作用域结束时自动进行。这种方式可以避免资源泄漏和资源使用后未释放的问题，提高程序的可靠性和安全性。
   在 Rust 语言中，RAII 通常通过结构体和实现 Drop trait 来实现。结构体中包含了需要管理的资源，Drop trait 中实现了资源的释放逻辑。当结构体的实例超出作用域时，Drop trait 中的 drop 方法会被自动调用，从而释放资源。
   RAII 是 Rust 语言的一个重要特性，它可以帮助程序员避免资源泄漏和资源使用后未释放的问题，提高程序的可靠性和安全性。
9. 单元结构体(Unit-like Struct)
   如果你定义一个类型，但是不关心该类型的内容, 只关心它的行为时，就可以使用 单元结构体：

```RUST
struct AlwaysEqual;

let subject = AlwaysEqual;

// 我们不关心 AlwaysEqual 的字段数据，只关心它的行为，因此将它声明为单元结构体，然后再为它实现某个特征
impl SomeTrait for AlwaysEqual {
}
```

10. 结构体数据的所有权
    在之前的 User 结构体的定义中，有一处细节：**我们使用了自身拥有所有权的 String 类型而不是基于引用的 &str 字符串切片类型。**
    这是一个有意而为之的选择：因为**我们想要这个结构体拥有它所有的数据，而不是从其它地方借用数据**。
    **如果你想在结构体中使用一个引用，就必须加上生命周期，否则就会报错：**
11. 任何类型的数据都可以放入枚举成员中: 例如字符串、数值、结构体甚至另一个枚举。
    利用枚举同一化类型。
12. Option 确实是一个比较不错的设计，在 JavaScript 中比较混淆的就有 undefined 和 null，undefined 代表变量的值未定义，null 代表变量的值默认为空值。

## 模式匹配

## 方法

## 泛型和特征

## 集合类型

## 生命周期

生命周期语法用来将函数的多个引用参数和返回值的作用域关联到一起，一旦关联到一起后，Rust 就拥有充分的信息来确保我们的操作是内存安全的。
有时候，'static 确实可以帮助我们解决非常复杂的生命周期问题甚至是无法被手动解决的生命周期问题，那么此时就应该放心大胆的用，只要你确定：你的所有引用的生命周期都是正确的，只是编译器太笨不懂罢了。

总结下：
生命周期 'static 意味着能和程序活得一样久，例如字符串字面量和特征对象
实在遇到解决不了的生命周期标注问题，可以尝试 T: 'static，有时候它会给你奇迹

&'static 表示的是任意一种从头活到尾的类型比如 string，实际用途更偏向于指针，表示指向一种生命周期为'static 的数据类型，T:'static 更偏向于泛型，用于定义一种类型

## 返回值和错误处理

## 包和模块

## 注释和文档

## 格式化输出

## 文件搜索工具

1. 所有的用户输入都不可信！不可信！不可信！
2. 关注点分离(Separation of Concerns)
   main.rs 负责启动程序，lib.rs 负责逻辑代码的运行。从测试的角度而言，这种分离也非常合理： lib.rs 中的主体逻辑代码可以得到简单且充分的测试，至于 main.rs ？确实没办法针对其编写额外的测试代码，但是它的代码也很少啊，很容易就能保证它的正确性。
3. 返回 Result 来替代直接 panic
4. 多用迭代器、模式匹配

# 闭包与迭代器

1. 捕获作用域中的值
2. 闭包对内存的影响
3. 三种 Fn 特征
   不可变借用：Fn
   可变借用：FnMut
   转移所有权：FnOnce

# 智能指针

1. Rc/Arc 是不可变引用；如果要修改，需要配合后面章节的内部可变性 RefCell 或互斥锁 Mutex
2. Cell 和 RefCell 在功能上没有区别，区别在于 Cell<T> 适用于 T 实现 Copy 的情况,在实际开发中，Cell 使用的并不多，因为我们要解决的往往是可变、不可变引用共存导致的问题，此时就需要借助于 RefCell 来达成目的
   RefCell 正是用于你确信代码是正确的，而编译器却发生了误判时。
   当你确信编译器误报但不知道该如何解决时，或者你有一个引用类型，需要被四处使用和修改然后导致借用关系难以管理时，都可以优先考虑使用 RefCell
   当非要使用内部可变性时，首选 Cell，只有你的类型没有实现 Copy 时，才去选择 RefCell
3. 内部可变性
   对一个不可变的值进行可变借用

# 循环引用与自引用

1. Rust 的安全性是众所周知的，但是不代表它不会内存泄漏。一个典型的例子就是同时使用 Rc<T> 和 RefCell<T> 创建循环引用，最终这些引用的计数都无法被归零，因此 Rc<T> 拥有的值也不会被释放清理。
   当你使用 RefCell<Rc<T>> 或者类似的类型嵌套组合（具备内部可变性和引用计数）时，就要打起万分精神
2. Weak 非常类似于 Rc，但是与 Rc 持有所有权不同，Weak 不持有所有权，它仅仅保存一份指向数据的弱引用
   使用方式简单总结下：**对于父子引用关系，可以让父节点通过 Rc 来引用子节点，然后让子节点通过 Weak 来引用父节点**。

   Weak 通过 use std::rc::Weak 来引入，它具有以下特点:

   - 可访问，但没有所有权，不增加引用计数，因此不会影响被引用值的释放回收
   - 可由 Rc<T> 调用 downgrade 方法转换成 Weak<T>
   - Weak<T> 可使用 upgrade 方法转换成 Option<Rc<T>>，如果资源已经被释放，则 Option 的值是 None
   - 常用于解决循环引用的问题

   ```rust
   #[derive(Debug)]
   struct Node {
       value: i32,
       parent: RefCell<Weak<Node>>,
       children: RefCell<Vec<Rc<Node>>>,
   }
   ```

3. unsafe 解决循环引用
   虽然 unsafe 不安全，但是在各种库的代码中依然很常见用它来实现自引用结构，主要优点如下:

   性能高，毕竟直接用裸指针操作
   代码更简单更符合直觉: 对比下 Option<Rc<RefCell<Node>>>

# unsafe

FFI(Foreign Function Interface) 外部函数接口
FFI 之所以存在是由于现实中很多代码库都是由不同语言编写的，如果我们需要使用某个库，但是它是由其它语言编写的，那么往往只有两个选择：

- 对该库进行重写或者移植
- 使用 FFI

当然，除了 FFI 还有一个办法可以解决`跨语言调用的问题`，那就是将其作为一个独立的服务，然后使用网络调用的方式去访问，HTTP，gRPC 都可以。

```rs
// 调用 C 标准库中的 abs 函数：
extern "C" {
    // ABI(Application Binary Interface) 是指应用程序二进制接口，它定义了函数的调用约定，包括参数传递、返回值等
    fn abs(input: i32) -> i32;
}

fn main() {
    unsafe {
        println!("Absolute value of -3 according to C: {}", abs(-3));
    }
}
```

---

- 命名规范
  https://course.rs/practice/naming.html

  - 变量命名
    对于 type-level 的构造 Rust 倾向于使用驼峰命名法，而对于 value-level 的构造使用蛇形命名法
    对于驼峰命名法，复合词的缩略形式我们认为是一个单独的词语，所以只对首字母进行大写：**使用 Uuid 而不是 UUID**，Usize 而不是 USize，Stdin 而不是 StdIn。
    对于蛇形命名法，缩略词用全小写：**is_xid_start。**
    对于蛇形命名法（包括全大写的 SCREAMING_SNAKE_CASE），除了最后一部分，其它部分的词语都**不能由单个字母组成： btree_map 而不是 b_tree_map，PI_2 而不是 PI2.**

  - 特征命名
    特征的名称应该使用动词，而不是形容词或者名词，例如 Print 和 Draw 明显好于 Printable 和 Drawable。
  - 类型转换要遵守 `as_，to_，into_` 命名惯例(C-CONV)
    类型转换应该通过方法调用的方式实现，其中的前缀规则如下：

    - `as_`：表示类型转换，无性能开销，borrowed -> borrowed.
    - `to_`：表示类型转换，性能开销大，返回一个新的值。
    - `into_`：表示类型转换，转换本身是零消耗的，ownership 转移，返回一个新的值。

    例子：

    - str::as_bytes() 把 str 变成 UTF-8 字节数组，性能开销是 0。输入是一个借用的 &str，输出也是一个借用的 &str
    - Path::to_str 会执行一次昂贵的 UTF-8 字节数组检查，输入和输出都是借用的。对于这种情况，如果把方法命名为 as_str 是不正确的，因为这个方法的开销还挺大
    - str::to_lowercase() 在调用过程中会遍历字符串的字符，且可能会分配新的内存对象。输入是一个借用的 str，输出是一个有独立所有权的 String
    - String::into_bytes() 返回 String 底层的 Vec<u8> 数组，转换本身是零消耗的。该方法获取 String 的所有权，然后返回一个新的有独立所有权的 Vec<u8>
    - 当一个单独的值被某个类型所包装时，访问该类型的内部值应通过 `into_inner()`方法来访问

    `如果 mut 限定符在返回类型中出现，那么在命名上也应该体现出来。`
    例如，Vec::as_mut_slice 就说明它返回了一个 mut 切片，在这种情况下 as_mut_slice 比 as_slice_mut 更适合

    ```rust
    // 返回类型是一个 `mut` 切片
    fn as_mut_slice(&mut self) -> &mut [T];
    ```

  - 读访问器(Getter)的名称遵循 Rust 的命名规范(C-GETTER)
    除了少数例外，在 Rust 代码中 **get 前缀不用于 Getter**。

    ```rust
    pub struct S {
        first: First,
        second: Second,
    }

    impl S {
        // 而不是 get_first
        pub fn first(&self) -> &First {
            &self.first
        }

        // 而不是 get_first_mut，get_mut_first，or mut_first
        pub fn first_mut(&mut self) -> &mut First {
            &mut self.first
        }
    }
    ```

    当有且仅有一个值能被 Getter 所获取时，才使用 get 前缀。
    例如，Cell::get 能直接访问到 Cell 中的内容。

  - 一个集合上的方法，如果返回迭代器，需遵循命名规则：iter，iter_mut，into_iter (C-ITER)

  ```RUST
  fn iter(&self) -> Iter             // Iter implements Iterator<Item = &U>
  fn iter_mut(&mut self) -> IterMut  // IterMut implements Iterator<Item = &mut U>
  fn into_iter(self) -> IntoIter     // IntoIter implements Iterator<Item = U>
  ```

  - Cargo Feature 的名称不应该包含占位词(C-FEATURE)
    不要在 Cargo feature 中包含无法传达任何意义的词，例如 use-abc 或 with-abc，直接命名为 abc 即可。
  - 命名要使用**一致性的词序**(C-WORD-ORDER)
    这是一些标准库中的错误类型:

    JoinPathsError
    ParseBoolError
    ParseCharError
    ParseFloatError
    ParseIntError
    RecvTimeoutError
    StripPrefixError
    它们都使用了 **谓语-宾语-错误** 的词序，如果我们想要表达一个网络地址无法分析的错误，由于词序一致性的原则，命名应该如下 ParseAddrError，而不是 AddrParseError。
    词序和个人习惯有很大关系，想要注意的是，你可以选择合适的词序，**但是要在包的范畴内保持一致性，就如标准库中的包一样**。

---

- #[no_mangle]，它用于告诉 Rust 编译器：不要乱改函数的名称
