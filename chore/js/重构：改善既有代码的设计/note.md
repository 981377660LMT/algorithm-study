[重构 改善既有代码的设计第二版](https://book-refactoring2.ifmicro.com/)

需要批判性地阅读本书

- `便于理解，便于修改`，是重构这个方法最直白的解释了
- 伪代码反而是最好的代码，因为它不注重细节，注重业务，便于人们理解，“我不管你怎么实现，我就要这种效果”，几乎没有 bug
- 好的代码，应该是接近人类语言或者说人的思维，让人能够很顺畅的读懂、去理解代码是去做什么的，怎么做的。`可读性`在某种意义上，是好代码和坏代码的区别。

## 第 1 章 重构，第一个示例

- 无论每次重构多么简单，养成重构后即运行测试的习惯非常重要
- `好代码的检验标准就是人们是否能轻而易举地修改它`
- 抽离函数减小复杂度，提高可读性

```js
function statement(invoice, plays) {
  let result = `Statement for ${invoice.customer}\n`
  for (let perf of invoice.performances) {
    result += ` ${playFor(perf).name}: ${usd(amountFor(perf))} (${perf.audience} seats)\n`
  }
  result += `Amount owed is ${usd(totalAmount())}\n`
  result += `You earned ${totalVolumeCredits()} credits\n`
  return result

  function totalAmount() {
    let result = 0
    for (let perf of invoice.performances) {
      result += amountFor(perf)
    }
    return result
  }
  function totalVolumeCredits() {
    let result = 0
    for (let perf of invoice.performances) {
      result += volumeCreditsFor(perf)
    }
    return result
  }
  function usd(aNumber) {
    return new Intl.NumberFormat('en-US', { style: 'currency', currency: 'USD', minimumFractionDigits: 2 }).format(aNumber / 100)
  }
  function volumeCreditsFor(aPerformance) {
    let result = 0
    result += Math.max(aPerformance.audience - 30, 0)
    if ('comedy' === playFor(aPerformance).type) result += Math.floor(aPerformance.audience / 5)
    return result
  }
  function playFor(aPerformance) {
    return plays[aPerformance.playID]
  }
  function amountFor(aPerformance) {
    let result = 0
    switch (playFor(aPerformance).type) {
      case 'tragedy':
        result = 40000
        if (aPerformance.audience > 30) {
          result += 1000 * (aPerformance.audience - 30)
        }
        break
      case 'comedy':
        result = 30000
        if (aPerformance.audience > 20) {
          result += 10000 + 500 * (aPerformance.audience - 20)
        }
        result += 300 * aPerformance.audience
        break
      default:
        throw new Error(`unknown type: ${playFor(aPerformance).type}`)
    }
    return result
  }
}
```

## 第 2 章 重构的原则

TODO

- 如果经理不能接受代码重构需要花时间，那就不要给经理说
- 何时重构：三次法则(事不过三，三则重构)

## 第 3 章 代码的坏味道

- 神秘命名（Mysterious Name）
  计算科学中最难的两件事是命名和缓存失效（There are only two hard things in computer science: cache invalidation and naming things）
- 重复代码（Duplicated Code）?
  有时候不一定，只是碰巧相似，但是未来可能会变得不同
- 过长函数（Long Function）
  有时候不一定，做好抽象即可
  `条件表达式和循环常常也是提炼的信号`
- 全局数据。（Global Data）(PS：慎用单例模式)
- 发散式变化（Divergent Change）
- 霰弹式修改
- 纯数据类（Data Class）

## 第 4 章 构筑测试体系

## 第 5 章 介绍重构名录

## 第 6 章 第一组重构

## 第 7 章 封装

## 第 8 章 搬移特性

## 第 9 章 重新组织数据

## 第 10 章 简化条件逻辑

## 第 11 章 重构 API

## 第 12 章 处理继承关系
