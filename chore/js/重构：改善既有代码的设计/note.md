[重构 改善既有代码的设计第二版](https://book-refactoring2.ifmicro.com/)
https://github.com/notes-folder/refactoring
https://juejin.cn/post/7255581952874102821
https://wxyclark.github.io/Work/Architect/refactor.html
需要批判性地阅读本书

- `便于理解，便于修改`，是重构这个方法最直白的解释了
- 伪代码反而是最好的代码，因为它不注重细节，注重业务，便于人们理解，“我不管你怎么实现，我就要这种效果”，几乎没有 bug
- 好的代码，应该是接近人类语言或者说人的思维，让人能够很顺畅的读懂、去理解代码是去做什么的，怎么做的。`可读性`在某种意义上，是好代码和坏代码的区别。
- DDD：构建合理的业务模型，如果建模是错的，不能反映客观规律事实，后续就是噩梦

## 第 1 章 重构，第一个示例

- 消除局部变量:将业务流转逻辑与计算逻辑拆分，简单点说，流转逻辑中使用的计算参数使用计算函数进行拆分出来
- 无论每次重构多么简单，养成重构后即运行测试的习惯非常重要
- `好代码的检验标准就是人们是否能轻而易举地修改它`
- 抽离函数减小复杂度，提高可读性

## 第 2 章 重构的原则

- 何时不应该重构
  1. `丑陋的代码能被隐藏在一个 API 之下：`
     如果我看见一块凌乱的代码，但并不需要修改它，那么我就不需要重构它。如果丑陋的代码能被隐藏在一个 API 之下，我就可以容忍它继续保持丑陋。只有当我需要理解其工作原理时，对其进行重构才有价值
  2. `重写比重构还容易`
- 如果经理不能接受代码重构需要花时间，那就不要给经理说
- 何时重构：三次法则(事不过三，三则重构)
- 重构的两种核心思路：`找到缺乏间接层的地方, 加入一个间接层. 找到多余的间接层, 删去这种过度设计`
- 普通接口的改名只要直接文本替换就能实现, 但是那些已经发布给外部人员使用的接口则不容易. 对于这种情况最好先让旧接口调用新接口并让两者共存, 标记为“deprecated”, 然后在足够久的未来后删去旧接口

## 第 3 章 代码的坏味道

- 设计 API 时，可以将查询和修改函数分离确保调用者不会调到有副作用的代码
- 神秘命名（Mysterious Name）
  计算科学中最难的两件事是命名和缓存失效（There are only two hard things in computer science: cache invalidation and naming things）
- 重复代码（Duplicated Code）?
  有时候不一定，只是碰巧相似，但是未来可能会变得不同
- 过长函数（Long Function）
  有时候不一定，做好抽象即可
  `条件表达式和循环常常也是提炼的信号`
  编码中要善于提炼小函数，`通过对小函数的复用或者流程编排，达到面向对象的效果。`遇到想写注释，那就是可以抽取为函数的信号，这是面向过程和面向对象的最显著的区别
- 过大的类：
  表示一个类的功能太多，可以将类的功能拆分为两个不同功能单一的类；
- switch 惊悚现身
  出现 switch、if…else 的时候，需要联想到多台和策略的抽象来替换；
- 全局数据。（Global Data）(PS：慎用单例模式)
  一般需要通过封装变量的方式将这种全局的数据进行封装，通过方法进行改变
  封装变量的好处：1. 利于监控与定位问题 2. 如果这个变量需要跟一些协同变量同时改变，函数可以做到收口
- 发散式变化（Divergent Change）
  某个类经常因为不同的原因在不同的方向上发生变化
  违反单一职责
- 霰弹式修改(shotgun surgery)
  如果每遇到某种变化，你都必须在`许多不同的类内做出许多小修改`
  `模块拆分的过细`
- **纯数据类（Data Class）**
  类比 DDD 思想，尽可能的让我们的数据类也能拥有行为
  btw，`纯数据类用 type，否则用 interface?`
- 被拒绝的遗赠（Refused Bequest）
- 依恋情结（feature envy）
  函数对某个类的兴趣高过对自己所处类的兴趣
- 数据泥团（Data Clumps）
  `总是绑在一起出现的数据应该有自己的对象`
  删掉众多数据中的一项。如果这么做，其他数据有没有因而失去意义？如果它们不再有意义，这就是一个明确信号：你应该为它们产生一个新对象
- 基本类型偏执
  不要吝啬使用小对象, `将一些基本类型包装为类很实用`
- 令人迷惑的临时字段
  如果一个类中有一个复杂算法, `需要好几个临时变量的协作, 那么我们应该将这些操作和变量拆分提取到一个独立类(或者属性?)中`, 提炼后称为函数对象. 由这个类去维护那些临时变量, 也方便我们传递这些变量
- 变换函数（Change Preventers）
  如果变换函数返回的本质是原来对象上`添加了更多的信息，用 enrich 命名，新类型也可以套一层rich<T>=T&{...}`。
  如果它生成的是跟原来`完全不同的对象，就用 transform/convert/to 命名`
- 中间人（Middle Man）
  如果某个类接口有`一半的函数都委托给其它类，就属于过度委托`
- 不适当的亲密（Inappropriate Intimacy）
  两个类关系过于亲密，花费太多时间去探究彼此的 private 成分，此时建议要么合并，要么将共同点提取到一个新类

## 第 4 章 构筑测试体系

- 每当你收到一个 bug 时，请先写一个单元测试来暴露这个问题
- 不要幻想把所有场景都覆盖测试，只覆盖那些关键、容易出错的地方就够了
- 将 UI 与业务逻辑隔离，这样才能写测试用例直接测逻辑部分
- 确保测试 case 在不该通过时真的会失败

## 第 5 章 介绍重构名录(Most Useful Set Of Refactoring)

每个重构手法都有如下 5 个部分

- 名称（name）
- 概括（sketch）
- 动机（motivation）
- 做法（mechanics）
- 范例（examples）

## 第 6 章 第一组重构

- 提炼函数：如果你需要花时间浏览一段代码才能弄懂他到底在干什么，那么就应该将其提炼到一个函数中，并根据他所做的事为其命名。读代码时间就可以`一眼看到函数的用途而不用关心他具体的实现(声明式，关注结果，通过封装降低复杂度)`。

## 第 7 章 封装(Encapsulate)

封装集合：集合的不当修改可能会带来意想不到的 bug，我们建议将集合的修改操作限制在封装集合的类中（通常是在类中设置“添加”和“移除”方法）

## 第 8 章 搬移特性( Moving Features)

- 搬移函数到更合适的类中
  Move a function when it references elements in other contexts more than the one it currently resides in

## 第 9 章 重新组织数据(Organizing Data)

## 第 10 章 简化条件逻辑(Simplifying Conditional Logic)

- 分解条件表达式 Decompose Conditional

```js
let charge
if (!aDate.isBefore(plan.summerStart) && !aData.isAfter(plan.summerEnd)) {
  charge = quantity * plan.summerRate
} else {
  charge = quantity * plan.regularRate + plan.regularServiceCharge
}

const charge = isSummer() ? summerCharge() : regularCharge()
```

- 以卫语句取代件套表达式 Replace Nested Conditional with Guard Clauses

```js
function getPayAmount() {
  let result
  if (isDead) {
    result = deadAmount()
  } else {
    if (isSeparated) {
      result = separatedAmount()
    } else {
      if (isRetired) {
        result = retiredAmount()
      } else {
        result = normalPayAmount()
      }
    }
  }
  return result
}

function getPayAmount() {
  if (isDead) return deadAmount()
  if (isSeparated) return separatedAmount()
  if (isRetired) return retiredAmount()
  return normalPayAmount()
}
```

- 以多态取代条件表达式 Replace Conditional with Polymorphism
  统一成接口方法,不同的子类实现不同的方法

## 第 11 章 重构 API(Refactoring APIS)

- 将查询函数和修改函数分离：
  将修改函数中的`查询代码块提取成一个查询函数，修改函数中根据查询函数的结果做相应的修改操作`。
  `Immutable function (query only) is easy to test and reuse`
- 移除默认参数?
  有一种说法是，默认参数造成了信息泄露。像 golang 就不支持默认参数。
- 保持对象完整：
  如果我看见代码从一个记录结构中导出几个值，然后又把这几个值一起传递给一个函数，我会`更愿意把整个记录传给这个函数，在函数体内部导出所需的值`
- 以工厂函数取代构造函数：
  `使用工厂函数可以避免一些构造函数的局限性。`像 golang 就没有构造函数，copy 更容易写。

## 第 12 章 处理继承关系(Dealing With Inheritance)

- 用委托取代子类或以委托取代超类，从而将继承体系转化成委托调用
- 函数上移：
  如果某个函数在各个子类中的函数体都相同，这就是显而易见的函数上移适用场合。
- 函数下移：
  如果超类中某个函数只与一个（或少数几个）子类有关，那么最好将其`从超类中挪走，放到真正关心它的子类中去。`
  如果不知道或者不明确哪些子类需要这个函数时，那就得用`多态`取代条件表达式，只留些共用的行为在超类
- **如果超类的一些函数对子类并不适用，就说明我不应该通过继承来获得超类的功能**
- 隐藏委托关系 Hide Delegate
  调用者不需要知道委托对象(数据结构的 size 封装)

  ```js
  manager = aPerson.department.manager

  manager = aPerson.manager
  class Person {
    get manager() {
      return this.department.manager
    }
  }
  ```

- 移除中间人 Remove Middle Man
  **delegating methods 过多时，需要移除(多维表格问题)**

  ```js
  manager = aPerson.manager

  class Person {
    get manager() {
      return this.department.manager
    }
  }

  manager = aPerson.department.manager
  ```

- 以委托取代超类 Replace Superclass with Delegate
- 以委托取代子类 Replace Subclass with Delegate
  - 继承这张牌只能打一次，如果后续子类朝着不同的方向演化，那继承就明显不合适
  - 继承给类之间引入了非常重的耦合关系，在超类上做任何修改都会影响到子类
