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
- 面向对象编程（object-oriented programming）
- 函数式编程（functional programming）

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
