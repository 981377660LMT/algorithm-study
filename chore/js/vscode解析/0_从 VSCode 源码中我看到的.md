https://developer.aliyun.com/article/1219001

## 代码结构

功能有清晰的分层，依赖关系是清晰单向的
每层针对不同端还有自己的实现分层

## 依赖注入

VSCode 中将各种功能解耦成一个个的 service，每个能力组合一个个 service 实现。这里的 service 通过依赖注入的方式附着到使用方上

> VSCode 中通过装饰器注解的方式来声明依赖关系，不过它并没有直接使用 reflect-metadata，而是基于 decorator 标注元信息实现了一套自己的依赖注入模式。
