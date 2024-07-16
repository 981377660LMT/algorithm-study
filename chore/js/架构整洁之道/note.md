架构整洁之道
https://github.com/leewaiho/Clean-Architecture-zh/blob/master/docs/part1.md

旧瓶装新酒

## Part1. INTRODUCTION 概述

## Chap1. WHAT IS DESIGN AND ARCHITECTURE? 设计与架构到底是什么

设计（Design）：
架构（Architecture）：

- 乱麻系统的特点(THE SIGNATURE OF A MESS)
  这种系统一般都是没有经过设计，匆匆忙忙被构建起来的。然后为了加快发布的速度，拼命地往团队里加入新人，同时加上决策层对代码质量提升和设计结构优化存在着持续的、长久的忽视，这种状态能持续下去就怪了。
- 研发团队最好的选择是`清晰地认识并避开工程师们过度自信的特点，开始认真地对待自己的代码架构，对其质量负责`

## Chap2. A TALE OF TWO VALUES 两个价值维度

- 对于每个软件系统，我们都对以通过`行为`和`架构`两个维度来休现它的实际价值
  行为价值：满足需求，系统正常工作；
  架构价值：满足需求的同时，保持软件系统的灵活性。software。“ware” 的意思是“产品”，而 “soft” 的意思，不言而喻，是指软件的灵活性。软件系统必须够“软” 也就是说，`软件应该容易被修改。`

  好的系统架构设计应该尽可能做到与“形状”无关。

- 系统行为更重要，还是系统架构的灵活性更重要；哪个价值更大?
  FIGHT FOR THE ARCHITECTURE 为好的软件架构而持续斗争

## Part2. STARTING WITH THE BRICKS: PROGRAMMING PARADIGMS 从基础构件开始：编程范式

## Chap3. PARADIGM OVERVIEW 编程范式总览

直到今天，我们也一共只有三个编程范式，而且未来几乎不可能再出现新的：

- 结构化编程（structured programming）

  Structured programming imposes discipline on direct transfer of control.
  结构化编程对程序`控制权的直接转移`进行了限制和规范。
  限制了 goto 语句

- 面向对象编程（object-oriented programming）

  Object-oriented programming imposes discipline on indirect transfer of control.
  面向对象编程对程序`控制权的间接转移`进行了限制和规范。
  限制了函数指针

- 函数式编程（functional programming）
  Functional programming imposes discipline upon assignment.
  函数式编程对程序中的`赋值`进行了限制和规范
  限制了赋值语句

多态是我们跨越架构边界的手段，函数式编程是我们规范和限制数据存放位置与访问权限的手段，结构化编程则是各模块的算法实现基础
Notice how well those three align with the three big concerns of architecture: function, separation of components, and data management.
这和软件架构的三大关注重点不谋而合：`功能性、组件独立性以及数据管理`

## Chap4. STRUCTURED PROGRAMMING 结构化编程

All programs can be constructed from just three structures: `sequence, selection, and iteration`
人们可以用顺序结构、分支结构、循环结构这三种结构构造出任何程序

## Chap5. OBJECT-ORIENTED PROGRAMMING 面向对象编程

多态其实不过就是函数指针的一种应用：
在 C++中，类中的每个虚函数（virtual function）的地址都被记录在一个名叫 vtable 的数据结构里。我们对虚函数的每次调用都要先查询这个表，其衍生类的构造函数负责将该衍生类的虚函数地址加载到整个对象的 vtable 中

## Chap6. FUNCTIONAL PROGRAMMING 函数式编程

**一个架构设计良好的应用程序应该将状态修改的部分和不需要修改状态的部分隔离成单独的组件，然后用合适的机制来保护可变量**

## Part3. DESIGN PRINCIPLES 设计原则

## Chap7. SRP: THE SINGLE RESPONSIBILITY PRINCIPLE SRP：单一职责原则

任何一个软件模块都应该只对一个用户（User）或系统利益相关者（Stakeholder）负责。
任何一个软件模块都应该只对某一类行为者(actor)负责。

## Chap8. OCP: THE OPEN-CLOSED PRINCIPLE OCP：开闭原则

## Chap9. LSP: THE LISKOV SUBSTITUTION PRINCIPLE LSP：里氏替换原则

## Chap10. ISP: THE INTERFACE SEGREGATION PRINCIPLE ISP：接口隔离原则

## Chap11. DIP: THE DEPENDENCY INVERSION PRINCIPLE DIP：依赖反转原则

- 在源代码层次的依赖关系中引用抽象类型
- 反例：java String 类，原因：
  - 软件系统在实际构造中不可避免地需要依赖到一些具体实现
  - String 类本身是非常稳定的
- 在应用 DIP 时，我们也`不必考虑稳定的操作系统或者平台设施，因为这些系统接口很少会有变动`
  ` 我们主要应该关注的是软件系统内部那些会经常变动的（volatile）具体实现模块`，这些模块是不停开发的，也就会经常出现变更。

## Part4. COMPONENT PRINCIPLES 组件构建原则

## Chap12. COMPONENTS 组件

组件是软件的部署单元，是整个软件系统在部署过程中可以独立完成部署的最小实体

## Chap13. COMPONENT COHESION 组件聚合

The granule of reuse is the granule of release.
软件复用的最小粒度应等同于其发布的最小粒度

## Chap14. COMPONENT COUPLING 组件耦合

- 无依赖环原则
- 每周构建
- 打破循环依赖
  1. 依赖倒置原则(DIP)：调用接口
  2. 抽离公共依赖
- 稳定性指标
  根据入向依赖和出向依赖的数量来判断组件的稳定性
- **并不是所有组件都应该是稳定的**
  取决于抽象程度
- **一个组件的抽象化程度应该与其稳定性保持一致**
  衡量抽象化程度：组件中类的数量/组件中抽象类和接口的数量

## Part5. ARCHITECTURE 软件架构

**软件架构设计的主要目标是支撑软件系统的全生命周期，设计良好的架构可以让系统便于理解、易于修改、方便维护，并且能轻松部署。软件架构的终极目标就是最大化程序员的生产力，同时最小化系统的总运营成本**

## Chap15. WHAT IS ARCHITECTURE? 什么是软件架构

- 软件架构师自身需要是程序员，并且必须一直坚持做一线程序员，绝对不要听从那些说应该让软件架构师从代码中解放出来以专心解决高阶问题的伪建议。也许软件架构师生产的代码量不是最多的，但是他们必须不停地承接编程任务。如果不亲身承受因系统设计而带来的麻烦，`就体会不到设计不佳所带来的痛苦，接着就会逐渐迷失正确的设计方向`

  - 开发(Development)
  - 部署(Deployment)
  - 运行(Operation)
  - 维护(Maintenance)

## Chap16. INDEPENDENCE 独立性

## Chap17. BOUNDARIES: DRAWING LINES 划分边界

## Chap18. BOUNDARY ANATOMY 边界剖析

## Chap19. POLICY AND LEVEL 策略与层次

## Chap20. BUSINESS RULES 业务逻辑

## Chap21. SCREAMING ARCHITECTURE 尖叫的软件架构

## Chap22. THE CLEAN ARCHITECTURE 整洁架构

## Chap23. PRESENTERS AND HUMBLE OBJECTS 展示器和谦卑对象

## Chap24. PARTIAL BOUNDARIES 不完全边界

## Chap25. LAYERS AND BOUNDARIES 层次与边界

## Chap26. THE MAIN COMPONENT Main 组件

## Chap27. SERVICES: GREAT AND SMALL 服务：宏观与微观

## Chap28. THE TEST BOUNDARY 测试边界

## Chap29. CLEAN EMBEDDED ARCHITECTURE 整洁的嵌入式架构

## Part6. DETAILS 实现细节

## Chap30. THE DATABASE IS A DETAIL 数据库只是实现细节

## Chap31. THE WEB IS A DETAIL Web 是实现细节

## Chap32. FRAMEWORKS ARE DETAILS 应用程序框架是实现细节

## Chap33. CASE STUDY: VIDEO SALES 案例分析：视频销售网站

## Chap34. THE MISSING CHAPTER 拾遗

---

架构整洁之道网友笔记
https://lailin.xyz/post/go-training-week4-clean-arch.html

- 架构的主要目的是支持系统的生命周期。良好的架构使系统易于理解，易于开发，易于维护和易于部署。最终目标是最小化系统的寿命成本并最大化程序员的生产力

https://www.meetkiki.com/archives/%E5%86%8D%E8%AF%BB%E3%80%8A%E6%9E%B6%E6%9E%84%E6%95%B4%E6%B4%81%E4%B9%8B%E9%81%93%E3%80%8B

- 面向对象的软件设计到底是什么，怎么用一句话形容这个行为？很多人只能意会，无法言表，看完这本书，终于知道怎么说了，那就是两个控制：

1. 分离系统中变与不变的内容，对变化的部分进行控制
2. 分析系统中的各种依赖关系（对象、组件），对这些依赖关系进行控制

https://juejin.cn/post/7160552589250166820

https://juejin.cn/post/7124624443279671332#heading-22
