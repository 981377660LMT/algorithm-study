- 定义
  MVC 即模型 Model、视图 View、管理器 Presenter，MVP 模式从 MVC 模式演变而来，通过管理器将视图与模型巧妙地分开，即将 `Controller 改名为 Presenter`，同时改变了通信方向，MVP 模式模式不属于一般定义的 23 种设计模式的范畴，而通常将其看作广义上的架构型设计模式

  View <-> Controller <-> Model

  MVP 模式通过解耦 View 和 Model，完全分离视图和模型使职责划分更加清晰，由于 View 不依赖 Model，可以将 View 抽离出来做成组件，其只需要提供一系列接口提供给上层操作

- 实现
