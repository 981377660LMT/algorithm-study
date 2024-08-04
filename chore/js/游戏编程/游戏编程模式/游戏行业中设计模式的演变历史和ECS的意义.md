如何评价《守望先锋》架构设计？ - 廖鑫炜的回答 - 知乎
https://www.zhihu.com/question/61169850/answer/185652605

@猴与花果山童鞋已经阐述了 ECS 的主要概念。此文主要从技术和工程角度简单探讨游戏行业中设计模式的演变历史和 ECS 的意义。

1. 综述

设计模式产生的动机，往往源于尝试解决一些存在的问题。游戏开发领域技术架构的宏观目标，大体包括以下目标：

- 适合快速迭代。无论是上线前的敏捷开发流程还是上线后根据市场反馈的调整，迭代都是基本需求。
- 易于保证产品质量。良好的架构，可以降低 Bug 和 Crash 出现概率。
- 开发效率。重要性不必多说。即使是更重视游戏质量的公司，越高的开发效率也有助于更短的时间打造出质量更高的游戏。
- 运行效率。大部分游戏对实时响应和运行流畅度都有很高的要求，同时游戏中又存在大量吃性能的模块（比如渲染、物理、AI 等）。
- 协作扩展性。能够在开发团队扩张时尽可能无痛，同时方便支持美术、策划、音效等非程序同事的开发需求。

现代 Entity Component System 的概念，以及对游戏开发领域的意义：

- Entity：代表游戏中的实体，是 Component 的容器。本身并无数据和逻辑。
- Component：代表实体“有什么”，一个或多个 Component 组成了游戏中的逻辑实体。只有数据，不涉及逻辑。
- System：对 Component 集中进行逻辑操作的部分。一个 System 可以操作一类或多类 Component。同一个 Component 在不同的 System 中，可以关联不同的逻辑。

ECS 并非《守望先锋》所独有和原创，事实上近年来以 ECS 为基础架构逐渐成为国际游戏开发领域的主流趋势。

采用 ECS 的范式进行开发，思路上跟传统的开发模式有较大的差别：

- Entity 是个抽象的概念，并不直接映射为具体的事物：比如可以不存在 Player 类，而是由多个相关 Component 所组成的 Entity 代表了 Player。如 Entity { PositionComponent, RenderComponent, StateMachineComponent, ... } 等等。
- 行为通过对 Component 实施操作来表达。比如简单重力系统的实现，可以遍历所有 Position Component 实施位移，而不是遍历所有 玩家、怪物、场景物件，或者它们统一的基类。
- 剥离数据和行为，数据存储于 Component 中，而 Component 的相关行为，和涉及多个 Component 的交互和耦合，由 System 进行实施。

ECS 框架，至少有以下优点：

- 模式简单。如果还是觉得复杂，推荐看看 GoF 的《设计模式》。
- 概念统一。不再需要庞大臃肿的 OOP 继承体系和大量中间抽象，有助于迅速把握系统全貌。同时，统一的概念也有利于实现**数据驱动**（后面会提到）。
- 结构清晰。Component 即数据，System 即行为。`Component 扁平的表达有助于实现 Component 间的正交。`而封装数据和行为的做法，不仔细设计就会导致 Component 越来越臃肿。
- 容易组合，高度复用。Component 具有高度可插拔、可复用的特性。而 System 主要关心的是 Component 而不是 Entity，`通过 Component 来组装新的 Entity，对 System 来说是无痛的。`
- `扩展性强。增加 Component 和 System，不需要对原有代码框架进行改动。`
- 利于实现面向数据编程（DOP）。对于游戏开发领域来说，`面向数据编程是个很重要的思路。天然亲和数据驱动的开发模式，有助于实现以编辑器为核心的工作流程。`
- 性能更好优化。接上条，相比 OOP 来说，DOP 有更大的性能优化空间。（详见后面章节）

若要了解为何会出现 ECS 这样的模式，以及它所试图解决的问题，需要考虑一下历史进程：

2. 演化路径

## 简单粗暴的上个世纪开发模式

注重于实现相关算法、功能和逻辑，代码只要能实现功能就行，怎么直观怎么来。比如

```cpp
class Player {
    int hp;
    Model* model;
    void move();
    void attack();
};
```

类似这样完全没有或很少架构设计的代码，在项目规模增大后，很快变得臃肿、难以扩展和维护。

## OOP 设计模式的泛滥 案例：OGRE

> 设计模式是语言表达能力不足的产物。 —— 某程序员

那么，作为他山之石，GoF 基于 Java 提出的设计模式，能否有效解决游戏开发领域的问题？

大家还记得当年国内风靡一时的游戏引擎 OGRE 么？

OGRE 总有那么些`学院派的味道，试图通过设计模式的广泛使用，来提高代码的可维护性和可扩展性。`

然而，个人对游戏开发领域大规模使用 OOP 设计模式的看法：

- 设计模式的六大原则大部分仍值得遵循。
- 基于 Java 实现的设计模式，未必适合其它语言和领域。想想 C# 的 event、delegate、lambda 可以简化或者消除多少种 GoF 的模式，再想想 Golang 的隐式接口。
- C++ 是游戏开发领域最主要的语言，可以 OOP 但并不那么 OO，比如缺少语言层面纯粹的 interface，也缺少 GC、反射等特性。照抄 Java 的设计模式未免有些东施尿频，而且难以实现 C++ 所推崇的零代价抽象。（template 笑而不语）
- 局部使用 OOP 设计模式来实现模块，并暴露简单接口，是可以起到提升代码质量和逼格的效果。然而在架构层面滥用，往往只是把逻辑中的复杂度转移到架构复杂度上。
- `滥用设计模式导致的复杂架构，并不对可读性和可维护性有帮助。`比如原本 c style 只要一个文件顺序读下来就能了解清楚的模块，滥用设计模式的 OOP 实现，阅读代码时有可能需要在十几个文件中来回跳转，还需要人脑去正确保证阅读代码的上下文...
- 过多的抽象导致过多的中间层次，却只是把耦合一层一层传递。直到最后结合反射 + IoC 框架 + 数据驱动，才算有了靠谱的解决方案。然而一提到反射，C++表示我的蛋蛋有点疼。

那么，有没有办法简化和沉淀出游戏开发领域较通用的模式？

## 未脱离 OO 思想的 Entity Component 模式 案例：Unity3D

Unity3D 是个使用了 Entity Component 模式的成功的商业引擎。

相信使用过 Unity3D 的童鞋，都知道 Unity3D 的 Entity Component 模式是怎么回事。（在 Unity3D 中，Entity 叫 GameObject）。

其优点：

- `组件复用`。体现了 ECS 的基本思想之一，Entity 由 Component 组成，而不是具体逻辑对象。设计得好的 Component 是可以高度复用的。
- `数据驱动`。场景创建、游戏实体创建，主要源于数据而不是硬编码。以此为基础，引擎实现了以编辑器为中心的开发模式。
- `编辑器为中心`。用户可在编辑器中可视化地编辑和配置 Entity 和 Component 的关系，修改属性和配置数据。在有成熟 Component 集合的情况下，新的关卡和玩法的开发，都可以完全不需要改动代码，由策划通过编辑器实现。

看起来，Unity3D 已经在很大程度上解决了游戏设计领域通用模式的问题。然而，其 Entity Component 模式仍然存在一些问题：Component 仍然延续了一些 OOP 的思路。比如：

- Component 是数据和行为的封装。虽然此概念容易导致的问题可以通过其它方式避免，但以不加思考照着最自然的方式去做，往往会造成 Component 后期的膨胀。比如 Component 需要支持不同的行为就定义了不同的函数和相关变量；Component 之间有互相依赖的话逻辑该写在哪个 Component 中；多个 Component 逻辑上互相依赖之后，就难以实现单个 Component 级别的复用，最后的引用链有可能都涉及了代码库中大部分 Component 等等。
- Component 是支持多态的引用语义。这意味着单个 Component 需要单独在堆上分配，难以实现下文所提到的，对同类型多个 Component 进行数据局部性友好的存储方式。这样的存储方式好处在于，批量处理可以减少 cache miss 和内存换页的情况。

## 当前主流的 Entity Component System 架构 案例：EntityX

那么，综合以上所说的各种问题，一个基于 C++ 的现代 Entity Component System，应该是什么样子？

具体案例，可以参考 [EntityX](https://github.com/alecthomas/entityx)，一个开源的 C++ ECS 框架。

**一一实现了前述现代 ECS 的各种概念：Entity 只是个 ID，Component 存储数据，System 实现关联多个 Component 的行为。**

代码味道：

```cpp
struct Position {
  Position(float x = 0.0f, float y = 0.0f) : x(x), y(y) {}

  float x, y;
};

struct Direction {
  Direction(float x = 0.0f, float y = 0.0f) : x(x), y(y) {}

  float x, y;
};

struct MovementSystem : public System<MovementSystem> {
  void update(entityx::EntityManager &es, entityx::EventManager &events, TimeDelta dt) override {
    es.each<Position, Direction>([dt](Entity entity, Position &position, Direction &direction) {
      position.x += direction.x * dt;
      position.y += direction.y * dt;
    });
  };
};
```

如上，实现了两类 Component：Position 和 Direction。

MovementSystem 只关心同时具有两类 Component 的 Entity。

一些值得说的特点：

- 低抽象代价。C++ 的模板特性，便于把不少在其他语言中难以避免的运行时开销，转移到编译时。
- 同类的多个 Component 实现了紧凑连续的内存布局。这个特性为什么重要？请参考[这个问题](https://www.zhihu.com/question/20275578) @Milo Yip 的回答。同时这也是 Unity3D 的 Entity Component 模式难以做到的。当遍历同类 Component 时，数据存储于连续的内存空间中，可以大大提高缓存命中率。
- Component 只有数据，行为是 System 的事。这样的模式，避免了上一节提到的 Unity3D 中容易出现的问题。Component 没有逻辑上的互相引用，Component 的耦合和依赖由 System 处理。此外，由 System 进行统一的状态修改，也有利于定位和隔离问题。
- System 间的解耦，主要通过事件回调。`System 之间不提倡互相引用`，通过 Signal 来实现 publish / subscribe 进行处理。《守望先锋》也提到了关于 System 间发生了耦合的麻烦情况通常用 Singleton 模式和把`共用代码放进 Utils 解决`。

## ECS 的进一步优化

除了可以提高缓存命中率外，新世代的 ECS 还可以通过`分析数据依赖和读写关系，来实现 System 间的并行`。比如更新时， System A 需要读 组件 1，System B 需要读 组件 1、写组件 2，System C 需要写 组件 1，那么调度时可以把 System A 和 System B 分配到不同线程处理，之后再处理 System C。原贴中也一笔带过提到了这方面的优化。然而对于复杂的 C++ 游戏来说，这个目标在实践上的可行性具有比较大的障碍：难以确保团队中的熊孩子不小心写出非线程安全的代码。

不过，Rust 给这个问题带来了解决方案。可以参考 Rust 实现的并行 ECS 框架：slide-rs/specs

Rust 的语言特性在编译期保证了线程安全，只需声明一下 System 对 Component 的访问权限如：

```rust
    type SystemData = (ReadStorage<'a, Velocity>,
                       WriteStorage<'a, Position>);
```

这样，即可安全地获得多线程带来的性能提升。

---

1. 守望团队在将整个游戏转成 ECS 之前也不确定 ECS 是不是真的好使。现在他们说 ECS 可以管理快速增长的代码复杂性，也是事后诸葛亮。
2. 如果可以把整个游戏世界都抽象成数据，存档/读档功能的实现也变得容易了。存档时只需要将所有 Component 数据保存下来，读档时只需要将所有 Component 数据加载进来，然后 System 照常运行。想想就觉得强大，这就是 DOP 的魅力。
