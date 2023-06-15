## C++11_14 高级编程 BOOST 程序库探秘

### 新特性

- 左值: 有名字的对象，可以取地址的对象，可以放在赋值号左边的对象。
- 右值: 临时对象，没有名字的对象，不能取地址的对象，不能放在赋值号左边的对象；生命周期即将结束的对象。
- **转移语义**：只要对象被 move()标记为可移动，就可以使用 std::move()将其转换为右值引用，从而调用移动构造函数或移动赋值操作符，不用担心深拷贝的性能问题。
- **完美转发**: 既可以接受左值，也可以接受右值的函数参数。std::`forward`()函数可以将左值引用转换为左值引用，将右值引用转换为右值引用。

- 自动类型推导 auto/decltype (declare type)
  auto 类似于 any
  decltype 类似于 typeof
- NULL 与 nullptr
  NULL 是一个宏定义，值为 0 (是一个整数，有安全问题)
  nullptr 是一个关键字，值为 0
- 统一使用{}初始化
- default 和 delete
  default: 显示声明默认构造函数
  delete: 禁止使用默认构造函数
- **类型别名**：using/typedef
  using 是加强版的 typedef，可以用来定义别名模板，而 typedef 不行。
- 编译期常量 constexpr
  const 和 constexpr 的区别：
  const 是运行时不可修改的变量，而 constexpr 是编译期常量，在编译期就能确定其值。
- 静态断言 static_assert
  特点是编译期断言，只能用于编译期，不能用于运行期。运行期断言用 assert()。通常需要配合 type_traits 库使用。
  static_assert(表达式, "错误信息")
- 可变参数模板
- lambda 表达式 []
  ```cpp
  std::for_each(v.begin(), v.end(), [](int i) { std::cout << i << std::endl; });
  ```
  捕获外部变量
  [**capture list**] (params) mutable exception attribute -> ret { body }
- thread_local
  用于声明线程局部变量，每个线程都有一份独立的实例，互不干扰。
- noexcept
  noexcept 无异常保证
- 属性 [[attribute]] 标记编译特征，方便编译器优化代码
  [[carries_dependency]] 传递依赖
  [[deprecated]] 废弃
  [[fallthrough]] 贯穿
  [[nodiscard]] 不要忽略返回值
  [[maybe_unused]] 可能未使用
  [[likely]] 可能性大
  [[unlikely]] 可能性小
  [[no_unique_address]] 空类优化

### 模板元编程

result_of: 用于获取函数调用的返回值类型
unwarp_reference: 用于获取引用类型的原始类型
call_traits: 用于获取函数调用的参数类型

### 类型特征萃取

type_traits: 用于萃取类型`特征`,位于 boost/type_traits.hpp

- 检查基本类型
- 复合类型
- 元数据属性
- 操作符重载属性

### 实用工具

- compressed_pair: 用于压缩两个对象的存储空间
  空类优化技术

### 迭代器

迭代器模式：
提供一种方法顺序访问一个容器对象中的各个元素，而又不**暴露**该对象的内部表示(对用户不透明)。

- 四个基本接口:
  **初始化、前进、是否结束、访问当前元素**
  cpp: begin/end/++/--/!=/==/\
  js: next()/done/value
  python: iter()/next()
- 更强大的接口：
  比较迭代位置(距离)、前进后退 n 个位置、交换两个迭代器所指元素的值
- 迭代器分类

  1. 输入迭代器：只读，从迭代器读取
  2. 输出迭代器：只写，向迭代器输出
  3. 前向迭代器：可读写，只能单向前进 (forward_list)
  4. 双向迭代器：可读写，可前进后退 (list/set/map)
  5. 随机访问迭代器：可读写，可前进后退，可取下标 (vector/deque/string)

  为了区分不同类型的迭代器，STL 提供了 iterator_traits 模板，用于萃取迭代器的类型特征。采用 tag 类(空类)作为迭代器的标签，如
  struct input_iterator_tag{};

- 迭代器适配器(adaptor)
  1. 逆向迭代器 reverse_iterator
  2. 插入迭代器 back_inserter
  3. 流迭代器 istream_iterator/ostream_iterator
- next 和 prior
  next: 前进 n 个 distance
  prior: 后退 n 个 distance
  虽然简单，但是**统一了迭代器的操作方式**
  因为只有随机访问迭代器才能够编写 iter+n 这样的代码，**不利于泛型编程**

### 区间

只想操作容器中的一部分元素，而不是整个容器，这时候就需要区间(range/islice)了。

- 操作函数
  begin/end/rbegin/rend/size/empty/distance

### 函数对象

- hash 计算任意 cpp 对象的哈希值
- mem_fn
- factory 消除了 new 的使用

### 指针容器

???

### 侵入式容器

非侵入式:不要求对容纳的元素做修改，使用方便
侵入式:对内存管理要求低，没有非侵入式容器的拷贝和克隆等要求

- 侵入式容器不负责内存分配, 元素的创建时容器之外的事情
- 侵入式容器并不是真正的容纳对象，只是用指针链接起来的对象，"链接视图"

### 多索引容器

### 流处理

### 泛型编程

### 预处理元编程

### 现代 c++开发浅谈

## CPP-17-STL-cookbook

https://github.com/xiaoweiChen/CPP-17-STL-cookbook/tree/master/content
