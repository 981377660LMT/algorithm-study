https://juejin.cn/post/6876943860988772360
BetterScroll 是一款重点解决移动端（已支持 PC）各种滚动场景需求的插件。它的核心是借鉴的 iscroll 的实现，它的 API 设计基本兼容 iscroll，在 iscroll 的基础上又扩展了一些 feature 以及做了一些性能优化。
BetterScroll 1.0 共发布了 30 多个版本，npm 月下载量 5 万，累计 star 数 12600+。那么为什么升级 2.0 呢？
v2 版本的初衷源于社区的一个需求：

`BetterScroll 能不能支持按需加载？`

为了支持插件的`按需加载`，BetterScroll 2.0 采用了 `插件化` 的架构设计。

CoreScroll 作为最小的滚动单元，`暴露了丰富的事件以及钩子`，其余的功能都由不同的插件来扩展，这样会让 BetterScroll 使用起来更加的灵活，也能适应不同的场景。

该项目采用的是 monorepos 的组织方式，使用 lerna 进行多包管理，`每个组件都是一个独立的 npm 包：`

与西瓜播放器一样，`BetterScroll 2.0 也是采用 插件化 的设计思想`，CoreScroll 作为最小的滚动单元，其余的功能都是通过插件来扩展。比如长列表中常见的上拉加载和下拉刷新功能，在 BetterScroll 2.0 中这些功能分别通过 pull-up 和 pull-down 这两个插件来实现。

插件化的好处之一就是可以支持按需加载，此外把`独立功能`都拆分成独立的插件，会让核心系统更加稳定，拥有一定的健壮性。

1. 更好的智能提示
   BetterScroll 团队充分利用了 TypeScript 接口自动合并的功能，让开发者在使用某个插件时，能够有对应的 Options 提示以及 bs（BetterScroll 实例）能够有对应的方法提示。
2. 插件化架构的优点
   灵活性高：整体灵活性是对环境变化快速响应的能力。由于插件之间的低耦合，改变通常是隔离的，可以快速实现。通常，核心系统是稳定且快速的，具有一定的健壮性，几乎不需要修改。
   可测试性：插件可以独立测试，也很容易被模拟，不需修改核心系统就可以演示或构建新特性的原型。
   性能高：虽然插件化架构本身不会使应用高性能，但通常使用插件化架构构建的应用性能都还不错，因为可以自定义或者裁剪掉不需要的功能。

3. 三大件

   1. 插件管理
   2. 插件连接

   ```shell
   $ npm install @better-scroll/pull-up --save
   ```

   ```JS
    import BScroll from '@better-scroll/core'
    import Pullup from '@better-scroll/pull-up'
    // BScroll.use 方法来注册插件：

    BScroll.use(Pullup)

    // 实例化 BetterScroll 时需要传入 PullUp 插件的配置项。
    new BScroll('.bs-wrapper', {
      pullUpLoad: true
    })

    当我们调用 BScroll.use(Pullup) 方法时，会先获取当前插件的名称，然后判断当前插件是否已经安装过了。如果已经安装则直接返回 BScrollConstructor 对象，否则会对插件进行注册。即把当前插件的信息分别保存到 pluginsMap（{}） 和 plugins（[]） 对象中：
   ```

   3. 插件通信
      BScrollConstructor 类，该类继承了 EventEmitter 事件派发器：

4. 在工程化方面，BetterScroll 使用了业内一些常见的解决方案：
   - lerna：Lerna 是一个管理工具，用于管理包含多个软件包（package）的 JavaScript 项目。
   - prettier：Prettier 中文的意思是漂亮的、美丽的，是一个流行的代码格式化的工具。
   - tslint：TSLint 是可扩展的静态分析工具，用于检查 TypeScript 代码的可读性，可维护性和功能性错误。
   - commitizen & cz-conventional-changelog：用于帮助我们生成符合规范的 commit message。
   - husky：husky 能够防止不规范代码被 commit、push、merge 等等。
   - jest：Jest 是由 Facebook 维护的 JavaScript 测试框架。
   - coveralls：用于获取 Coveralls.io 的覆盖率报告，`并在 README 文件中添加一个不错的覆盖率按`钮。
   - vuepress：Vue 驱动的静态网站生成器，它用于生成 BetterScroll 2.0 的文档。
