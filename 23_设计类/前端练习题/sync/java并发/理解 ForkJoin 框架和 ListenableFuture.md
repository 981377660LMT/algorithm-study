# 理解 ForkJoin 框架和 ListenableFuture

## 1. RecursiveTask 和 RecursiveAction - Fork/Join 框架的基础

Fork/Join 框架是 Java 7 引入的用于并行执行任务的框架，专为可以递归分解的工作负载设计。它基于"分而治之"的原则：将大任务拆分成小任务并行执行，然后合并结果。

### RecursiveTask<V> - 有返回值的递归任务

`RecursiveTask` 是一个抽象类，用于需要返回结果的任务。

#### 关键特点

1. **返回值**：任务执行完成后返回计算结果
2. **递归分解**：大问题被分解为小问题
3. **结果合并**：子任务结果被合并成最终结果

#### 工作原理

```java
public abstract class RecursiveTask<V> extends ForkJoinTask<V> {
    protected abstract V compute();  // 子类必须实现此方法
}
```

使用 `RecursiveTask` 时，你需要：

1. 继承 `RecursiveTask` 类
2. 实现 `compute()` 方法，包含：
   - 问题足够小时直接解决
   - 问题较大时分解并fork子任务
   - 等待子任务完成并合并结果

#### 实际例子：并行计算数组和

```java
class SumTask extends RecursiveTask<Long> {
    private final long[] array;
    private final int start;
    private final int end;
    private static final int THRESHOLD = 10000;

    SumTask(long[] array, int start, int end) {
        this.array = array;
        this.start = start;
        this.end = end;
    }

    @Override
    protected Long compute() {
        // 问题足够小，直接计算
        if (end - start <= THRESHOLD) {
            long sum = 0;
            for (int i = start; i < end; i++)
                sum += array[i];
            return sum;
        }

        // 问题较大，分解为子问题
        int middle = (start + end) / 2;

        // 创建子任务
        SumTask leftTask = new SumTask(array, start, middle);
        SumTask rightTask = new SumTask(array, middle, end);

        // 异步执行其中一个子任务
        leftTask.fork();

        // 当前线程直接执行另一个子任务
        long rightResult = rightTask.compute();

        // 等待之前异步执行的子任务结果
        long leftResult = leftTask.join();

        // 合并结果
        return leftResult + rightResult;
    }
}
```

### RecursiveAction - 无返回值的递归任务

`RecursiveAction` 与 `RecursiveTask` 非常相似，但它不返回结果（void返回类型）。

#### 关键特点

1. **无返回值**：用于执行不需要返回结果的操作
2. **递归分解**：与 `RecursiveTask` 相同
3. **适用场景**：适用于原地修改数据或执行副作用操作

#### 工作原理

```java
public abstract class RecursiveAction extends ForkJoinTask<Void> {
    protected abstract void compute();  // 子类必须实现此方法
}
```

使用方式与 `RecursiveTask` 类似，但不需要处理返回值。

#### 实际例子：并行数组排序

```java
class SortTask extends RecursiveAction {
    private final int[] array;
    private final int start;
    private final int end;
    private static final int THRESHOLD = 1000;

    SortTask(int[] array, int start, int end) {
        this.array = array;
        this.start = start;
        this.end = end;
    }

    @Override
    protected void compute() {
        // 问题足够小，直接排序
        if (end - start <= THRESHOLD) {
            Arrays.sort(array, start, end);
            return;
        }

        // 问题较大，分解为子问题
        int middle = (start + end) / 2;

        // 创建并执行子任务
        invokeAll(
            new SortTask(array, start, middle),
            new SortTask(array, middle, end)
        );

        // 合并两个已排序的子数组
        merge(array, start, middle, end);
    }

    // 合并两个已排序数组的方法
    private void merge(int[] array, int start, int middle, int end) {
        // 合并逻辑...
    }
}
```

### RecursiveTask 与 RecursiveAction 的对比

| 特性     | RecursiveTask                    | RecursiveAction                     |
| -------- | -------------------------------- | ----------------------------------- |
| 返回值   | 有（泛型类型 V）                 | 无（void）                          |
| 用途     | 需要计算结果的任务               | 仅执行操作的任务                    |
| 合并阶段 | 需要合并子任务结果               | 无需合并结果                        |
| 方法签名 | `protected abstract V compute()` | `protected abstract void compute()` |
| 典型应用 | 计算、查找、收集结果             | 排序、修改数据、并行副作用操作      |

## 2. ListenableFuture - 可监听的异步计算结果

`ListenableFuture` 是 Google Guava 库引入的接口，它扩展了 Java 的 `Future` 接口，增加了监听器功能，使你可以注册回调以在 Future 完成时执行。

### 关键特点

1. **回调支持**：允许注册在 Future 完成时执行的回调函数
2. **链式操作**：可以组合多个异步操作，形成处理管道
3. **转换和组合**：提供丰富的 API 处理异步结果

### 工作原理

```java
// ListenableFuture 接口基本结构
public interface ListenableFuture<V> extends Future<V> {
    void addListener(Runnable listener, Executor executor);
}
```

使用 `ListenableFuture` 的主要方式：

1. **创建**：通常通过 `ListeningExecutorService` 创建
2. **添加回调**：使用 `Futures.addCallback` 添加成功/失败处理
3. **转换**：使用 `Futures.transform` 转换结果
4. **组合**：使用 `Futures.allAsList` 等方法组合多个 Future

### 实际例子

```java
// 创建 ListeningExecutorService
ListeningExecutorService service = MoreExecutors.listeningDecorator(
    Executors.newFixedThreadPool(10));

// 提交任务，获得 ListenableFuture
ListenableFuture<Integer> future = service.submit(() -> {
    // 复杂计算
    Thread.sleep(1000);
    return 42;
});

// 添加回调
Futures.addCallback(future, new FutureCallback<Integer>() {
    @Override
    public void onSuccess(Integer result) {
        System.out.println("计算结果: " + result);
    }

    @Override
    public void onFailure(Throwable t) {
        System.err.println("计算失败: " + t.getMessage());
    }
}, MoreExecutors.directExecutor());

// 转换结果
ListenableFuture<String> transformed = Futures.transform(future,
    (Integer input) -> "结果是: " + input,
    MoreExecutors.directExecutor());

// 组合多个 Future
ListenableFuture<List<Object>> combined = Futures.allAsList(future, transformed);
```

## 3. 三者关系与适用场景

### RecursiveTask 与 RecursiveAction

这两个类是 **Fork/Join 框架** 的核心，适用于：

- **CPU 密集型任务**：充分利用多核 CPU
- **可分解问题**：能够递归拆分的大规模计算
- **数据并行**：对大数据集执行相同操作

主要应用场景：

- 并行排序
- 矩阵和数组计算
- 图像处理
- 树和图的遍历

### ListenableFuture

它是 **异步编程模型** 的扩展，适用于：

- **I/O 密集型任务**：网络请求、数据库操作
- **事件驱动编程**：响应事件而非阻塞等待
- **复合异步操作**：依赖关系复杂的异步任务链

主要应用场景：

- 网络服务调用
- 异步 UI 更新
- 复杂的依赖链工作流
- 响应式编程模式

## 4. 类比理解

### RecursiveTask

想象你在组织一个大型拼图游戏。你可以：

1. 把拼图分成几个区域（分解问题）
2. 让不同人负责不同区域（并行执行）
3. 最后把完成的区域拼在一起（合并结果）
4. 得到完整的图像（返回结果）

### RecursiveAction

想象你在组织一次大规模清洁活动：

1. 把整个区域分成几个小区域（分解问题）
2. 派不同团队清洁不同区域（并行执行）
3. 每个团队直接清洁自己的区域（原地修改）
4. 不需要"合并"任何东西 - 工作完成后环境已被清洁（无返回值）

### ListenableFuture

想象你在餐厅点了一份需要时间准备的菜：

1. 下单后，你不会一直盯着厨房等待（非阻塞）
2. 你告诉服务员："菜准备好后，请通知我"（注册回调）
3. 你可以去做其他事情（执行其他代码）
4. 当服务员通知你菜准备好了，你可以回到餐桌（回调执行）
5. 你甚至可以提前告诉服务员："菜好了后，请直接加些胡椒再给我"（转换操作）

通过这些类比，这些看似复杂的并发工具就变得更容易理解了。它们各自解决不同的并发编程问题，选择合适的工具取决于你的具体需求。
