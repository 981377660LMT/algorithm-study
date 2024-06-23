cyclejs

# Rxjs 和 Promise 区别

- 核心思想: 数据响应式

Promise => 允诺
Rxjs => 由订阅/发布模式引出来

- 执行和响应

Promise 需要 then() 或 catch()执行，并且是一次性的。
Rxjs 数据的流出不取决于是否 subscribe()，并且可以多次响应

- 数据源头和消费

Promise 没有确切的数据消费者，每一个 then 都是数据消费者，同时也可能是数据源头，适合组装流程式（A 拿到数据处理，完了给 B，B 完了把处理后的数据给 C，以此类推）
Rxjs 则有明确的数据源头，以及数据消费者
