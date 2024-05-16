`24天从Node转Rust`
https://zhuanlan.zhihu.com/p/456650665

1. 第 1 天：从 nvm 到 rustup
2. 第 2 天：从 npm 到 cargo
3. 第 3 天：设置 VS Code
4. 第 4 天：Hello World (以及你的前两个疑惑 )
5. 第 5 天：借用和所有权

6. 第 6 天：字符串，第 1 部分
7. → 第 7 天：语法和语言，第 1 部分
8. 第 8 天：语言，第 2 部分：从对象(objects), 类(classes) 到散列映射(HashMaps) 和结构(structs)
   我们不能像 TypeScript 那样，写一个返回值是 string | undefined 的函数。
   作为对比，Rust 有**枚举类型**，这也就是 Option 的用处所在.
   数据与行为的分离很重要
9. 第 9 天：语言，第 3 部分：适用于 Rust 结构（+枚举！）的"类方法"
   枚举(enums)是 Rust 解决缺少联合类型(union types)问题的一种方案
   另一种方案是特质(traits)。
   `两者的差别在于，想要的到底是类型，还是行为的一个子集。`
10. 第 10 天：从混合(Mixins) 到特质(Traits)
    Rust 的特质(traits) 很像 JavaScript 的混合(mixins)，他们都是一组方法（或方法签名）
    特质仅仅是一组方法而已。
    函数返回值动态类型不能使用 impl [Trai]t，而是使用 dyn [trait]。
11. 第 11 天：模块系统
12. 第 12 天：字符串，第 2 部分
13. 第 13 天：结果(Result) 和选项(Option)
14. 第 14 天：管理错误
15. 第 15 天：闭包
16. 第 16 天：生命周期，引用，和 'static
17. 第 17 天：数组，循环和迭代器
18. 第 18 天：异步
19. 第 19 天：开始一个大型项目
20. 第 20 天：命令行界面(CLI)参数和日志
21. 第 21 天：创建和运行 WebAssembly
22. 第 22 天：使用 JSON
23. 第 23 天：骗过借用检查器
24. 第 24 天：箱(Crates)和工具
