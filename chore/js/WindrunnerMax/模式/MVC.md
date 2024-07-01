- 定义
  MVC 即模型 Model、视图 View、控制器 Controller，用一种将`业务逻辑、数据、视图`分离的方式组织架构代码，通`过分离关注点的方式来支持改进应用组织方式`，其促成了业务数据 Model 从用户界面 View 中分离出来，还有第三个组成部分 Controller 负责管理传统意义上的业务逻辑和用户输入，通常将 MVC 模式看作架构型设计模式。
  View -> Controller -> Model -> View
- 实现
  View 传送指令到 Controller。
  Controller 完成业务逻辑后，要求 Model 改变状态。
  Model 将新的数据发送到 View，用户得到反馈。
