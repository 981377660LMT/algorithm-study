# TypeScript 实现 CountDownLatch

`CountDownLatch` 是一种同步工具，它允许一个或多个线程等待，直到在其他线程中执行的一组操作完成。这个实现在 TypeScript 中模拟了 Java 中的 `CountDownLatch` 功能。

### 核心属性

1. `_count`: 表示需要等待完成的操作数量
2. `_waiters`: 等待计数归零的线程队列

### 主要方法

1. **构造函数**:

   - 接受一个非负整数作为初始计数
   - 验证参数有效性，不允许负数计数

2. **countDown()**:

   - 将计数减一
   - 如果计数达到零，释放所有等待的线程
   - 如果计数已经是零，则不执行任何操作

3. **countDownMultiple(delta)**:

   - 一次将计数减少指定的值
   - 验证减少量是合法的（正数且不会使计数变为负数）
   - 如果计数达到零，释放所有等待的线程

4. **await(timeout?)**:

   - 异步方法，使当前线程等待，直到计数归零
   - 如果计数已经是零，则立即返回
   - 否则，创建一个Promise并加入等待队列
   - 支持可选的超时参数，超时后会拒绝Promise

5. **awaitWithBoolean(timeout)**:

   - 提供更便捷的超时处理，返回布尔值而不是抛出异常
   - 如果成功等待计数归零，返回true
   - 如果等待超时，返回false

6. **countDownAll()**:
   - 一次性将计数减少到零
   - 释放所有等待的线程

### 实现特点

1. **一次性使用**:

   - 与Java版本一样，这个CountDownLatch设计为一次性使用
   - 计数一旦达到零，无法重置

2. **异步/Promise支持**:

   - 完全基于Promise设计，适用于JavaScript/TypeScript的异步编程模型
   - 无缝集成async/await语法

3. **超时处理**:

   - await方法支持超时参数
   - 提供两种超时处理方式：抛出异常或返回布尔值

4. **扩展功能**:
   - countDownMultiple: 一次减少多个计数
   - countDownAll: 一次将计数归零
   - 这些是对Java版本的扩展，提供更多灵活性

### 使用场景

1. **启动信号**:

   - 主线程可以创建CountDownLatch(1)，并在准备就绪后调用countDown()
   - 工作线程在开始前调用await()，确保所有线程同时开始

2. **完成信号**:

   - 主线程创建CountDownLatch(N)，其中N是工作线程数
   - 每个工作线程完成时调用countDown()
   - 主线程调用await()等待所有工作完成

3. **资源初始化**:

   - 等待多个资源或服务初始化完成
   - 每个资源初始化完成时调用countDown()

4. **分段任务协调**:
   - 将大任务分解为多个部分并行处理
   - 使用CountDownLatch协调所有部分完成后的后续处理

### 与其他同步工具的比较

1. **vs CyclicBarrier**:

   - CountDownLatch: 一次性使用，计数不能重置
   - CyclicBarrier: 可循环使用，自动重置，适用于多次重复的同步点

2. **vs Semaphore**:

   - CountDownLatch: 计数只能减少，不能增加，主要用于等待操作完成
   - Semaphore: 计数可增可减，主要用于控制资源访问

3. **vs 条件变量**:
   - CountDownLatch: 简单的倒计时机制，适合一次性等待
   - 条件变量: 更通用的等待/通知机制，适合复杂条件和重复使用

这个TypeScript实现提供了Java CountDownLatch的全部功能，并添加了一些额外的便利方法，使其更适合JavaScript/TypeScript的异步编程模型。
