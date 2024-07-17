https://go7hic.github.io/A-Philosophy-of-Software-Design/#/
https://zjykzk.github.io/posts/cs/design/aposd/
https://book.douban.com/review/13169504/
https://blog.csdn.net/yh88623131/category_12078681.html

作者写了这本书，并开设了斯坦福 CS 190 的课程
**软件开发，一个是应对复杂性，一个是应对变化**
《软件复杂度应对之道》这个书名可能会更贴切一些。

## 第一章：介绍

如果软件开发人员应始终考虑设计问题，而**降低复杂性**是软件设计中最重要的要素，
则软件开发人员应始终考虑复杂性。
`这本书是关于如何使用复杂性来指导软件设计的整个生命周期。`

两种解决复杂性的通用方法

1. 使代码更简单和更明显来消除复杂性
2. 封装它（模块化设计）

改善设计技能的最好方法之一就是学会识别危险信号：信号表明一段代码可能比需要的复杂。

## 第二章：复杂性的本质(The Nature of Complexity)

- 复杂性的定义
  软件复杂度是任何和`系统结构相关的并且让系统难以理解和修改`的事情
  Complexity is anything related to the `structure of a software system` that makes it hard to `understand and modify` the system.

  复杂性可以采取多种形式。

  - `难以理解`：很难理解一段代码是如何工作的
  - `难以修改`：需要花费很多精力才能实现较小的改进，或者可能不清楚必须修改系统的哪些部分才能进行改进，如果不引入其他错误，可能很难修复

  读者比作家更容易理解复杂性
  作为开发人员，您的工作不仅是创建可以轻松使用的代码，而且还要创建其他人也可以轻松使用的代码

- 复杂性的症状

  - 修改放大(Change amplification，令人恼火的)
    看似简单的变更需要在许多不同地方进行代码修改(例如，没有统一的配置文件)
    模块依赖，缺少抽象时，就容易发生。
  - 认知负担(Cognitive load，增加改变的成本)
    指开发人员需要多少知识才能完成一项任务。较高的认知负担意味着开发人员必须花更多的时间来学习所需的信息(例如，垂悬指针带来的内存泄漏)
    有时，需要更多代码行的方法实际上更简单，因为它减少了认知负担
  - 未知的未知(Unknown unknowns，`最糟糕的`，“坑”)
    必须修改哪些代码才能完成任务，或者开发人员必须获得哪些信息才能成功地执行任务，这些都是不明显的

  修改放大，如果告诉你具体需要修改哪些地方，那么还是能够正确实现功能的。
  同样，认知负担，如果能够明确依赖哪些信息，也是能够理解代码的。
  不确定是否修改正确是最糟糕的，除了通读代码全局了解以外，你是没法确定怎么修改。

  良好设计的最重要目标之一就是使系统`显而易见`
  一个显而易见的系统是，开发人员可以在不费力地思考的情况下快速猜测要做什么，同时又可以确信该猜测是正确的
  `第 18 章讨论使代码更明显的技术`

- 复杂性的原因(两个)

  - `依赖性`(dependencies)
    无法孤立地理解和修改给定的一段代码
    软件设计的目标之一是`减少依赖关系的数量，并使依赖关系保持尽可能简单和明显`
  - `模糊性`(obscurity)
    当重要的信息不明显时，就会发生模糊(文档不足，变量名不清晰，代码不一致)
    减少模糊性的最佳方法是简化系统设计

  本质：复杂度是递增的 -> “零容忍”

- 总结
  Complexity comes from an accumulation of `dependencies and obscurities`.
  As complexity increases, it leads to `change amplification, a high cognitive load, and unknown unknowns.`
  As a result, it takes more code modifications to implement each new feature.
  In addition, developers spend more time acquiring enough information to make the change safely and, in the worst case, they can’t even find all the information they need.
  The bottom line is that complexity makes it difficult and risky to `modify an existing code base`.

## 第三章：工作代码是不够的，战略与战术编程（Working Code Isn’t Enough，Strategic vs. Tactical Programming)

如果您进行战术编程，则每个编程任务都会带来一些此类复杂性
几乎每个软件开发组织都有至少一个将战术编程发挥到极致的开发人员：战术龙卷风
成为一名优秀的软件设计师的第一步是要意识到仅工作代码是不够的。`引入不必要的复杂性以更快地完成当前任务是不可接受的`。最重要的是系统的长期结构。任何系统中的大多数代码都是通过扩展现有代码库编写的，因此，作为开发人员，最重要的工作就是促进这些将来的`扩展`。
`战略性编程需要一种投资心态。您必须花费时间来改进系统的设计，而不是采取最快的方式来完成当前的项目。`
这些投资会在短期内让您放慢脚步，但从长远来看会加快您的速度。
不是实施想到的第一个想法，请尝试几种替代设计并选择最简洁的设计。

- 投资多少？
  大量的前期投资（例如尝试设计整个系统）将不会有效
  随着您对系统的了解，`理想的设计趋于零碎出现`。因此，最好的方法是连续进行大量小额投资
  将总开发时间的 10％到 20％用于投资
- 创业与投资
  一旦代码库变成了意大利面条，几乎是不可能修复的。您可能会为产品的使用寿命付出高昂的开发成本
  最好的工程师对良好的设计深感兴趣。如果您的代码库很残酷，那么单词会变得毫无用处，这将使您难以招募
  Facebook 是否能够成功清除多年来战术编程中积累的问题：“快速行动并打破困境” -〉“以坚实的基础架构快速移动”
  在一家关心软件设计并拥有清晰代码基础的公司中工作会有趣得多
- 总结
  Good design doesn’t come for free.
  It has to be something you invest in continually, so that small problems don’t accumulate into big ones.
  The most effective approach is one where `every engineer makes continuous small investments in good design`.

## 第四章：模块应该是深的(Modules Should Be Deep)

- 模块化设计(modular design)
  模块是具有接口和实现的任何代码单元，可以指一个服务、一个包、一个类甚至一个函数
  模块化设计的目标是最大程度地减少模块之间的依赖性
  每个模块独立处理其中一部分复杂度
  为了管理依赖关系，我们将每个模块分为两个部分：接口和实现，接口描述模块能“做什么”，实现部分描述“怎么做”
  **最好的模块是那些其接口比其实现简单得多的模块**
  `如果开发人员需要了解特定信息才能使用模块，则该信息是模块接口的一部分`
  明确指定接口的好处之一是，它可以`准确指示开发人员使用关联模块所需要知道的内容`(有助于消除 unknown unknowns 问题)

  `接口的设计要让常规的场景变得简单`。像 JAVA 中的从一个文件名创建一个有缓存的数据流需要特别的构造 BufferedInputStream 使得缓存接口的使用变得复杂，容易遗漏这个环节。

- 抽象
  抽象与模块化设计的思想紧密相关，模块化设计本质上就是在做抽象这件事情
  抽象是实体的简化视图，其中`省略了不重要的细节`。抽象是有用的，因为它们`使我们更容易思考和操纵复杂的事物`
  在模块化编程中，`每个模块以其接口的形式提供抽象`
  抽象错误的两种情况：包含并非真正重要的细节->增加认知负担；忽略了真正重要的细节->未知的未知
- 深模块(good)
  `最好的模块很深：它们允许通过简单的接口访问许多功能。`
  The best modules are those that provide powerful functionality yet have simple interfaces.
  I use the term deep to describe such modules.
  接口很 narrow，内涵很丰富
  浅层模块是具有相对复杂的接口的模块，但功能不多：它不会掩盖太多的复杂性
  模块深度是考虑成本与收益的一种方式。模块提供的好处是其功能。模块的成本（就系统复杂性而言）是其接口。`模块的接口代表了模块强加给系统其余部分的复杂性：接口越小越简单，引入的复杂性就越小`。最好的模块是那些收益最大，成本最低的模块。接口不错，但更多或更大的接口不一定更好！

  Unix 操作系统及其后代（例如 Linux）提供的文件 I/O 机制是深层接口的一个很好的例子.
  多年来，Unix I/O 接口的实现已经发生了根本的发展，但是五个基本内核调用并没有改变

  ```cpp
  int open(const char* path, int flags, mode_t permissions);
  ssize_t read(int fd, void* buffer, size_t count);
  ssize_t write(int fd, const void* buffer, size_t count);
  off_t lseek(int fd, off_t offset, int referencePosition);
  int close(int fd);
  ```

  深度模块的另一个示例是诸如 Go 或 Java 之类的语言中的垃圾收集器。这个模块根本没有接口。
  垃圾收集器的实现非常复杂，但是使用该语言的程序员无法发现这种复杂性。
  诸如 Unix I/O 和垃圾收集器之类的深层模块`提供了强大的抽象，因为它们易于使用，但隐藏了巨大的实现复杂性。`

- 浅模块(bad)
  浅层模块是其接口与其提供的功能相比相对复杂的模块。
  浅类有时是不可避免的，但是它们`在管理复杂性方面没有提供太多帮助`。

  > 多维表格有太多这样的浅类，增加了认知负担，没有提供任何帮助。

  ```java
  private void addNullValueForAttribute(String attribute) {
      data.put(attribute, null);
  }
  ```

  从管理复杂性的角度来看，此方法会使情况变得更糟，而不是更好。`该方法不提供任何抽象，因为其所有功能都可以通过其接口看到。`
  `该方法增加了复杂性（以供开发人员学习的新接口的形式），但没有提供任何补偿。`

- 经典主义
  不幸的是，深度类的价值在今天并未得到广泛认可。
  `编程中的传统观点是，类应该小而不是深`。经常告诉学生，类设计中最重要的事情是将较大的类分成较小的类。
  对于方法，通常会给出相同的建议：“任何长于 N 行的方法都应分为多种方法”（N 可以低至 10）。
  这种方法导致了大量的浅类和方法，这增加了整体系统的复杂性。(?)
  小类不会贡献太多功能，因此必须有很多小类，每个小类都有自己的接口。
  `这些接口的累积会在系统级别产生巨大的复杂性`。小类也导致冗长的编程风格，这是由于每个类都需要样板。

  - bad case: Java 读取文件的装饰器模式
    最常见的分类病实例之一是 Java 类库。Java 语言不需要很多小类，但是分类文化似乎已在 Java 编程社区中扎根。
    例如，要打开文件以便从文件中读取序列化的对象，必须创建三个不同的对象：

    ```java
    FileInputStream fileStream = new FileInputStream(fileName);
    BufferedInputStream bufferedStream = new BufferedInputStream(fileStream);
    ObjectInputStream objectStream = new ObjectInputStream(bufferedStream);
    ```

    FileInputStream 对象仅提供基本的 I/O：它不能执行缓冲的 I/O，也不能读取或写入序列化的对象。
    BufferedInputStream 对象将`缓冲`添加到 FileInputStream，而 ObjectInputStream 添加了`读取和写入序列化对象`的功能。
    一旦文件被打开，上面代码中的前两个对象 fileStream 和 bufferedStream 将永远不会被使用。以后的所有操作都使用 objectStream。

    Java 开发人员会争辩说，并不是每个人都希望对文件 I/O 使用缓冲，因此不应将其内置到基本机制中。
    他们可能会争辩说，最好分开保持缓冲，以便人们可以选择是否使用它。
    **提供选择是好的，但是应该设计接口以使常见情况尽可能简单（请参阅第 6 页的公式）。**
    **几乎每个文件 I/O 用户都希望缓冲，因此默认情况下应提供缓冲。对于不需要缓冲的少数情况，该库可以提供一种禁用它的机制。**

  - good case: Unix 文件
    Unix 系统调用的设计者使常见情况变得简单。
    例如，他们认识到`顺序 I/O 是最常见的，因此他们将其作为默认行为`。
    使用 lseek 系统调用，随机访问仍然相对容易实现，但是仅执行顺序访问的开发人员无需了解该机制。
    **如果一个接口具有许多功能，但是大多数开发人员只需要了解其中的一些功能，那么该接口的有效复杂性就是`常用功能的复杂性`。**

- 总结
  By separating the interface of a module from its implementation, we can hide the complexity of the implementation from the rest of the system.
  Users of a module need only understand the abstraction provided by its interface.
  The most important issue in designing classes and other modules is to `make them deep`, so that they have simple interfaces for the common use cases, yet still provide significant functionality.
  This maximizes the amount of complexity that is concealed.

## 第五章：信息隐藏(和泄露)（Information Hiding (and Leakage)）

不同模块负责的内容应该正交（互不相关）
这章指导我们正确的进行模块`信息的隐藏（实现）`和`信息暴露（文档）`
在日常开发中，需要仔细了解某个 service 或者方法的实现才敢去用它，这正是因为他们的接口设计的不好，该泄露的没有泄露，该隐藏的没有隐藏，理想情况下我们直接看文档直接使用就行了

- 实现深层模块最重要的技术是`信息隐藏`
  基本思想是每个模块应封装一些`知识`，这些知识代表设计决策。(有点 RDD 的思想了)
  该知识嵌入在模块的实现中，但不会出现在其界面中，因此`其他模块不可见`
  知识可以理解为数据结构和算法(如何实现 TCP 网络协议;如何在多核处理器上调度线程;如何解析 JSON 文档)
  不能把实现的信息暴露到接口。比如具体的数据结构和算法细节

- 信息泄漏
  信息隐藏的反面是信息泄漏
  如果您在类之间遇到信息泄漏，请自问“我如何才能重新组织这些类，使这些特定的知识只影响一个类?”
  解决信息泄漏的两个方法：

  1. union：重新组织类，使特定的知识只影响一个类
  2. extract：从所有受影响的类中提取信息，创建一个只封装这些信息的新类

  当在多个地方使用相同的知识时，例如两个都理解特定类型文件格式的不同类，就会发生信息泄漏。
  **软件设计人员比较重要的技能是确定谁在什么时候需要什么信息！需要暴露的信息就需要让它变得明显！**

- `Temporal decomposition(时间分解)` 导致信息泄漏
  时间分解：按照过程顺序拆分模块(bad)
  `按照角色知识拆分模块，而不是按照过程拆分(good)`
  `设计模块的时候，关注每个需要执行的任务上，而不是他们的执行顺序。`
  比如：读一个文件，修改文件内容，最后写文件。如果把三个操作封装成三个模块，就会把编码算法这个信息分布在读模块和写模块当中。

## 第六章：通用化（General-Purpose Modules are Deeper）

通用化的模块更`容易变得深度以及隐藏信息`

- 做一定的抽象，不仅仅满足当前的需求；同时不过度抽象，当实现当前需求时感到难用
  比如：文本编辑器中插入、删除字符串。比较合适的设计的是

  ```cpp
  void insert(Position s, String s)
  void delete(Position s, Position e)
  ```

  比设计成

  ```cpp
  void insert(Cursor c, String s)
  void delete(Cursor s, Cursor e)
  ```

  要好，因为前者需要的信息是本身固有的，可以用在其他不是 UI 的场景，`更加通用`。而后者只能在 UI 的基础上使用。

- 通用性可以更好地隐藏信息
  当细节很重要时，最好使它们明确且尽可能明显; 当细节不重要时，隐藏它们
- 设计通用化接口的一些启发式问题？

  - 满足当前需求需要最简单的接口是什么？如果`接口数量较少`那很有可能是通用化的。
  - 这个方法有多少个不同的地方使用？专用的方法通常只会在`一个地方使用`。
  - 接口使用起来方便吗？这个可以防止通用化做的太过，如果`使用不方便`有可能是过度通用化了。

- 总结
  通用接口比专用接口具有许多优点。
  通用接口提供了类之间的更清晰的分隔。
  而专用接口则倾向于在类之间泄漏信息。
  `使模块具有某种通用性`是降低整体系统复杂性的最佳方法之一。

## 第七章: `分层抽象`（Different Layers, Different Abstractions）

每一层都关注不同的问题，对应不同的抽象，上层只能调用相邻下层的接口，下层不能调用上层，最典型的案例就是 OSI 分层模型
当分层后发现相邻层存在`只是传递了参数给另一层的调用情况时(Pass-through methods，直通方法)`，说明两层的模块之间职责划分有问题，增加了复杂性但没有增加能力，需要调整。
对于 Dispatcher、Decorator 的设计来说，虽然 interface 重复，但`增加了能力`，所以认为没有问题
存在一些参数需要跨层传递时，推荐直接使用共享 Context 对象.

使用装饰者模式模式的一些启发式问题？

1. 能否`把新功能放在当前类中`？如果功能相对通用，或者当前类的逻辑类似，在或者当前类会使用这个功能，那么就应该放在当前类中。继续吐槽 JDK 中 BufferInputStream。
2. 如果新功能专门为了当前用户场景的，那么放在使用这个接口的那一层比较合适？
3. 能否和现有的一个`装饰者合并`，从而让那个装饰者更加深度？
4. 能否重新实现一个类呢？

## 第八章: 下推复杂逻辑(Pull Complexity Downwards)

- 目的是`简化接口`：对于一些配置参数，如果模块内部自己决策能比使用该模块的用户决策得更好，那就内部消化掉，不要暴露给用户去决策。

- 如果一个复杂度没法避免，尽量放在实现部分，可以简化接口，降低复杂度。放在接口部分会导致所有的用户都需要感知这个复杂度，增加了整体复杂度。比如：抛出一个异常就是把处理异常的交给了用户，任何使用者都必须感知。

## 第九章: 合并还是分开？

给定两个功能，它们应该在同一位置一起实现，还是应该分开实现？
这章讲了什么情况下模块应该合并，什么情况下应该拆分。
合并和分开模块的目的是为了最少的依赖，最好的`信息隐藏以及深度接口`

- 分开导致的一些可能问题：

  1. 代码重复
  2. 需要额外管理这些分开的模块的代码
  3. 分开以后不容易寻找相关的模块

- 如果两块代码是相关就应该放在一起，一些相关的暗示：
  - 共享信息(例如最好在同一位置读取和解析请求)
  - 一起使用；使用了 A 就一定会使用 B，使用了 B 就一定会使用 A。
  - 重复概念
  - 如果不看一个模块另外一个模块就无法理解
- 方法代码行数问题
  长度其实不重要关键是要`把抽象做好`，一个方法做一个事情而要做完整。
  如果抽象做的清晰，自然方法的行数也不会太多，因为可以把一个抽象再做细分么。

- 长方法并不总是坏的
  例如，假设一个方法包含按顺序执行的五个 20 行代码块。如果这些块是`相对独立`的，则可以一次读取并理解该方法的一个块。将每个块移动到单独的方法中并没有太大的好处。如果这些块具有复杂的交互作用，则将它们`保持在一起`就显得尤为重要，这样读者就可以一次看到所有代码。如果每个块使用单独的方法，则读者将不得不在这些扩展方法之间`来回切换`，以了解它们如何协同工作。如果方法具有简单的签名并且易于阅读，则包含数百行代码的方法就可以了。这些方法很深入（很多功能，简单的接口），很好。
- 拆分的方法
  从单一职责原则、解耦、信息隐藏等角度考虑，拆分是有必要的。
  1. 提取子任务
  2. 拆分为两个单独的方法 -> 每个结果方法的接口应该比原始方法的接口更简单

## 第十章: 错误处理

异常处理所带来的复杂度很大，所以需要**减少不必要的异常**

- 可能出现过多的场景：
  - 过度防御编程。
  - 逃避处理复杂的场景，直接把错误扔给了使用者处理。
  - 接口中过多的异常信息。
- 减少异常的方式
  - 把异常当作正常逻辑的一部分，简化接口。例如删除一个并不存在的文件；
  - 隐藏异常，如果在发生异常情况的较低层次能处理好，就不要对高层抛出异常；比如 TCP 超时重传；
  - 合并异常，抽象一些异常，把多种异常用统一的逻辑来处理。spring 中 ResponseEntityExceptionHandler 就是典型应用
  - 直接 crash，这个只使用于一个错误无法处理，或者处理的复杂度太高的情况。`打印错误日志，保留现场，停止应用程序`(logFatal)

## 第十一章: 设计两次

无论成本多高，或者多难以实现，都记得多设计几个方案，横向对比和思考一下每个方案的优缺点，或者你就能找到当前场景下最佳的方案。

## 第十二章: 写好注释

一定要写注释，所谓的：好的代码是自注释的，这个观点是扯淡的，大家不要相信，再好的代码也得有注释
The overall idea behind comments is to capture information that was in the mind of the designer `but couldn’t be represented in the code`.
**注释应该表达代码无法表达的信息** -> 一个技巧尽量不用代码中使用的单词
一般需要包含一个"高"和一个"低"：高是 summary，低是 detail

- 接口的注释需要和实现的注释分开，如果接口注释还需要写实现的注释，那么这个接口是比较浅的

- 接口的注释
  高层的注释、使用者视角下的行为、调用的前置条件、副作用
- 实现的注释
  做什么以及为什么

## 第十三章: 注释应该描述代码中不明显的内容

编写注释的原因是，使用编程语言编写的语句`无法捕获编写代码时开发人员想到的所有重要信息`。
注释记录了这些信息，以便后来的开发人员可以轻松地理解和修改代码。注释的指导原则是，注释应描述代码中不明显的内容。

- 注释的最重要原因之一是抽象
  注释可以提供一个更简单，更高级的视图（“调用此方法后，网络流量将被限制为每秒 maxBandwidth 字节”）
- 不要重复注释

## 第十四章: 命名

- 当不知道怎么命名的时候，往往会意味着代码的设计有一些问题，需要重构了！
- 好的命名三个要求：反映`本质`、`精确`以及`一致`

## 第十五章: 先写注释

## 第十六章: 修改老代码

## 第十七章: 保持一致性(Consistency)

维护一致性不容易，有几个方法

1. 文档
2. 工具自动化
3. 入乡随俗

## 第十八章: 让代码的意图更加明显

几个可能导致代码不清晰的地方：

1. 异步编程:主要是控制流不清晰，请求和返回不连贯，需要额外的信息去串起来。主要通过注释来解决
2. 泛型容器；当用来做一个返回的类型时，它提供的方法很难准确表达具体的业务信息
3. 定义和声明类型不同
4. 代码行为超出阅读者预期

## 第十九章: 软件开发趋势和软件复杂度

每当您遇到有关新软件开发范例的提案时，就必须从复杂性的角度对其进行挑战：`该提案确实有助于最大程度地降低大型软件系统的复杂性吗？`从表面上看，许多建议听起来不错，但是如果您深入研究，您会发现其中一些会使复杂性恶化，而不是更好。

## 第二十章: 为性能而设计

- 干净的设计和高性能是兼容的(golang 后缀数组实现)
- 复杂的代码通常会很慢，因为它会执行多余或多余的工作
- `找到对性能最重要的关键路径并使它们尽可能简单`

## 第二十一章: 结论

这本书是关于一件事的：复杂性。

## 总结

- Summary of Design Principles 设计原则摘要

1. Complexity is incremental: you have to sweat the small stuff (see p. 11).
   软件复杂度是慢慢积累起来的，你必须锱铢必较
2. Working code isn’t enough (see p. 14).
   能工作的代码还不够
3. Make continual small investments to improve system design (see p. 15).
   持续进行少量投资以改善系统设计
4. Modules should be deep (see p. 22)
   模块应该是深的
5. Interfaces should be designed to make the most common usage as simple as possible (see p. 27).
   **接口设计应当使得最常用的路径越简单越好**
6. It’s more important for a module to have a simple interface than a simple implementation (see pp. 3. 55, 71).
   **一个模块具有一个简单的接口比一个简单的实现更重要**
7. General-purpose modules are deeper (see p. 39).
   模块应当尽量设计的通用
8. Separate general-purpose and special-purpose code (see p. 62).
   分离通用的代码和特定需求的代码(往小了说就是抽取公共函数或者公共类，往大了说就是建立共享服务或者中台系统)
9. Different layers should have different abstractions (see p. 45).
   不同的层应具有不同的抽象(如果不同层次的抽象是一样的，那必然会导致大量的重复代码和没有意义的类)
10. Pull complexity downward (see p. 55).
    把复杂性放在底层
11. Define errors (and special cases) out of existence (see p. 79).
    通过定义让错误（和特殊场景）没有机会发生
12. Design it twice (see p. 91).
13. Comments should describe things that are not obvious from the code (see p. 101).
    注释应描述代码中不明显的内容
14. Software should be designed for ease of reading, not ease of writing (see p. 149).
    软件的设计应易于阅读而不是易于编写
15. The increments of software development should be abstractions, not features (see p. 154).
    软件开发的增量应该是抽象而不是功能

- Summary of Red Flags 危险信号

1. Shallow Module: the interface for a class or method isn’t much simpler than its implementation (see pp. 25, 110).
   浅模块：类或方法的接口并不比其实现简单得多
2. Information Leakage: a design decision is reflected in multiple modules (see p. 31).
   信息泄漏：设计决策反映在多个模块中
3. Temporal Decomposition: the code structure is based on the order in which operations are executed, not on information hiding (see p. 32).
   时间分解：代码结构基于执行操作的顺序，而不是信息隐藏
4. Overexposure: An API forces callers to be aware of rarely used features in order to use commonly used features (see p. 36).
   过度暴露：为了使用常用功能，API 强制调用者注意很少使用的功能，
5. Pass-Through Method: a method does almost nothing except pass its arguments to another method with a similar signature (see p. 46).
   Pass-Through Method：一种方法几乎不执行任何操作，只是将其参数传递给具有相似签名的另一种方法
6. Repetition: a nontrivial piece of code is repeated over and over (see p. 62).
   重复：一遍又一遍的重复代码
7. Special-General Mixture: special-purpose code is not cleanly separated from general purpose code (see p. 65).
   特殊通用混合物：特殊用途代码未与通用代码完全分开
8. Conjoined Methods: two methods have so many dependencies that its hard to understand the implementation of one without understanding the implementation of the other (see p. 72).
   联合方法：两种方法之间的依赖性很大，以至于很难理解一种方法的实现而又不理解另一种方法的实现
9. Comment Repeats Code: all of the information in a comment is immediately obvious from the code next to the comment (see p. 104).
   注释重复代码：注释旁边的代码会立即显示注释中的所有信息
10. Implementation Documentation Contaminates Interface: an interface comment describes implementation details not needed by users of the thing 1. being documented (see p. 114).
    实施文档污染了界面：界面注释描述了所记录事物的用户不需要的实施细节
11. Vague Name: the name of a variable or method is so imprecise that it doesn’t convey much useful information (see p. 123).
12. Hard to Pick Name: it is difficult to come up with a precise and intuitive name for an entity (see p. 125).
    难以选择的名称：很难为实体提供准确而直观的名称
13. Hard to Describe: in order to be complete, the documentation for a variable or method must be long. (see p. 131).
    难以描述：为了完整起见，变量或方法的文档必须很长
14. Nonobvious Code: the behavior or meaning of a piece of code cannot be understood easily. (see p. 148).
    非显而易见的代码：一段代码的行为或含义不容易理解

---

读完这本书，收获更多的不是具体 design 的方法，而是 mindset。

现在回去看本书的封面觉得很有意思，上方是杂乱的线，代表复杂的 implementation，而下方是整齐的线，代表 interface；这表示本书最重要的一个观点，module should be deep，一个 deep module 有简单的 interface，但是有复杂的 implementation，隐藏了很多 module 使用者不需要的信息。deep module 有助于减少 complexity，对于 module 的使用者来说，需要了解得更少，但是能获得更多的。

另一个收获就是 investment mindset。我工作中也常常会以不知道如何 design，觉得 design 太花时间，改动太大可能会破坏以往已经实现的功能为由，避免 design，停留在 feature-driven development 的舒适区。但长时间地注重实现功能的开发（tactical programming）会导致 complexity 的叠加，此时需要修改一个很小的部分，也常常需要花很多的时间。而 investment mindset 的想法是可以花 20%左右的时间去思考这些问题：为什么要实现这个功能；怎样实现这个 feature 才能让它好像一开始就在 design 中一样；如何用别的办法实现这个功能，和原方法相比有什么利弊；怎样才能让代码的读者更容易理解我的代码；有哪些代码是冗余的。可能我刚开始需要花更多的时间，但我相信长期的训练可以有益于自己对更复杂软件的实现。

investment mindset 是一种视角的转变，之前我 evaluate 我工作的完成可能是以单一的时间维度：如果我能在很快的时间让代码工作起来，就是好的；但是现在加入了 complexity 的维度：我要更多得考虑 design，考虑减少系统的复杂度，让之后 maintain 起来更方便，也要提升自己对 system abstraction 的理解。

investment mindset 也是一种长期主义的心态，我应该做一些长期有益的事情。做项目并不是很短暂的过程，所以需要写 unit test 方便之后项目的重构，写 comment & document 方便之后修改的时候参考。同时生活中也是一样，我没必要纠结于一时的得失因为它会过去；我不能逃避面对问题因为之后同样的问题可能会再出现；我不用焦虑于当下自己做不好一些事情，因为只要我 keep practice，总有一天自然而然就能做好。虽然好像是一些大道理，但是花一些时间想清楚，接下来只要去实践。

总之，这本书相比别的具体介绍 design pattern 之类的设计书而言，并不是非常实用，但它用一整本书讲了为什么我要 design，究竟有什么好处；我没有看过其他 design 类的书籍，但现在也有兴趣去探索一下相关的书籍。而且全书的逻辑非常清晰，例子很多，读起来也不吃力，是一本值得推荐的好书。

---

软件复杂度是日积月累的，如果只是打补丁总有一天会失控。
可以容忍一开始不好的设计，但是之后每次都可以试着去优化一点。
但也不要想动不动就整个大新闻，毕竟业务的稳定和投入产出也是要考虑的。
你必须一直去思考这个问题，如果一切只会完成短期任务，那么复杂度一定会失控。

`模块和类要有深度`。
我觉得这个问题在 Java 的框架里特别明显，一层套一层，理论上是面向对象，可复用可替换，但是对于 99%的应用来说，这个需求根本不存在。
一个模块或类，甚至是函数，要解决一个完整的问题。`能原子化的就不要再为了拆分而拆分`。

`信息隐藏`。在接口层面要暴露的是做什么，至于怎么做留到实现层面。复杂度尽可能下沉，接口的调用者不要去想，也不应该去操心底层实现的细节。

尽量为通用的场景提供`默认的行为`，不要让用户自己去理解。

错误处理不要过度设计。
不要把错误扔给用户处理，很多时候用户也不知道该怎么处理。
如果`索引范围不符合返回一个空字符串`就能解决大部分问题，非得扔个 IndexOutOfBoundsException 异常出来。

写注释，注释一定要有用，要把那些不容易理解的东西解释清楚，而不是翻译代码。

`一致性`，命名，代码风格，处理逻辑尽可能统一。哪怕是低水平的统一也好过各自放飞。

测试驱动开发是伪命题。
`大多数人根本没有先写用例，再写代码的能力`。
另外用例不可能一次写对，但用例摆在那里，开发的时候代码就会已通过测试为目标，很多时候容易导致拼凑补丁。
但是有一种情况必须先写用例，就是处理历史遗留屎山。
当前代码支持的业务场景，以及你对代码功能的所有猜测必须有用例作为支撑，这样在修改完成后可以快速检查有没有影响现有功能

最重要的原则，`把自己的第一个想法毙掉`。至少要想两个方案，有选择才有比较

---

软件设计
