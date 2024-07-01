- 定义
  MVVM 全称 Model-View-ViewModel 是基于 MVC 和 MVP 体系结构模式的改进，MVVM 就是 MVC 模式中的 View 的状态和行为抽象化，`将视图 UI 和业务逻辑分开`，更清楚地将用户界面 UI 的开发与应用程序中业务逻辑和行为的开发区分开来，MVP 模式模式不属于一般定义的 23 种设计模式的范畴，而通常将其看作广义上的架构型设计模式。

- 实现
  View <- ViewModel <-> Model
  在 Model 更新时，ViewModel 通过绑定器将数据更新到 View，在 View 触发指令时，会通过 ViewModel 传递消息到 Model

  MVVM 模式与 MVP 模式行为基本一致，`主要区别是其通常采用双向绑定 data-binding，即将 View 和 Model 的同步逻辑自动化了`，以前 Presenter 负责的 View 和 Model 同步不再手动地进行操作
