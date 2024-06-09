# Programming paradigms (编程范式)

Imperative programs describe `how` to do something, whereas declarative programs describe `what` to do
命令式程序说明了`如何做`，然而声明式程序说明`做了什么`

命令式程序说明了如何做，然而声明式程序说明做了什么

## 命令查询分离原则 Command-Query-Separation(CQS)

函数不应该产生抽象的副作用，只允许命令（过程）产生副作用——Bertrand Meyer:《面向对象软件构造》

## 最小惊奇原则 Principle of least astonishment (POLA)

系统的组件应该像人们期望的那样工作，而不应该给用户一个惊奇。

## 统一访问原则 Uniform-Access

一个模块提供的所有服务都应该通过一个统一的符号来提供，而这个符号并不表明它们是通过存储还是通过计算来实现的。——Bertrand Meyer:《面向对象软件构造》
