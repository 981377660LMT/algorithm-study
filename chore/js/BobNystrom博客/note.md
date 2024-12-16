https://journal.stuffwithstuff.com/archive/

1.  C# Extension Methods
    无侵入式扩展方法

    - “Extension method” = “friendlier calling convention”

    - Reuse methods without inheritance
      重用方法而不继承
      `类似 rust 的 为某个 trait 实现方法`
      “if this class provides this capability, then it also has this capability”.
    - prefer static methods of a helper class to instance methods
      `更喜欢工具类的静态方法而不是实例方法`

2.  Checking Flags in C# Enums
    I like C# enums and I also like using them as `bitfields`
3.  Using an Iterator as a Game Loop
    使用迭代器作为游戏循环

    ```cpp
    void GameLoop()
    {
        while (mPlaying)
        {
            HandleUserInput();
            UpdateGameState();
            Render();
        }
    }
    ```

    - Separating out the UI 分离 UI
      没有显式调用 Render() ，而是让引擎在事情发生时引发事件（怪物移动、使用物品等）
      ```cpp
      void ProcessGame(UserInput input)
      {
          HandleUserInput(input);
          UpdateGameState();
      }
      ```
    - Enter iterators 输入迭代器

      ```cpp
      IEnumerable<bool> CreateProcessIterator()
      {
          while (true)
          {
              foreach (Entity entity in Entities)
              {
                  if (entity.NeedsUserInput)
                  {
                      yield return true;
                  }
                  else
                  {
                      entity.Move();
                  }
              }

              foreach (Item item in Items)
              {
                  item.Move();
              }
          }
      }
      ```

    总结：

    - 通常实现的方式是游戏循环为每个演员每轮提供一小部分时间。每个实体都会有一个 Update()方法，该方法，执行一步然后返回；缺点是需要维护状态
    - 如果您的系统支持协程，由于`纤程维护自己的整个调用堆栈`，因此您甚至可以从其他函数调用中在它们之间进行切换

4.  Amaranth, an Open Source Roguelike in C#
    Amaranth，一款用 C# 编写的开源 Roguelike 游戏

    - Because I’m crazy about decoupling, it’s actually split into three separate projects:
      因为我对解耦很着迷，所以它实际上分为三个独立的项目

5.  Naming Things in Code 在代码中命名事物
    When I’m designing software, I spend a lot of time thinking about names. For me, thinking about names is inseparable from the process of design. To name something is to define it.
    当我设计软件时，我会花很多时间思考名字。对我来说，对名字的思考，离不开设计的过程。命名某物就是定义它。

    If a class doesn’t represent something concrete, consider a metaphor.
    如果一个类不代表具体的东西，考虑一个隐喻。

    - Bad:
      - IncomingMessageQueue
      - CharacterArray
      - SpatialOrganizer
    - Good:
      - Mailbox
      - String
      - Map

    If you use a metaphor, use it consistently.
    如果你使用一个比喻，请始终如一地使用它。

    - Bad: Mailbox, DestinationID
    - Good: Mailbox, Address

    Name functions that return a Boolean (i.e. predicates) like questions.
    命名返回布尔值（即谓词）的函数，如问题。

    - Bad: list.Empty();
    - Good: list.IsEmpty();
      list.Contains(item);

    Name functions that return a value and don’t change state using nouns.
    **命名返回值,且不使用名词更改状态的函数(仅查询)。**

    - Bad: list.GetCount();
    - Good: list.Count();

    Don’t make the name redundant with an argument.
    **不要用参数使名称多余。**

    - Bad:
      list.AddItem(item);
      handler.ReceiveMessage(msg);
    - Good:
      list.Add(item);
      handler.Receive(msg);

    Don’t use “and” or “or” in a function name.
    **不要在函数名称中使用“and”或“or”。**
    如果在名称中使用连词，则该函数可能做得太多。
    将其分解成较小的部分并相应地命名。如果要确保这是一个原子操作，请考虑为整个操作创建一个名称，或者可能创建一个封装它的类。

    - Bad:
      mail.VerifyAddressAndSendStatus();
    - Good:
      mail.VerifyAddress();
      mail.SendStatus();

    **命名良好的代码更容易与其他程序员讨论，有助于传播代码知识；**
    一个带有良好名称部件的模块可以快速教会您它的作用。通过只阅读一小部分代码，您将快速构建整个系统的完整心理模型。如果它将某物称为“邮箱”，您将希望看到“邮件”和“地址”，而无需阅读它们的代码。
    另一方面，糟糕的名称会创建一堵不透明的代码墙，迫使您在脑海中煞费苦心地运行该程序，观察其行为，然后创建自己的私人命名法。“哦，DoCheck（）看起来正在遍历连接，看看它们是否都是实时的。我称之为 AreConnectionsLive（）“。这不仅速度慢，而且不可转移。

    **当我无法命名某物时，很有可能我试图命名的东西设计不佳。也许它试图一次做太多事情，或者错过了一个关键部分来完成它。**
    很难说我的设计是否好，但我发现我`做得不好的最可靠的指南之一是当名字不容易出现时。`
    当我现在设计时，我试着注意这一点。`一旦我对名字感到满意，我通常会对设计感到满意。`

    Name functions that return a value and don’t change state using nouns.
    **命名返回值,且不使用名词更改状态的函数(仅查询)。**

    - Bad: list.GetCount();
    - Good: list.Count();

    Don’t make the name redundant with an argument.
    **不要用参数使名称多余。**

    - Bad:
      list.AddItem(item);
      handler.ReceiveMessage(msg);
    - Good:
      list.Add(item);
      handler.Receive(msg);

    Don’t use “and” or “or” in a function name.
    **不要在函数名称中使用“and”或“or”。**
    如果在名称中使用连词，则该函数可能做得太多。
    将其分解成较小的部分并相应地命名。如果要确保这是一个原子操作，请考虑为整个操作创建一个名称，或者可能创建一个封装它的类。

    - Bad:
      mail.VerifyAddressAndSendStatus();
    - Good:
      mail.VerifyAddress();
      mail.SendStatus();

    **命名良好的代码更容易与其他程序员讨论，有助于传播代码知识；**
    一个带有良好名称部件的模块可以快速教会您它的作用。通过只阅读一小部分代码，您将快速构建整个系统的完整心理模型。如果它将某物称为“邮箱”，您将希望看到“邮件”和“地址”，而无需阅读它们的代码。
    另一方面，糟糕的名称会创建一堵不透明的代码墙，迫使您在脑海中煞费苦心地运行该程序，观察其行为，然后创建自己的私人命名法。“哦，DoCheck（）看起来正在遍历连接，看看它们是否都是实时的。我称之为 AreConnectionsLive（）“。这不仅速度慢，而且不可转移。

    **当我无法命名某物时，很有可能我试图命名的东西设计不佳。也许它试图一次做太多事情，或者错过了一个关键部分来完成它。**
    很难说我的设计是否好，但我发现我`做得不好的最可靠的指南之一是当名字不容易出现时。`
    当我现在设计时，我试着注意这一点。`一旦我对名字感到满意，我通常会对设计感到满意。`

6.  Class in the Front, Proto in the Back
    Class 在前面，Proto 在后面
7.  “void”, “null”, “Maybe” and “nothing”
    每种编程语言都必须为两种缺失提供一种机制：当函数总是不返回任何内容时，以及当它有时不返回任何内容时。
    新兴语言不使用 null，而是使用 Maybe/Option 类型。这是一个包装器，它要么包含一个值，要么包含一个空值。
8.  **Future-Proofing**, Uniform Access, and Masquerades
    面向未来、统一访问和伪装

    如果我们希望它们是可变的，并且在它们更改时没有任何验证或逻辑可执行，该怎么办？
    除开封装 getter、setter，还有哪些做法是`面向未来(future-proofing)`的？

    - Hiding constructors behind factories
      将构造者隐藏在工厂后面

      ```CS
      class PersonFactory {
        public Person create() {
          return new Person();
        }
      }

      void doSomething(PersonFactory factory) {
        Person person = factory.create();
        // ...
      }
      ```

    - Hiding classes behind interfaces
      将类隐藏在接口后面

    - Why do we do this?
      我们不面向未来。假设我们到处都使用公共字段。然后后来我们意识到我们确实需要验证该字段。我们必须做些什么来做出这种改变？
    - Depth not volume 深度而非体积
      哪个更容易：在一个文件中更改一百行代码，还是在一百个文件中更改一行代码？

    - The uniform access principle(统一访问原则)
      模块提供的所有服务都应该通过统一的符号提供，这不会暴露它们是通过存储还是通过计算实现的

    Name functions that return a value and don’t change state using nouns.
    **命名返回值,且不使用名词更改状态的函数(仅查询)。**

    - Bad: list.GetCount();
    - Good: list.Count();

    Don’t make the name redundant with an argument.
    **不要用参数使名称多余。**

    - Bad:
      list.AddItem(item);
      handler.ReceiveMessage(msg);
    - Good:
      list.Add(item);
      handler.Receive(msg);

    Don’t use “and” or “or” in a function name.
    **不要在函数名称中使用“and”或“or”。**
    如果在名称中使用连词，则该函数可能做得太多。
    将其分解成较小的部分并相应地命名。如果要确保这是一个原子操作，请考虑为整个操作创建一个名称，或者可能创建一个封装它的类。

    - Bad:
      mail.VerifyAddressAndSendStatus();
    - Good:
      mail.VerifyAddress();
      mail.SendStatus();

    **命名良好的代码更容易与其他程序员讨论，有助于传播代码知识；**
    一个带有良好名称部件的模块可以快速教会您它的作用。通过只阅读一小部分代码，您将快速构建整个系统的完整心理模型。如果它将某物称为“邮箱”，您将希望看到“邮件”和“地址”，而无需阅读它们的代码。
    另一方面，糟糕的名称会创建一堵不透明的代码墙，迫使您在脑海中煞费苦心地运行该程序，观察其行为，然后创建自己的私人命名法。“哦，DoCheck（）看起来正在遍历连接，看看它们是否都是实时的。我称之为 AreConnectionsLive（）“。这不仅速度慢，而且不可转移。

    **当我无法命名某物时，很有可能我试图命名的东西设计不佳。也许它试图一次做太多事情，或者错过了一个关键部分来完成它。**
    很难说我的设计是否好，但我发现我`做得不好的最可靠的指南之一是当名字不容易出现时。`
    当我现在设计时，我试着注意这一点。`一旦我对名字感到满意，我通常会对设计感到满意。`

9.  Class in the Front, Proto in the Back
    Class 在前面，Proto 在后面
10. “void”, “null”, “Maybe” and “nothing”
    每种编程语言都必须为两种缺失提供一种机制：当函数总是不返回任何内容时，以及当它有时不返回任何内容时。
    新兴语言不使用 null，而是使用 Maybe/Option 类型。这是一个包装器，它要么包含一个值，要么包含一个空值。
11. **Future-Proofing**, Uniform Access, and Masquerades
    面向未来、统一访问和伪装

    如果我们希望它们是可变的，并且在它们更改时没有任何验证或逻辑可执行，该怎么办？
    除开封装 getter、setter，还有哪些做法是`面向未来(future-proofing)`的？

    - Hiding constructors behind factories
      将构造者隐藏在工厂后面

      ```CS
      class PersonFactory {
        public Person create() {
          return new Person();
        }
      }

      void doSomething(PersonFactory factory) {
        Person person = factory.create();
        // ...
      }
      ```

    - Hiding classes behind interfaces
      将类隐藏在接口后面

    - Why do we do this?
      我们不面向未来。假设我们到处都使用公共字段。然后后来我们意识到我们确实需要验证该字段。我们必须做些什么来做出这种改变？
    - Depth not volume 深度而非体积
      哪个更容易：在一个文件中更改一百行代码，还是在一百个文件中更改一行代码？

    - The uniform access principle(统一访问原则)
      模块提供的所有服务都应该通过统一的符号提供，这不会暴露它们是通过存储还是通过计算实现的

12. 用户能否定义自己的抽象，在语法上与内置行为相同？我们可以用我们自己的逻辑替换默认语言提供的行为，而不必更改调用约定吗？
    在所有主流静态 OOP 语言中，调用 new 总是返回固定类的实例。没有办法用我们自己的逻辑来代替它。我们被工厂困住了。
13. 如何轻松地将新数据类型和新行为添加到现有系统？
    假设我们正在编写一个文档编辑器。我们有几种可以使用的文档：文本、绘图和电子表格。我们需要对文档执行一些操作：将其绘制到屏幕上、加载它并将其保存到光盘。它们形成一个网格，如下所示：

    ```
                Text       Drawing   Spreadsheet
            ┌───────────┬───────────┬───────────┐
    draw()  │           │           │           │
            ├───────────┼───────────┼───────────┤
    load()  │           │           │           │
            ├───────────┼───────────┼───────────┤
    save()  │           │           │           │
            └───────────┴───────────┴───────────┘

    ```

    There are a couple of questions to answer:
    有几个问题需要回答：

    - How do we organize the code for this?
      我们如何为此组织代码？

    - How do we add new columns—new types of documents?
      我们如何添加新列——新类型的文档？

    - How do we add new rows—new operations you can perform on any document?
      我们如何添加新行——可以在任何文档上执行的新操作？

    - How do we ensure all of the cells are implemented?
      我们如何确保所有单元都得到实施？

    Magpie’s answers for the original four questions are:
    喜鹊对原四个问题的回答是：

    - How do we organize code? However you like. Put stuff together where it makes sense.
      我们如何组织代码？不过你喜欢就好。将东西放在有意义的地方。

    - How do we add new columns—new types of documents? Like a typical OOP language: define a new class. If it has the necessary methods, it’s a Document.
      我们如何添加新列——新类型的文档？就像典型的 OOP 语言一样：定义一个新类。如果它具有必要的方法，那么它就是一个 Document 。

    - How do we add new rows—new operations you can perform on any document? Add new methods to the classes that need them. This can be done outside of the file where the class is defined.
      我们如何添加新行——可以在任何文档上执行的新操作？将新方法添加到需要它们的类中。`这可以在定义类的文件外部完成。`

    - How do we ensure all of the cells are covered? Add the new operation to the interface too. The static checker will then make sure only classes that have the operation are used in places that expect a Document.
      我们如何确保所有单元格都被覆盖？也将新操作添加到界面中。然后，`静态检查器将确保只有具有该操作的类才会在需要 Document 的地方使用。`

    When you’re defining things, you get the flexibility of a dynamic language. Before it runs, though, you get the safety of a static language.
    当您定义事物时，您将获得动态语言的灵活性。不过，在它运行之前，您可以获得静态语言的安全性。

    **有点像 rust?**

14. The Language I Wish Go Was

    - Named/keyword arguments 命名/关键字参数
    - Block arguments 块参数
    - Operator overloading 运算符重载

    Go has two really neat type system features: implicitly implemented interfaces and a flat type hierarchy. There are two other simple additions I’d dig: tuples and unions.
    Go 有两个非常简洁的类型系统功能：隐式实现的接口和平面类型层次结构。我还想添加另外两个简单的补充：元组和联合。

    - Tuples 元组
    - Unions 联合
    - Constructors 构造函数
      构造函数的主要职责是初始化结构体的所有字段。在分配构造函数中的字段之前访问该字段是一个静态错误，并且在方法体末尾未能分配给所有字段也是一个静态错误。
    - Eliminating nil 消除 nil
      原因和 Null 一样

    - Error-handling 错误处理
      In addition to the aforementioned automatic stack-unwinding, there’s two other things I like about exceptions:
      除了前面提到的自动堆栈展开之外，我还喜欢异常的另外两点：

      - No zombie variables 没有僵尸变量
        如果您进入 file.read() ，您就知道您有一个有效的文件。换句话说，我们可以使用块来限制变量的范围，使其仅在正确的情况下才存在。
      - Type safety without coupling 无耦合类型安全
        我将一个我定义的类型的对象传递给某个第三方库，然后第三方库调用该对象的方法
        出错时，第三方库不需要关心我的错误类型；
        但是用错误码的方式，第三方库需要将该错误存储在某些 interface{}中，接收代码必须进行动态转换。
        换句话说，我们必须放弃类型安全。

    - Generics 泛型
    - future-proofing
      But for many other things, Go has taken steps backwards:
      但对于许多其他事情，Go 已经倒退了：

      - Field access is different from method calls (which always take ()).
        字段访问与方法调用不同（方法调用始终采用() ）。
      - Subscript syntax like array[index] cannot be overloaded (unlike C++ and C#).
        像 array[index]这样的下标语法不能重载（与 C++ 和 C# 不同）。
      - Object allocation uses special new syntax and can only zero-initialize the object.
        对象分配使用特殊的 new 语法，并且只能对对象进行零初始化。如果稍后需要更复杂的初始化，则必须将每个 new(Foo)调用替换为 NewFoo() 。

      Sometimes a little syntactic sugar goes a long way.

      Right now, Go avoids this by having a culture of not future-proofing. That culture is only sustainable as long as all of the code that your code touches is very easy for you to modify. That’s true within some small or very agile organizations, but once Go starts moving to wider enterprise use, I fear we’ll start seeing “best practices” like “always wrap every field in a getter method” and “always hide constructors behind New**” functions and then it’s Java all over again.
      现在，Go 通过拥有一种不面向未来的文化来避免这种情况。只有当您的代码涉及的所有代码都非常容易修改时，这种文化才可持续。在一些小型或非常敏捷的组织中确实如此，但一旦 Go 开始转向更广泛的企业使用，我担心我们将开始看到“最佳实践”，例如“始终将每个字段包装在 getter 方法中”和“始终将构造函数隐藏在 New**后面”函数，然后又是 Java。

15. getMemberType 和 canAssignFrom
16. 将实例创建为单个原子操作。为此，我们需要传入实例所需的所有字段，它将返回一个完全初始化的、可供使用的对象。这就是 construct 。

17. https://journal.stuffwithstuff.com/2011/03/19/pratt-parsers-expression-parsing-made-easy/
    Pratt Parsers: Expression Parsing Made Easy
    Pratt 解析器：表达式解析变得简单
18. ''We use a common base class, because we have/had no better tool in mind for sharing functions among classes.

    - “我们使用一个公共基类，因为我们没有更好的工具来在类之间共享函数。
    - 不允许两次继承同一个类(无环)
    - **允许多继承，但最后必须保证依赖是树，而不是图(Mixin?)**

    - Do we really need a common base class for all objects?
      我们真的需要所有对象都有一个公共基类吗？
    - No common base class => no deadly diamonds
      没有共同的基类 => 没有致命的钻石
    - Multimethods: dynamic dispatch
      多种方法：动态调度
    - Extenible classes: No member member functions
      可扩展类：无成员成员函数

19. Higher Order Macros in C++
    C++ 中的高阶宏(X Macro)
    https://en.wikipedia.org/wiki/X_Macro
20. Rooms and Mazes: A Procedural Dungeon Generator
    房间和迷宫：程序地下城生成器
21. **What color is your function?**

    - 五个规则
      - 同步函数返回值，异步函数不返回值，而是调用回调。
      - 同步函数将其结果作为返回值提供，异步函数通过调用您传递给它的回调来提供结果。
      - 您不能从同步函数调用 async 函数，因为在 async 函数稍后完成之前，您将无法确定结果。
      - 由于回调的原因，异步函数不会在表达式中组合，具有不同的错误处理，并且不能与 try/catch 一起使用，也不能在许多其他控制流语句中使用。
      - Node 的全部问题是核心库都是异步的。（尽管他们确实减少了这个功能并开始添加很多东西的 \_\_\_Sync（） 版本。
    - 什么语言没有颜色？
      JS、Dart、C# 和 Python 都有异步传染问题。
      Go、Lua 和 Ruby 没有。

      他们有什么共同点吗？
      Threads. Or, more precisely: `multiple independent callstacks that can be switched between`. It isn’t strictly necessary for them to be operating system threads. Goroutines in Go, coroutines in Lua, and fibers in Ruby are perfectly adequate.
      线程。或者，更准确地说：可以在多个独立的调用堆栈之间切换。它们并非绝对必要成为操作系统线程。`Go 中的 goroutines、Lua 中的 coroutines 和 Ruby 中的 fibers 都足够了。`

      The fundamental problem is `“How do you pick up where you left off when an operation completes”?` You’ve built up some big callstack and then you call some IO operation. For performance, that operation uses the operating system’s underlying asynchronous API. You cannot wait for it to complete because it won’t. You have to return all the way back to your language’s event loop and give the OS some time to spin before it will be done.
      `根本问题是“当操作完成时，您如何从上次中断的地方继续”？`您已经构建了一些大的调用堆栈，然后调用了一些 IO 操作。为了提高性能，该操作使用操作系统的基础异步 API。您不能等待它完成，因为它不会。您必须一直返回到语言的事件循环，并在完成之前给操作系统一些时间来旋转。

      Once operation completes, you need to resume what you were doing. `The usual way a language “remembers where it is” is the callstack. That tracks all of the functions that are currently being invoked and where the instruction pointer is in each one.`
      操作完成后，您需要继续执行您正在执行的操作。`语言 “记住它的位置” 的通常方式是 callstack。这将跟踪当前正在调用的所有函数以及指令指针在每个函数中的位置。`

      But to do async IO, you have to unwind and discard the entire C callstack.
      但要执行异步 IO，您必须展开并丢弃整个 C 调用堆栈。

      Node with its ever-marching-to-the-right callbacks stuffs all of those callframes in closures. When you do:
      Node 及其不断向右行进的回调将所有这些调用帧塞入闭包中。当您执行以下操作时：

      ```js
      function makeSundae(callback) {
        scoopIceCream(function (iceCream) {
          warmUpCaramel(function (caramel) {
            callback(pourOnIceCream(iceCream, caramel))
          })
        })
      }
      ```

      Each of those function expressions closes over all of its surrounding context. `That moves parameters like iceCream and caramel off the callstack and onto the heap. When the outer function returns and the callstack is trashed, it’s cool. That data is still floating around the heap.`
      这些函数表达式中的每一个都在其所有周围的上下文中关闭。这会将 iceCream 和 caramel 等参数从调用堆栈移动到堆上。当外部函数返回并且 callstack 被丢弃时，这很酷。这些数据仍然在堆周围浮动。

      The problem is you have to manually reify every damn one of these steps. There’s actually a name for this transformation: `continuation-passing style`. It was invented by language hackers in the 70s as an intermediate representation to use in the internals of their compilers. It’s a really bizarro way to represent code that happens to make some compiler optimizations easier to do.
      `问题是你必须手动具体化这些步骤中的每一个。这种转换实际上有一个名称： continuation-passing style` 。它是语言黑客在 70 年代发明的，作为用于编译器内部的中间表示。这是一种非常奇怪的代码表示方式，恰好使某些编译器优化更容易实现。

      Async-await does help. If you peel back your compiler’s skull and see what it’s doing when it hits an await call you’d see it actually doing the CPS-transform. That’s why you need to use await in C#: it’s a clue to the compiler to say, “break the function in half here”. Everything after the await gets hoisted into a new function that the compiler synthesizes on your behalf.
      Async-await 确实有帮助。如果你剥开编译器的头骨，`看看它在点击 await 调用时在做什么，你会看到它实际上在做 CPS 转换`。这就是为什么你需要在 C# 中使用 await：编译器可以说“在这里将函数分成两半”。等待之后的一切 被提升到编译器代表您合成的新函数中。

      `This is why async-await didn’t need any runtime support in the .NET framework. The compiler compiles it away to a series of chained closures that it can already handle.` (Interestingly, closures themselves also don’t need runtime support. They get compiled to anonymous classes. In C#, closures really are a poor man’s objects.)
      `这就是 async-await 在 .NET Framework 中不需要任何运行时支持的原因。编译器将其编译为一系列它已经可以处理的链式闭包。`（有趣的是，闭包本身也不需要运行时支持。它们被编译为匿名类。在 C# 中，闭包确实是穷人的对象。

      `使用回调、promise、async-await 和生成器，你最终会把你的异步函数拖成一堆闭包，这些闭包都存在于堆中。`

22. 代码审查
    谷歌做的一件聪明的事情是严格的代码审查。

    - 首先，团队中的某个人进行常规的审查，以确保代码完成了它应该完成的任务。
    - 但是，还有第二层审查，称为可读性。它确保代码是可读的：它容易理解和维护吗？它是否符合语言的风格和习惯用法？是否有详细记录？

    A name has two goals:
    一个名字有两个目标：

    It needs to be clear: you need to know what the name refers to.
    这需要明确：你需要知道这个名字指的是什么。

    It needs to be precise: you need to know what it does not refer to.
    它需要精确：你需要知道它不指什么。

    1. Omit words that are obvious given a variable’s or parameter’s type
       省略 那些在给定变量或参数类型时显而易见的单词
       特别是，对于集合，使用复数名词描述内容总是比使用单数名词描述集合更好
       这也适用于方法名称。方法名不需要描述它的参数或它们的类型--参数列表会为你做这些。

       ```java
       // Bad:
        mergeTableCells(List<TableCell> cells)
        sortEventsUsingComparator(List<Event> events,
            Comparator<Event> comparator)

        // Better:
        merge(List<TableCell> cells)
        sort(List<Event> events, Comparator<Event> comparator)
       ```

    2. Omit words that don’t disambiguate the name
       省略不能消除名称歧义的单词
       有些人倾向于把他们所知道的一切都塞进它的名字里。
       记住，名称是一个标识符：它指向定义它的位置。它并不是一个详尽的目录，读者可能想知道的关于对象的一切。
       定义就是这样。这个名字只是让他们到达那里。

       ```java
       // Bad:
       finalBattleMostDangerousBossMonster;
       weaklingFirstEncounterMonster;

       // Better:
       boss;
       firstMonster;
       ```

    3. Omit words that are known from the surrounding context
       省略从上下文中已知的单词

       ```JAVA
       // Bad:
       class AnnualHolidaySale {
         int _annualSaleRebate;
         void promoteHolidaySale() { ... }
       }

       // Better:
       class AnnualHolidaySale {
         int _rebate;
         void promote() { ... }
       }
       ```

    4. Omit words that don’t mean much of anything
       省略那些没有什么意义的词
       In many cases, the words carry no meaningful information. They’re just fluff or jargon.
       Usual suspects include: data, state, amount, value, manager, engine, object, entity, and instance.
       在许多情况下，这些词没有任何有意义的信息。它们只是些废话或者行话。
       可疑对象包括：数据、状态、数量、值、管理器、引擎、对象、实体和实例。
       Ask yourself “Would this identifier mean the same thing if I removed the word?”
       If so, the word doesn’t carry its weight. Vote if off the island.
       **问问你自己：“如果我把这个词去掉，这个标识符的意思是一样的吗？”如果是这样的话，这个词就没有分量了。**

    https://www.swift.org/documentation/api-design-guidelines/#omit-needless-words
    这种事情在阅读代码时相当分散注意力

23. 假设我们有两个切片类型[]E1 和[]E2，分别是 E1 和 E2 元素的切片。如果 E1 可以赋值给 E2，是否意味着[]E1 可以赋值给[]E2？可赋值性是否从内部类型“传播”到外部类型？
24. **每一种具有子类型的语言中，函数类型在其返回类型中是协变的，在其参数类型中是逆变的。**
    但是在 Go 中，函数类型是不变的。
    不仅仅是函数类型。Go 语言中所有的复合类型都是`不变`的：数组、切片、通道、映射、函数。因此，基础类型（不包含任何其他类型的类型）具有某种类似于子类型的可赋值性概念。
    但是一旦你把一个类型包装在另一个类型中，任何可赋值的概念都消失了。
    https://journal.stuffwithstuff.com/2023/10/19/does-go-have-subtyping/

    为什么？需要看 golang 中接口的底层实现

    Go 语言有子类型，但它不支持变量，所有的复合类型都是`不变的`

25. 你可以想象一种语言需要三样东西：
    Non-uniform representation: Values in memory take up only as much space as they need and avoid pointer indirection when possible to maximize runtime efficiency.
    `非统一表示`：内存中的值只占用它们所需的空间，并尽可能避免指针间接，以最大限度地提高运行时效率。

    Polymorphism: The ability to reuse code to work with a range of values of different types.
    `多态性`：重用代码以处理一系列不同类型的值的能力。

    Variance: Sort of the “lifted” form of polymorphism: The ability to reuse code to work with composite types that contain a range of inner types.
    `Variance`：类似于多态性的“提升”形式：重用代码以处理包含一系列内部类型的复合类型的能力

    大多数面向对象的语言牺牲了第一个来得到另外两个。这为您提供了灵活性和表现力，但运行时成本遍布整个程序。

    golang 牺牲了 Variance，但在个体值级别上保持多态性。
    与隐式转换相结合，可以实现非统一表示。
    在这三者中，Variance 对用户来说可能是最没有价值的，所以我认为这是一个非常聪明的权衡。
