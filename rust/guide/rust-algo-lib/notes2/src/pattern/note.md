https://refactoringguru.cn/design-patterns/rust

## 创建型模式提供`创建`对象的机制， 增加已有代码的灵活性和可复用性。

1. 工厂方法模式
   你需要有一个既能够创建新对象， 又可以重用现有对象的普通方法
2. 建造者模式
   一般是类的关联函数(静态方法)

## 结构型模式介绍如何将对象和类`组装`成较大的结构， 并同时保持结构的灵活和高效。

## 行为模式负责对象间的高效`通信`和职责委派。

---

重构

https://refactoringguru.cn/refactoring/what-is-refactoring

1. Clean code 干净的代码
   显而易见
   不包含重复
   更少的代码就是更少的东西。更少的代码就是更少的维护。代码越少，错误就越少。代码是责任，保持简短。
2. Causes of technical debt 技术债务的原因
   Business pressure 业务压力
   Lack of understanding of the consequences of technical debt 缺乏对技术债务后果的理解
   Failing to combat the strict coherence of components 未能与组件的严格一致性作斗争
   Lack of tests 缺乏测试
   Lack of documentation 缺乏文档
   Lack of interaction between team members 团队成员之间缺乏互动
   Long-term simultaneous development in several branches 多个分支的长期同时发展
   Delayed refactoring 延迟重构
   Lack of compliance monitoring 缺乏合规性监控
3. When to refactor 何时重构

   - Rule of Three 三分法则
     When you're doing something for the first time, just get it done.
     当你第一次做某事时，就把它做好。
     When you're doing something similar for the second time, cringe at having to repeat but do the same thing anyway.
     当你第二次做类似的事情时，对不得不重复但无论如何都要做同样的事情感到畏缩。
     When you're doing something for the third time, start refactoring.
     当你第三次做某事时，开始重构。
   - When adding a feature
     添加功能时
     Refactoring helps you understand other people's code. If you have to deal with someone else's dirty code, try to refactor it first. Clean code is much easier to grasp. You will improve it not only for yourself but also for those who use it after you.
     **重构可帮助你理解其他人的代码。如果你必须处理别人的脏代码，请先尝试重构它**。干净的代码更容易掌握。您不仅会为自己改进它，还会为您之后使用它的人改进它。
     Refactoring makes it easier to add new features. It's much easier to make changes in clean code.
     重构使添加新功能变得更加容易。在干净的代码中进行更改要容易得多。
   - When fixing a bug
     修复错误时
   - During a code review
     在代码审查期间

4. How to refactor 如何重构
   The code should become cleaner.
   代码应该变得更干净。
   New functionality shouldn't be created during refactoring.
   不应在重构期间创建新功能。
   All existing tests must pass after refactoring.
   重构后，所有现有测试都必须通过。

5. Code Smells 代码异味

- Bloaters 膨胀者
  膨胀器是代码、方法和类，它们已经增加到如此`庞大`的比例，以至于它们很难使用。通常，这些气味不会立即出现，而是随着程序的发展而随着时间的推移而积累（尤其是当没有人努力根除它们时）。

  - Long Method 长方法 >一个屏幕
  - Large Class 大类
  - Primitive Obsession 基本类型偏爱
  - Long Parameter List 长参数列表 >=4
  - Data Clumps 数据泥团
    Sometimes different parts of the code contain identical groups of variables (such as parameters for connecting to a database). These clumps should be turned into their own classes.
    有时，代码的不同部分包含相同的变量组（例如用于连接到数据库的参数）。**这些团块应该变成它们自己的类。**

- Object-Orientation Abusers 面向对象滥用者
  All these smells are incomplete or incorrect application of object-oriented programming principles.
  所有这些气味都是对面向对象编程原则的不完整或不正确的应用。

  - Switch Statements Switch 语句
    当你看到 switch 时，你应该想到多态性
  - Temporary Field 临时字段
  - Refused Bequest 拒绝继承
    如果子类仅使用从其父类继承的一些方法和属性，则层次结构将不正常。不需要的方法可能只是未使用或被重新定义并发出异常。
    有人之所以在类之间创建继承，只是因为希望在超类中重用代码。但是超类和子类是完全不同的。
    如果继承没有意义，并且子类确实与超类没有任何共同之处，请消除继承，以支持将继承替换为委托。
  - Alternative Classes with Different Interfaces 具有不同接口的替代类
    两个类执行相同的函数，但具有不同的方法名称。

- Change Preventers 变更预防者
  你需要在代码中的一个地方进行更改，你也必须在其他地方进行许多更改。因此，程序开发变得更加复杂和昂贵。

  - Divergent Change 发散变化
  - Shotgun Surgery 散弹手术
  - Parallel Inheritance Hierarchies 平行继承体系

- Dispensables 可废弃物
  可有可无的东西是毫无意义和不需要的东西，如果没有它将使代码更简洁、更高效、更易于理解。

  - Comments 注释
  - Duplicate Code 重复代码
  - Lazy Class 懒惰类
    理解和维护课程总是需要时间和金钱。因此，如果一个班级做得不够吸引你的注意力，就应该删除它。
  - Data Class 数据类
  - Dead Code 死代码
  - Speculative Generality 猜测性通用性

- Couplers 耦合者

  - Feature Envy 特性嫉妒
    一个方法访问另一个对象的数据比访问它自己的数据更多。
  - Inappropriate Intimacy 不恰当的亲密关系
  - Message Chains 消息链
    在代码中，您会看到一系列类似于 $a->b()->c()->d()
  - Middle Man 中间人
    如果一个类只执行一个操作，将工作委派给另一个类，为什么它存在？

- Other smells 其他异味
  Incomplete Library Class 不完整的库类

6. 重构技术

- Composing Methods 组合方法
  提取、删除方法，消除重复代码
  对不同的值使用不同的变量。每个变量应该只负责一件特定的事情。
  将 Method 替换为 Method 对象
- Moving Features between Objects 在对象之间移动特征
  移动方法、属性
  隐藏委托(Hide Delegate)
  删除中间人
- Organizing Data 组织数据
  将“单向关联”更改为“双向关联”
- Simplifying Conditional Expressions 简化条件表达式
  Replace Nested Conditional with Guard Clauses
  将嵌套条件替换为保护子句，你应该有一个“扁平”的条件列表，一个接一个。
- Simplifying Method Calls 简化方法调用
  将构造函数替换为工厂方法
- Dealing with Generalization 处理泛化
  组合代替继承
