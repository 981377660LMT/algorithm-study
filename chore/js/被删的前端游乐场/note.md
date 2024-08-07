被删的前端游乐场

https://godbasin.github.io/categories/
https://godbasin.github.io/front-end-playground/front-end-basic/front-end/front-end-1.html

# 前端 CHANGELOG 生成指南

https://godbasin.github.io/2019/11/10/change-log/
https://zhuanlan.zhihu.com/p/51894196

conventional-changelog
主要介绍 commitizen、conventional-changelog-cli 、standard-version 这三工具了。

# 数据驱动 vs 事件驱动

事件驱动

构建页面：设计 DOM => 生成 DOM => 绑定事件
监听事件：操作 UI => 触发事件 => 响应处理 => 更新 UI

数据驱动

构建页面：设计数据结构 => 事件绑定逻辑 => 生成 DOM
监听事件：操作 UI => 触发事件 => 响应处理 => 更新数据 => 更新 UI
在使用数据驱动的时候，模板渲染的事情会交给框架去完成，我们需要做的就是数据处理而已。

# 前端模板引擎流程

https://godbasin.github.io/2017/10/21/template-engine/
数据绑定已经作为一个框架最基础的功能。我们常常使用的单向绑定、双向绑定、事件绑定、样式绑定等，里面具体怎么实现，而当我们数据变动的时候又会触发怎样的底部流程呢

`模板生成 AST => AST 生成模板 => 数据/事件/属性绑定的监听 => 数据变更 Diff => 局部更新模板`

1. 模板数据绑定
   1.1. 解析语法生成 AST
   1.1.1. 捕获特定语法
   1.1.2. DOM 元素捕获
   1.1.3. 数据绑定捕获
   1.2. AST 生成模版
   1.2.1. 生成模版的方法
   1.2.2. 浏览器的渲染机制
2. 模版数据更新
   2.0.1. 数据更新监听
   脏监测/getter-setter/手动更新
   脏检测是 Angular 的一大特色。由于事件触发的时候，并不能知道哪些数据会有变化，所以会进行大面积数据的新旧值 Diff，这也毫无疑问会导致一些性能问题。在 Angular2 版本之后，由于使用了`zone.js对异步任务进行跟踪，把这个计算放进worker，完了更新回主线程`，是个类似多线程的设计，也提升了性能。

   2.0.2. 数据更新 Diff

# 页面区块化与应用组件化

1. 页面区块化
   - 1.1. 什么是区块
2. 应用组件化

   - 2.1. 什么是组件
     组件可以扩展 HTML 元素，封装可重用的代码。

   - 2.2. 组件的划分

     1. 视觉和交互上是一个完整的组件。
     2. 写代码的时候，可重复的内容即可视为一个组件。

     在一个团队内，最好是使用一种方式来进行划分。因为对于成员的相互配合和项目的维护来说，统一的规范是比较重要的。

   - 2.3. 组件的封装
     - 2.3.1. 组件内维护自身的数据和状态
     - 2.3.2. 组件内维护自身的事件
     - 2.3.3. 通过初始化事件，来初始化组件状态，激活组件
     - 2.3.4. 对外提供配置项，来控制展示以及具体功能
     - 2.3.5. 通过对外提供查询接口，可获取组件状态
   - 2.4. 组件的通信
   - 2.5. 组件化程度
     组件的封装是会消耗维护性的，过度的封装会导致代码难维护，可读性也差。所以我们需要根据项目的大小以及复杂度，来进行什么程度的封装
     适度也是很重要的，个人认为`好的架构是变化的，跟随着项目变化而变化，保持着拓展性和维护性`。如果说我们只是为了抽象而抽象，那想必会把简单的事情复杂化，整个应用和代码会变得难以理解。适度的抽象很重要，但相比错误的抽象过程，没有抽象可能会更好。

# 一个组件的自我修养

1. 组件的划分
   1.1. 通过视觉和交互划分
   1.2. 通过代码复用划分
2. 组件的封装
   2.1. 独立的组件 -> 维护自身数据、事件
   2.2. 组件与外界 -> 对外提供钩子、配置项、查询接口(获取组件状态)

# 组件配置化

1. 配置化思想

   - 1.1. 可配置的数据
     - 1.1.1. 应用中的可配置数据
       搭起一整套的运营管理平台，一些简单的文字或是数据，则可以通过平台进行配置
     - 1.1.2. 代码中的可配置数据
       常量配置
     - 1.1.3. 文件里的可配置数据
       把这样的可配置数据，单独写到某个文件里维护。当需要调整的时候，只需要下发一个配置文件就好啦
   - 1.2. 可配置的接口
     配置化的实现有两点很重要的东西：`规范`和`解决方案`。如果说目前较好的从前端到后台的规范，可能 GrapgQL 和 Restful ，大家不熟悉的也可以去看看啦。

   - 1.3. 可配置的页面

     这种页面的配置，基本上有两种实现方式：

     - 配置后`生成静态页面`的代码，直接加载生成的页面代码。
       适用于一些移动端的模版页面开发，例如简单的活动页面、商城页面等等
     - 写通用的配置化逻辑，在`加载页面的时候拉取配置数据，动态生成页面`。
       一些管理平台的实现，毕竟大多数都是增删查改，形式无非列表、表单和菜单等

   - 1.4. 可配置的应用

2. 组件配置化

   - 2.1. 可配置的数据
   - 2.2. 可配置的样式

     - 2.2.1. 根据子元素配置 CSS （缺点：`DOM 结构调整影响大`）
     - 2.2.2. 根据子 class 配置 CSS（BEM，`缺点：名字太长`）
       block：可以与组件和模块对应的命名，如 card、dialog 等
       element：元素，如 header、footer 等
       modifier：修饰符，可视作状态等描述，如 actived、closed 等

       ```less
       .my-dialog {
        background: white;
        &__header {}
        &__section {}
        &__footer {}
        &__btn {
          &--inactived
        }
       }
       ```

       当然，如今很多框架都支持样式的作用域，通常是通过在 class 里添加随机 MD5 等，来保持局部作用域的 class 样式。常见的话，我们是搭配第一和第二种方式一起使用的。

   - 2.3. 可配置的展示
   - 2.4. 可配置的功能
     灰度？

配置化思想：**把相似的部分提取出来抽象封装，把可变的部分结合配置来高效地调整（interface、config、event）**

# 数据抽离与数据管理

从事件驱动脱离，来到了数据驱动的世界
在把数据与逻辑分离到极致的时候，你再看一个应用，会看到一具静态的逻辑躯壳，以及数据如灵魂般地注入到应用里，使其获得生命。

**数据的抽离，其实与配置化的思想有相通的地方**，即把可变部分分离，然后通过注入的方式，来实现具体的功能和展示。

1. 应用数据抽离
   - 1.1. 状态数据(state)
     怎么定义状态数据？最浅显的办法就是`这些数据，可以直接影响模块的状态`，如对话框的出现、隐藏，标签的激活、失活，长流程的当前步骤等
   - 1.2. 动态数据(props)
     不会跟随着应用的生命周期而改变，也不会随着应用的关闭而消失。它们独立存在于外界，通过注入的方式进入应用，并影响具体的展示和功能逻辑
     例如从数据库拉取回来的
   - 1.3. 将数据与应用抽离
     办公室上班例子，人涌入公司，给公司注入灵魂，公司得以运作。
     我们要做的，不只是如何划分数据、将数据与应用`抽离`，我们还需要将其有规律地`管理`。
2. 应用数据管理

   - 2.1. 数据的流动

     - 2.1.1. 事件通知(event)
       事件通知机制很方便，可以随意地定义触发的时机，也可以任意地点使用监听或是触发。
       但事件机制的弊端也是很明显，就是每一个事件的触发对应一个监听，关系是一一对应。在整个应用中看，则是`散落在各处，随意乱窜的数据流动`。需要定位的时候，只能通过`全局搜索`的方式来跟踪数据的去向。
       当然，也有些人会定义一个**中转站，所有的事件数据流都会经过那**，这样的维护方式会有所改善。
     - 2.1.2. 共享对象(store) 通过注入对象的引用，来在不同组件中获取相同的数据源
     - 2.1.3. 单方向流动(setState) Vuex、Redux
     - 2.1.4. 树状作用域(context) golang 的 context

   - 2.2. 适度的管理
     一个状态管理工具则可以轻松解决乱糟糟的数据流问题

# 前端构建大型应用(2018 年文)

1. 项目设计
   定位: to B/C，大小，框架和工具的选型，项目和团队规范

   - 1.1. 框架选择
     - 1.1.1. Angular
       这里的 Angular 是指 Angular 2.0+ 版本，v1.0 我们通常称之为 AngularJS，目前已经不更新了
       项目中使用 Angular，最大的体验感受则是项目`有完备的结构和规范`，新加入的成员能很快地通过复制粘贴完成功能的开发。身边有- 人说过，`好的架构设计，能让高级程序员和初入门的程序员写出相似的代码`，这样对于整体管理和项目的维护有非常好的体验。
       很多人说 Angular 难上手，其实主要在于开始的项目搭建、以及 Angular 独有的一套设计方案的理解。
       但是依赖注入的设计方式，我们`几乎不用考虑很多数据和状态管理的问题`。当然脏检查的方式曾经也带来性能问题，后面在加入树状的模块化、Zone.js 之后，即使没有虚拟 DOM，性能也是有大大的提升。
       > https://godbasin.github.io/2021/05/01/angular-design-zonejs/ zonejs
     - 1.1.2. React
     - 1.1.3. Vue
     - 1.1.4. 开源框架？
   - 1.2. 项目代码结构
     就是公共组件、工具等同类的文件，放置一起维护会比较好。而且还有个小 tips，我们可以在搭建项目的时候，在 README.md 里面描述下该项目下的代码和文件结构
   - 1.3. 代码流程规范

2. 大型应用优化

   - 2.1. 路由管理
   - 2.2. 抽象和组件化
     在我们开始写重复的代码、或是进行较多的复制粘贴的时候，大概我们需要考虑对组件进行适当的抽象了。
     > 不过 copilot 用多了，好像注意不到？？？
   - 2.3. 状态和数据管理
   - 2.4. 代码打包
     - 2.4.1. 路由异步加载
     - 2.4.2. Webpack 分块打包
     - 2.4.3. Source map
     - 2.4.4. Tree-shaking
   - 2.5. 编写可测试代码
     项目中功能的快速迭代、开发工作量饱满等原因，导致甚至单元测试这种都很少编写。Emmmmm。。。

# Angular 框架解读--Zone 区域之 zone.js

https://blog.csdn.net/kingslave1/article/details/135630112

提供了一种跟踪和管理异步操作的机制。它的核心概念是 Zone，它可以帮助我们捕获和处理异步操作的上下文
Zone 在 Angular 中有很多用途

1. 变更检测
   Angular 的变更检测机制是依赖于 zone.js 的
   每当发生异步操作时，zone.js 会通知 Angular 进行变更检测，以确保视图能够及时更新
2. 错误处理
   捕获和处理异步操作中的错误
3. 性能监控
   监控异步操作的执行时间，以便评估和优化应用程序的性能

- NgZone：Angular 中的 zone.js
  当我们在 Angular 应用程序中执行异步操作时，NgZone 会自动创建一个 Zone，并把这些操作放入该 Zone 中。这样做的好处是，我们可以在异步操作完成后触发变更检测，以确保视图能够及时更新
  通过使用 NgZone 的 run()方法，我们确保异步任务的结束能够触发变更检测

## 原理

https://godbasin.github.io/2021/05/01/angular-design-zonejs/
https://juejin.cn/post/6859348400463314951
https://segmentfault.com/a/1190000044163634
https://github.com/JLQmiller/angularindepth/blob/master/articles/angular-35.%5B%E7%BF%BB%E8%AF%91%5D-%E7%BF%BB%E9%98%85%E6%BA%90%E7%A0%81%E5%90%8E%EF%BC%8C%E6%88%91%E7%BB%88%E4%BA%8E%E7%90%86%E8%A7%A3%E4%BA%86Zone.js.md
https://ithelp.ithome.com.tw/m/articles/10220772

context???

zone.js 接管了浏览器提供的异步 API，比如点击事件、计时器等等。也正是因为这样，它才能够对异步操作有更强的控制介入能力，提供更多的能力。
为每个点击函数安排了一个事件任务

先了解 run 与 runOutsideAngular 两个 API 即可

# 数据库事务

- 所谓事务是用户定义的一个数据库操作序列，这些操作要么全做要么全不做，是一个不可分割的工作单位。
  当事务被提交给了 DBMS（数据库管理系统），则 DBMS（数据库管理系统）需要确保该事务中的所有操作都成功完成且其结果被永久保存在数据库中，如果事务中有的操作没有成功完成，则事务中的所有操作都需要被回滚，回到事务执行前的状态;同时，该事务对数据库或者其他事务的执行无影响，所有的事务都好像在独立的运行。

- 数据库事务拥有以下四个特性，习惯上被称之为 ACID 特性。

  - 原子性（Atomicity）：包含在其中的对数据库的操作要么全部被执行，要么都不执行；一个事务是不可分割的，事务中的任何一条 SQL 执行失败，已经执行成功的语句也必须撤销，`状态回退`到执行事务之前。
  - 一致性（Consistency）：一个事务执行之前和执行之后都必须处于一致性状态。拿转账来说，假设用户 A 和用户 B 两者的钱加起来一共是 5000，那么不管 A 和 B 之间如何转账，转几次账，事务结束后两个用户的钱相加起来应该还得是 5000，这就是事务的一致性。
  - 隔离性（Isolation）：数据库中的数据应满足完整性约束，事务开始和结束之间的`中间状态不会被其他事务看到`
  - 持久性（Durability）：已被提交的事务对数据库的修改应该`永久保存在数据库中`。即便是在数据库系统遇到故障的情况下也不会丢失提交事务的操作。

    一般来说，事务的 ACID 由 RDBMS 来实现的，`RDBMS 采用日志来保证事务的原子性，一致性，持久性。采用锁的机制来实现事务的隔离性`。

- 事务的隔离级别
  在事务并发操作时，可能出现的问题有：

  - 脏读：一个事务读取到了另一个事务`未提交的数据`
  - 不可重复读：一个事务读取到了另一个事务已提交的数据，导致两次读取的数据`内容`不一致
    不可重复读出现的原因就是事务并发修改记录，要避免这种情况，最简单的方法就是对要修改的记录加锁，这回导致锁竞争加剧，影响性能。另一种方法是通过 `MVCC 可以在无锁的情况下，避免不可重复读`。
  - 幻读：一个事务读取到了另一个事务已提交的数据，导致两次读取的数据`总量`不一致。幻读是由于并发事务增加记录导致的，这个不能像不可重复读通过记录加锁解决，因为对于`新增的记录根本无法加锁。需要将事务串行化，才能避免幻读`。

  **不可重复读和幻读到底有什么区别呢？**
  (1) 不可重复读是读取了其他事务更改的数据，`针对 update 操作`
  解决：使用行级锁，锁定该行，事务 A 多次读取操作完成后才释放该锁，这个时候才允许其他事务更改刚才的数据。
  (2) 幻读是读取了其他事务新增的数据，`针对 insert 和 delete 操作`
  解决：使用表级锁，锁定整张表，事务 A 多次读取数据总量之后才释放该锁，这个时候才允许其他事务新增数据。

  脏读、幻读、不可重复读是在`并发事务`的情况下才发生的，为了解决这些问题，数据库引入了隔离级别，并且不同的隔离级别可以解决不同的问题。

  - 读未提交（Read Uncommitted）：什么都不需要做，允许脏读。所有的并发事务问题都会发生。
  - 读已提交（Read Committed）：只有在事务提交后，其更新结果才会被其他事务看见。可以解决脏读问题。Oracle 等多数数据库默认都是该级别 (不重复读)。
  - 可重复读（Repeatable Read）：保证在同一事务中多次读取同样记录的结果是一致的。可以解决不可重复读问题。MySQL/InnoDB 默认级别，可以解决脏读、不可重复读。
  - 可串行化（Serializable）：最高的隔离级别，`通过强制事务串行执行，避免了幻读问题`。可以解决脏读、不可重复读、幻读。

- MVCC（多版本并发控制）
  英文全称为 Multi-Version Concurrency Control，乐观锁为理论基础的 MVCC（多版本并发控制），`MVCC 的实现没有固定的规范。每个数据库都会有不同的实现方式。`
  mysql 中，默认的事务隔离级别是可重复读（repeatable-read），为了解决不可重复读，innodb 采用了 MVCC（多版本并发控制）来解决这一问题。
  MVCC 是利用在每条数据后面加了隐藏的两列（`创建版本号和删除版本号`）: create_version 和 delete_version，`每个事务在开始的时候都会有一个递增的版本号`。

  - 增：直接 insert，创建版本号设为当前事务的版本号
  - 删: 直接将数据的删除版本号更新为当前事务的版本号
  - 改：采用 delete+add 的方式来实现，首先将当前数据`标志为删除`，然后新增一条新的数据：
  - 查：`查询操作为了避免查询到旧数据或已经被其他事务更改过的数据`，需要满足如下条件：

    - 查询时当前事务的版本号需要大于或等于创建版本号
    - 查询时当前事务的版本号需要小于删除的版本号
      即：`create_version <= current_version < delete_version`
      这样就可以避免查询到其他事务修改的数据

- MySQL 如何解决幻读
  https://juejin.cn/post/6971741501273407518
  幻读：在一个事务中使用相同的 SQL 两次读取，第二次读取到了其他事务新插入的行，则称为发生了幻读。
  谈到幻读，首先我们要引入“当前读”和“快照读”的概念，通过名字就可以理解：

  `快照读`：生成一个事务快照（ReadView），之后都从这个快照获取数据。普通 select 语句就是快照读。(clone 一份数据?)
  `当前读`：读取数据的最新版本。常见的 update/insert/delete、还有 select ... for update、select ... lock in share mode 都是当前读。
  对于快照读，`MVCC 因为因为从 ReadView 读取，所以必然不会看到新插入的行，所以天然就解决了幻读的问题。`
  而对于当前读的幻读，MVCC 是无法解决的。需要使用 Gap Lock 或 Next-Key Lock（Gap Lock + Record Lock）来解决。
  用上面的例子稍微修改下以触发当前读：`select * from user where id < 10 for update`，当使用了 Gap Lock 时，Gap 锁`会锁住 id < 10 的整个范围，因此其他事务无法插入 id < 10 的数据`，从而防止了幻读。

- InnoDB 引擎通过什么技术来保证事务的这四个特性的呢？
  持久性是通过 redo log （重做日志）来保证的，保证可以在数据库重启后恢复数据；
  原子性是通过 undo log（回滚日志） 来保证的，保证事务失败后可以回滚；
  隔离性是通过 MVCC（多版本并发控制） 或锁机制来保证的，解决 data race；
  一致性则是通过持久性+原子性+隔离性来保证；

# 数据库索引

数据库索引，是数据库管理系统中一个排序的数据结构，以协助快速查询、更新数据库表中数据。
数据库管理系统（RDBMS）通常决定索引应该用哪些数据结构。

- B-TRee 索引
- 哈希索引
- R-Tree 索引
- 位图索引

基本原则是只如果表中某列在查询过程中使用的非常频繁，那就在该列上创建索引。

- 索引类型
  根据数据库的功能，可以在数据库设计器中创建三种索引：唯一索引、主键索引和聚簇索引。
  唯一索引：唯一索引确保列中的每个值都是唯一的。唯一索引允许空值。
  主键索引：主键索引是唯一索引的特定类型。主键索引要求列中的每个值都是唯一的，且不为空。
  聚簇索引：聚簇索引对表中的行进行排序，并将行存储在磁盘上。聚簇索引只能有一个。

# NoSQL 和关系数据库结合

一般把 NoSQL 和关系数据库进行结合使用，`各取所长`，需要使用关系特性的时候我们使用关系数据库，需要使用 NoSQL 特性的时候我们使用 NoSQL 数据库，各得其所。NoSQL 数据库是关系数据库在某些方面（性能，扩展）的一个弥补。
举个简单的例子吧，比如用户评论的存储，评论大概有主键 id、评论的对象 aid、评论内容 content、用户 uid 等字段。我们能确定的是评论内容 content 肯定不会在数据库中用 where content=’’ 查询，评论内容也是一个大文本字段。`那么我们可以把主键 id、评论对象 aid、用户 id 存储在数据库，评论内容存储在 NoSQL`，这样数据库就节省了存储 content 占用的磁盘空间，从而节省大量 IO，对 content 也更容易做 Cache。

# 多人协作如何进行冲突处理

Operational transformation(OT)
OT 算法最初是为在纯文本文档的协作编辑中的一致性维护和并发控制而发明的，在本文中我们也主要掌握一致性维护相关的一些方法。

- 1.1. 协同软件的冲突
- 1.2. 操作的拆分(step)
  只要拆分得足够仔细，对于子表的所有用户行为，都可以由这些操作来组合成最终的效果。
  例如，
  复制粘贴一张子表，可以拆分为插入-重命名-更新内容；
  剪切一张子表，可以拆分为插入-更新内容-删除-移动其他子表。
  通过分析用户行为，我们可以提取出这些基本操作。
- 1.3. 操作间的冲突处理
  `n*n`，多对多的处理思路是加虚拟节点(也可以及是多个类型虚拟节点)
- 1.4. 最终一致性的实现
  OT 算法的一个核心目标，是实现`最终一致性`。
  为什么会有最终一致性的需求呢？
  > 最终一致性=没有一致性，随便搞搞就行了

# 在线文档的网络层设计思考

https://godbasin.github.io/2020/08/23/online-doc-network/
像在线文档这样大型的项目，不管是从功能职责方面，还是从代码维护方面，分层和分模块都是必然的趋势。而网络层作为与服务端直接连接的一层，有多少是我们可以做到更好的呢？

1. 认识网络层
   涉及多人在线协作的场景，从用户交互到服务端存储都会特别复杂。
   对于前端来说，从后台获取的数据到展示，分别需要经过`网络层、数据层和渲染层`。

   - 1.1. 网络层职责
     网络层无非就是做一些与服务端通信的工作，例如发起请求、异常处理、自动重试、登录态续期等。如果说，除了 HTTP 请求，可能还涉及 socket、长连接、请求数据缓存等各种功能。
     在多人协作的场景下，为了保证用户体验，一般会采用 OT 算法来进行冲突处理。而为了保证每次的用户操作都可以按照正确的时序来更新，我们会维护一个自增的版本号，每次有新的修改，都会更新版本号。因此，在这样的场景下，网络层的职责大概包括：

     - `校验`数据合法性
     - 本地数据准确的`提交`给后台（涉及队列和版本控制）
     - 协同数据正确处理后分发给数据层（涉及`冲突处理`）

   - 1.2. 网络层设计
     - 1.2.1. 连接层：管理与服务端的连接（Socket、长连接等）
       前后端通信方式有很多种，常见的包括 HTTP 短轮询（pollinf）、Websocket、HTTP 长轮询（long-polling）、SSE（Server-Sent Events）等
       每种通信方式都有各自的优缺点，包括兼容性、资源消耗、实时性等，也有可能跟业务团队自身的后台架构有关系。因此我们在设计连接层的时候，`考虑接口拓展性，应该预留对各种方式的支持`
   - 1.3. 接入层：管理数据版本、冲突处理、与数据层的连接等
     接收数据（服务端 -> 数据层）：管理来自服务端的数据；冲突处理、应用
     发送数据（数据层 -> 服务端）：管理需要提交给服务端的数据：数据提交、拉补版本

# 大型前端项目要怎么跟踪和分析函数调用链

https://godbasin.github.io/2020/06/21/trace-stash/
https://course.rs/logs/observe/intro.html

`指标(metric)`：用于表示在某一段时间内，一个行为出现的次数和分布
`日志(log)`：记录在某一个时间点发生的一次事件
`链路(trace)`：记录一次请求所经过的完整的服务链路，可能会横跨线程、进程，也可能会横跨服务(分布式、微服务)

---

1. 方案设计
   1.1. 现状
   一般来说，对于大型项目或是新人加入，维护过程（熟悉代码、定位问题、性能优化等）比较痛的有以下问题：

   - 函数执行情况黑盒
     函数调用链不清晰
     函数耗时不清楚
     代码异常定位困难
   - 用户行为难以复现

     1.2. 目标
     1.3. 整体方案设计
     1.4. 方案细节设计

2. 函数调用链的设计和实现
   2.1. 单次追踪对象
   2.2. 追踪堆栈
   2.3. 装饰器逻辑

# 谈谈依赖和解耦

大型项目总避免不了各种模块间相互依赖，如果模块之间耦合严重，随着功能的增加，项目到后期会越来越难以维护。今天我们一起来思考下，大家常说的代码解耦到底要怎么做？

1. 依赖是怎么产生的
   1.1. 接口管理
   依赖其实在`接口设计完成的时候就出来了`，虽然这是我们自己设计的接口，但它依赖于上游按照约定来调用。而上游有调整的时候，我们是需要跟随者适配或者调整的。
   `这是来自于“甲方按照约定接口来调用服务”、“乙方按照约定接口来提供服务”的依赖。`

   1.2. 状态管理
   由接口管理产生的依赖通常来自外部，而应用内部也会有依赖的产生，常见的包括状态管理和事件管理。
   最简单的，生命周期就是一种状态。
   由于程序会有状态变化，因此我们的功能实现必然依赖程序的状态。
   `这是来自于对某个程序“按照预期运行”进行合理设计而产生的依赖。`

   1.3. 功能管理
   当我们根据功能将代码拆分成一个个模块之后，功能模块的管理也同样会产生一些依赖。

   1.4. 依赖来自于约束
   为了方便管理，我们设计了一些约定，并基于“大家都会遵守约定”的前提来提供更好、更便捷的服务。
   举个例子，前端框架中为了更清晰地管理渲染层、数据层和逻辑处理，常用的设计包括 MVC、MVVM 等。
   而要使这样的架构设计发挥出效果，我们需要遵守其中的使用规范，不可以在数据层里直接修改界面等。
   可以看到，`依赖来自于对代码的设计`。

2. 依赖可以解耦吗

   2.1. 依赖的划分
   我们想要减少的，是不合理的依赖。而通过合理的设计，可以进行恰当的解耦。

   2.1.1. 无状态的函数式编程？
   我们需要对功能模块进行划分，划分出有状态和无状态的功能，来将状态管理放置到更小的范围，避免“牵一发而动全身”。
   在这里，我们进行了`状态有无的划分`。

   2.1.2. 单向流的数据管理(dag)？
   在这里，我们进行了`模块内外数据`的划分。

   2.1.3. 服务化
   服务化，是系统解耦最常用的一种方式。
   通过将功能进行业务领域的拆分，我们得到了不同领域的服务，常见的例如电商系统拆分成订单系统、购物车系统、商品系统、商家系统、支付系统等等。
   而如今打得火热的“微服务”，也都是基于领域建模的一种实现方式。
   在这里，我们进行了`业务领域`的划分。

   2.1.4. 模块化与依赖注入？
   `功能应用`的划分。

# 响应式编程在前端领域的应用

响应式编程基于观察者模式，是一种面向数据流和变化传播的声明式编程方式。
以 rxjs 为例.

1. 什么是响应式编程
   - 1.1. 异步数据流
   - 1.2. 响应式编程在前端领域
     - 1.2.1. HTTP 请求与重试
     - 1.2.2. 用户输入
       在用户频繁交互的场景，数据的流式处理可以让我们很方便地进行节流和防抖。除此之外，模块间的调用和事件通信同样可以通过这种方式来进行处理。
   - 1.3. 比较其他技术
     - 1.3.1. Promise
       - 是否有状态
         Promise 会发生状态扭转，状态扭转不可逆；
         而 Observable 是无状态的，数据流可以源源不断，可用于随着时间的推移获取多个值
       - 是否立即执行
         Promise 在定义时就会被执行；而 Observable 只有在被订阅时才会执行
       - 是否可取消
         Promise 无法取消；而 Observable 可以通过 unsubscribe 来取消订阅
     - 1.3.2. 事件
2. 响应式编程提供了怎样的服务
   - 2.1. 热观察与冷观察
     Hot Observable，可以理解为现场直播，我们进场的时候只能看到即时的内容
     Cold Observable，可以理解为点播（电影），我们打开的时候会从头播放(如果我们想要在拉群后，自动同步之前的聊天记录，通过冷观察就可以做到。)
   - 2.2. 合流
     merge：多个 Observable 合并为一个
     combine：所有的 Observable 都有值之后，才会触发
   - 2.3. 其他使用方式
     - 2.3.1. timer
     - 2.3.2. 数组/可迭代对象

# VSCode 源码解读：事件系统设计

TL;DR

- 提供标准化的 Event 和 Emitter 能力
- 通过注册 Emitter，并对外提供类似生命周期的方法 onXxxxx 的方式，来进行事件的订阅和监听
- 通过提供通用类 Disposable，统一管理相关资源的注册和销毁
- 通过使用同样的方式`this._register()`注册事件和订阅事件，将事件相关资源的处理统一挂载到 dispose()方法中

---

前端大型项目中要怎么管理满天飞的事件、模块间各种显示和隐式调用的问题

**看源码的方式有很多种，带着疑问有目的性地看，会简单很多。**

1. VS Code 事件

   - 1.1. Q1: VS Code 中的事件管理代码在哪？base\common\event.ts
   - 1.2. Q2: VS Code 中的事件都包括了哪些能力？
     除了常见的 once 和 DOM 事件等兼容，还提供了比较丰富的事件能力：

     - 防抖动
     - 可链式调用
     - 缓存
     - Promise 转事件

   - 1.3. Q3: VS Code 中的事件的触发和监听是怎么实现的？
     Emitter 以 Event 为对象，以简洁的方式提供了事件的订阅、触发、清理等能力

     ```ts
     // 这是事件发射器的一些生命周期和设置

     export interface EmitterOptions {
       onFirstListenerAdd?: Function
       onFirstListenerDidAdd?: Function
       onListenerDidAdd?: Function
       onLastListenerRemove?: Function
       leakWarningThreshold?: number
     }

     export class Emitter<T> {
       // 可传入生命周期方法和设置
       constructor(options?: EmitterOptions) {}

       // 允许大家订阅此发射器的事件
       get event(): Event<T> {
         // 此处会根据传入的生命周期相关设置，在对应的场景下调用相关的生命周期方法
       }

       // 向订阅者触发事件
       fire(event: T): void {}

       // 清理相关的 listener 和队列等
       dispose() {}
     }
     ```

   - 1.4. Q4: 项目中的事件是怎么管理的？
     Emitter 似乎有些简单了，我们只能看到单个事件发射器的使用。那各个模块之间的事件订阅和触发又是怎么实现的呢？

     - 把 eventEmitter 绑在对象上(注册事件发射器)
     - 对外提供定义的事件
     - 在特定时机向订阅者触发事件

     ```ts
     // 这里我们只摘录相关的代码
     class WindowManager {
       public static readonly INSTANCE = new WindowManager()
       // 注册一个事件发射器
       private readonly _onDidChangeZoomLevel = new Emitter<number>()
       // 将该发射器允许大家订阅的事件取出来
       public readonly onDidChangeZoomLevel: Event<number> = this._onDidChangeZoomLevel.event

       public setZoomLevel(zoomLevel: number, isTrusted: boolean): void {
         if (this._zoomLevel === zoomLevel) {
           return
         }
         this._zoomLevel = zoomLevel
         // 当 zoomLevel 有变更时，触发该事件
         this._onDidChangeZoomLevel.fire(this._zoomLevel)
       }
     }

     const instance = new WindowManager(opts)
     instance.onDidChangeZoomLevel(() => {
       // 该干啥干啥
     })
     ```

   - 1.5. Q5: 事件满天飞，不会导致性能问题吗？
     如果在某个组件里做了事件订阅这样的操作，当组件销毁的时候是需要取消事件订阅的。
     否则该订阅内容会在内存中一直存在，除了一些异常问题，还可能引起内存泄露。
     一些地方的使用方式是：

     ```ts
     // 这里使用了this._register(new Emitter<T>())这样的方式注册事件发射器，我们能看到该方法继承自Disposable。
     export class Scrollable extends Disposable {
       private _onScroll = this._register(new Emitter<ScrollEvent>())
       public readonly onScroll: Event<ScrollEvent> = this._onScroll.event

       private _setState(newState: ScrollState): void {
         const oldState = this._state
         if (oldState.equals(newState)) {
           return
         }
         this._state = newState
         // 状态变更的时候，触发事件
         this._onScroll.fire(this._state.createScrollEvent(oldState))
       }
     }

     export abstract class Disposable implements IDisposable {
       // 用一个 Set 来存储注册的事件发射器
       private readonly _store = new DisposableStore()

       constructor() {
         trackDisposable(this)
       }

       // 处理事件发射器
       public dispose(): void {
         markTracked(this)

         this._store.dispose()
       }

       // 注册一个事件发射器
       protected _register<T extends IDisposable>(t: T): T {
         if ((t as unknown as Disposable) === this) {
           throw new Error('Cannot register a disposable on itself!')
         }
         return this._store.add(t)
       }
     }
     ```

     Dispose 模式主要用来资源管理，资源比如内存被对象占用，则会通过调用方法来释放。

     ```ts
     export interface IDisposable {
       dispose(): void
     }
     export class DisposableStore implements IDisposable {
       private _toDispose = new Set<IDisposable>()
       private _isDisposed = false

       // 处置所有注册的 Disposable，并将其标记为已处置
       // 将来添加到此对象的所有 Disposable 都将在 add 中处置。
       public dispose(): void {
         if (this._isDisposed) {
           return
         }

         markTracked(this)
         this._isDisposed = true
         this.clear()
       }

       // 丢弃所有已登记的 Disposable，但不要将其标记为已处置
       public clear(): void {
         this._toDispose.forEach(item => item.dispose())
         this._toDispose.clear()
       }

       // 添加一个 Disposable
       public add<T extends IDisposable>(t: T): T {
         markTracked(t)
         // 如果已处置，则不添加
         if (this._isDisposed) {
           // 报错提示之类的
         } else {
           // 未处置，则可添加
           this._toDispose.add(t)
         }
         return t
       }
     }
     ```

   - 1.6. Q6: 上面只销毁了事件触发器本身的资源，那对于订阅者来说，要怎么销毁订阅的 Listener 呢？

# VSCode 源码解读：IPC 通信机制

1. Electron 的 通信机制
2. Electron 与 NW.js

   - 2.0.1. NW.js 内部架构
   - 2.0.2. Electron 内部架构
     Electron 强调 Chromium 源代码和应用程序进行分离，因此并没有将 Node.js 和 Chromium 整合在一起。
     在 Electron 中，分为主进程(main process)和渲染器进程(renderer processes)：
     那么，不在一个进程当然涉及`跨进程通信`。于是，在 Electron 中，可以通过以下方式来进行主进程和渲染器进程的通信：

     - 利用 ipcMain 和 ipcRenderer 模块进行 `IPC` 方式通信，它们是处理应用程序后端（ipcMain）和前端应用窗口（ipcRenderer）之间的进程间通信的事件触发。
     - 利用 remote 模块进行 `RPC` 方式通信。
       > remote 模块返回的每个对象（包括函数），表示主进程中的一个对象（称为远程对象或远程函数）。当调用远程对象的方法时，调用远程函数、或者使用远程构造函数 (函数) 创建新对象时，实际上是在发送同步进程消息

3. VSCode 的通信机制

   - 3.1. VSCode 多进程架构
     VSCode 采用多进程架构，VSCode 启动后主要有下面的几个进程：

     主进程
     渲染进程，多个，包括 Activitybar、Sidebar、Panel、Editor 等等
     插件宿主进程
     Debug 进程
     Search 进程

   - 3.2. IPC 通信
     主进程和渲染进程的通信基础还是 Electron 的 webContents.send、ipcRender.send、ipcMain.on。
     - 3.2.1. 协议
       ```ts
       export interface IMessagePassingProtocol {
         send(buffer: VSBuffer): void
         onMessage: Event<VSBuffer>
       }
       ```
     - 3.2.2. 频道：VSCode 通过频道来区分不同的通信类型
       ```ts
       /**
        * IChannel是对命令集合的抽象
        * call 总是返回一个至多带有单个返回值的 Promise
        */
       export interface IChannel {
         call<T>(command: string, arg?: any, cancellationToken?: CancellationToken): Promise<T>
         listen<T>(event: string, arg?: any): Event<T>
       }
       ```
     - 3.2.3. 客户端与服务端
     - 3.2.4. 连接：一个连接(Connection)由 ChannelClient 和 ChannelServer 组成。
       ```ts
       interface Connection<TContext> extends Client<TContext> {
         readonly channelServer: ChannelServer<TContext> // 服务端
         readonly channelClient: ChannelClient // 客户端
       }
       ```

# 复杂渲染引擎架构与设计

## 收集与渲染

1. 渲染数据的收集
2. 收集与绘制的功能划分
3. 渲染数据享元
