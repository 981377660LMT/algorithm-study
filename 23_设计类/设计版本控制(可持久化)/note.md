# 版本控制

1. 记录状态(state,持久化)

   - 不能全量保存,需要 `structual sharing` 的可持久化数据结构（Persistent data structure）
   - **每次修改历史版本总会返回一个新的数据结构(immutable ,不可变对象)**
   - 一般使用树型结构实现(包含链表,非线性结构)
   - 持久化数据结构可以带来许多好处，比如异常安全（Exception Safety）和并发性（Concurrency）。

2. 记录变化(action/operation)
   - mutation/action
   - undo redo 的对顶栈
   - 一些部分可持久化数据结构

## 完全可持久化栈

- 链表实现 (git)
- Path copying (路径复制,将即将被修改的点路径上的所有节点克隆出一个新的。这些修改必须通过数据结构逐级连接)

## 可持久化线段树(静态)

静态查询区间第 k 小

## 完全可持久化数组

### 完全可持久化数组有**在线离线**两种维护方法

N:数组长度 M:更新次数 Q:查询次数

- 离线:

  - 预处理查询+dfs(前提是操作可逆) `O(N+M+Q)`

- 在线:

  - full backup (Copy on Write 写入时复制,使用复制整个数据结构的方式来记录每次更改) `O(N*M+Q)`
  - 状态复元(前提是操作可逆) `O(N+M*Q)`
  - 数据结构 `O((Q+M)*logN)`

### 完全可持久化数组的 api:

1. 在某个历史版本上修改某一个位置上的值(生成一个完全一样的版本，不作任何改动;从 1 开始编号，版本 0 表示初始状态数组)
2. 访问某个历史版本上的某一位置的值

将每个 version 视为结点,那么所有的 version 连接构成了一棵树

https://www.luogu.com.cn/problem/P3919
https://qiita.com/wotsushi/items/72e7f8cdd674741ffd61#%E5%8F%82%E8%80%83%E8%A8%98%E4%BA%8B
https://www.cnblogs.com/Icys/p/Persistence_SegmentTree.html

py 库 https://github.com/tobgu/pyrsistent
js 库 https://github.com/immutable-js/immutable-js

## 完全可持久化并查集

用可持久化数组维护的并查集，本质和可持久化数组是一样的。
https://www.luogu.com.cn/blog/SSerxhs/solution-p3402

fa[a] = b，为了可持久化，我们就用可持久化数组来维护 fa[i]。注意这里不能再使用路径压缩了，道理很简单，可持久化要尽可能减少修改的次数。但是我们依然保留了一种优化方式：在维护 fa[i] 的同时维护一个 dep[i] ，表示这个节点的深度，保证在合并时是深度较小的点向深度较大的点合并即可。

## 完全可持久化平衡树

https://www.luogu.com.cn/problem/P3835

## 参考

[Persistent Data Structure](https://fuzhe1989.github.io/2017/11/07/persistent-data-structure/) (推荐)
[永続データ構造](https://qiita.com/wotsushi/items/72e7f8cdd674741ffd61#%E5%8F%82%E8%80%83%E8%A8%98%E4%BA%8B)
[持久化数据结构学习笔记——序列](https://zhuanlan.zhihu.com/p/33859991) (主要在说如何优化编码来节省空间)
[记录历史：持久化数据结构](https://quant67.com/post/algorithms/ads/persistent/persistent.html)
[可持久化数据结构](https://zh.m.wikipedia.org/zh-hans/%E5%8F%AF%E6%8C%81%E4%B9%85%E5%8C%96%E6%95%B0%E6%8D%AE%E7%BB%93%E6%9E%84)
[陈立杰：可持久化数据结构研究.pdf](https://github.com/Misaka233/algorithm/blob/master/%E9%99%88%E7%AB%8B%E6%9D%B0%EF%BC%9A%E5%8F%AF%E6%8C%81%E4%B9%85%E5%8C%96%E6%95%B0%E6%8D%AE%E7%BB%93%E6%9E%84%E7%A0%94%E7%A9%B6.pdf)
[Re 永続データ構造が分からない人のためのスライド](https://www.slideshare.net/qnighy/re-15443018)
[競技プログラミングにおける永続データ構造問題まとめ](https://blog.hamayanhamayan.com/entry/2017/05/21/001252)

## Persistent Data Structure 博客笔记

### 什么地方需要用到持久性数据结构？

1. 函数式编程语言。它的定义就要求了不能有可变数据和可变数据结构。
2. 使用 Persistent Map/HashMap 有助于简化 Prototype 的实现。
3. Lazy Evaluation。(如果一个数据结构是可变的，我们肯定不会放心对它使用 Lazy Evaluation。)
4. **并发编程(eg:协同编辑)**。
   并发编程同步手段:锁/原子操作

锁的缺点:

1.  锁的粒度太大，会导致性能下降。
2.  死锁

原子操作的缺点:

1. 原子操作本质上也是锁（总线锁），因此高并发度时开销还是会很大。
   无锁的策略使用一种叫做比较交换的技术（CAS Compare And Swap）来鉴别线程冲突，一旦检测到冲突产生，就重试当前操作直到没有冲突为止。

   > CAS (Compare And Swap)∶ 解决多线程并行情况下使用锁造成性能损耗的一种机制，这是硬件实现的原子操作。CAS 操作包含三个操作数:内存位置、预期原值和新值。如果内存位置的值与预期原值相匹配，那么处理器会自动将该位置值更新为新值。否则，处理器不做任何操作。

2. 需要非常精细的实现一个 lock-free 的数据结构，维护难度大，且很难证明其正确性。一些久经考验的数据结构仍然可能存在 bug。

在不同线程间访问相同的持久性数据结构时，我们很清楚其中不会有任何的数据竞争，因为无论其他线程如何修改这个结构，当前线程看到的结构永远是不变的。这不需要异常复杂的实现和同样复杂的测试来保证。通过锁和原子操作实现的并发数据结构追求的是“没有明显的错误”，而持久性数据结构则是“明显没有错误”。

### 如何实现 Persistent Data Structure

- Copy Anything，开销太大。
- 记录 action，读取时的开销非常大。一种改进方案是`每 SQRT 个修改操作创建一次数据副本`。
- **Path Copy**，当我们进行修改时，我们会复制路径上经过的 Node，直到最终修改发生的 Node。这里我们用“构造”代替了“修改”

### `垃圾回收`与`引用计数`

持久性数据结构的一个核心思想是为当前每个持有的人保留一个版本，即对于相同的数据可能同时存在多个版本。这样我们就需要有`垃圾回收机制`，对于每个版本的数据，在没有人持有之后回收掉。
通常**引用计数就足够了**。因为持久化数据结构有一个特点，因为各个版本形成的是一颗树，无环，不需要可达性分析，同时，这个特点对 Java 类的分代 GC 也是很有好处的。

### Persistent Stack (链表实现)

当我们拥有了一个持久性的 List 后，我们就可以在其上实现一套函数式的操作：

- fmap
- filter
- foldl/foldr
- forEach

### Persistent Map (可持久化平衡树)

算法导论思考题 13-1：持久动态集合中的持久二叉搜索树

### Persistent Vector (可持久化数组)

三个优化:

1. 优化 logn 到常数：使用类似于 Trie 树的结构，即每个节点能容纳多于两个子节点，**通常为 32 个(想象每个结点都是一个 bitset 状态压缩)**，这样一个 6 层的数据结构就可以容纳最多 1073741824 个元素，每次 修改 最多需要复制 **6 个节点**。
2. 懒标记：对于 32 个子节点的 Persistent Vector 来说，有 31/32 的 修改 是真正的 O(1)操作，其它 1/32 的 修改 才需要 O(lgn)。需要在每个 node 里记录 tail 结点的 offset。
3. Transient 优化：我们即使每个操作都复制了一个新节点，下次操作这个节点时还是要再复制，因为我们要保证不破坏其它人的使用。`那么什么时候可以直接重用这个节点，不用复制呢？`
   - 只能被当前 Vector 访问到的节点:我们每次修改 Vector 时都申请一个 UUID，并放到这次修改创建的节点中，这样通过 UUID 就能判断节点是不是当前 Vector 创建出来的。
   - 当前 Vector 在这次操作后不会被其它人使用:我们没有办法阻止用户**在所有修改完成之前使用这个 Vector**。因此我们使用一个新类型 Transient 来进行批量修改。好处：显式的让用户知道，我们要做一些不一样的事情啦，在它结束之前不要访问这个 Vector。在批量修改结束之后，还需要一个操作来**把 Transient 变为 Persistent**，这步操作会把每个节点中的 ID 置为 NULL，保证合法的 Vector 的节点没有 ID，从而避免一个 Transient 被误用。

### Persistent HashMap

当我们有了一个 Persistent Vector 的实现之后，实现一个 Persistent HashMap 也就不困难了。

- 如果我们使用开放寻址法来实现 HashMap，它就是单纯的一个 Vector，只要用 Persistent Vector 来实现就可以了。
- 如果我们使用链地址法来实现 HashMap，它就是一个 Persistent Vector，每个元素是一个 Persistent Stack，只要用 Persistent Vector 和 Persistent Stack 来实现就可以了。

### Persistent Data Structure 的一个使用场景

类似乐观锁

- RMW 操作，先 Read，再本地 Modify，再 Write 写回。对于 RMW 操作，必不可少的操作就是 CompareAndSwap(CAS)，也就是在 Write 时比较一下原对象有没有在你的 Read 之后被修改过。如果有的话，需要重试整个 RMW 操作。
- 可以使用 STM（Software Transactional Memory）来更规范的使用 UserMap：

  1. 线程创建 Transaction t。
  2. 线程通过 t 获得 M0。
  3. 线程通过 t 修改 M0，得到 M1。
  4. 线程通过 t 替换 UserMap 为 M1。
     这个过程中，如果有数据冲突，根据 STM 的实现不同，可能在不同的地方失败，只有当没有数据冲突时，整个操作才能顺利走下去。

- 当我们使用 Persistent Data Structure 时，**数据冲突的概率决定了它的使用是否高效**。上面的例子中，对于全局唯一的 UserMap，如果有大量的修改操作同时进行，那么其中只会有非常少量的操作能成功，其它操作都会因为数据冲突而失败。那么我们在使用时，就要考虑能否减少 UserMap 的粒度，从而降低冲突概率，提高性能。比如，**我们将 UserMap 分成 4096 个桶，每个桶是一个 Persistent HashMap，那么冲突概率就会小很多**，整体性能就上去了。类似 JDK 中的 ConcurrentHashMap。

https://zhuanlan.zhihu.com/p/350099474
`JDK1.7 中的 ConcurrentHashMap` ：segment+entry,Segment 继承了 RerentrantLock，所以 Segment 是一种可重入锁(可重入就是说某个线程已经获得某个锁，可以再次获取锁而不会出现死锁)，扮演锁的角色。Segment 默认为 16，也就是并发度为 16。
`JDK1.8 中的ConcurrentHashMap`：选择了与 HashMap 相同的 Node 数组+链表+红黑树结构；在锁的实现上，抛弃了原有的 Segment 分段锁，采用 CAS + synchronized 实现更加细粒度的锁。将锁的级别控制在了更细粒度的哈希桶数组元素级别，也就是说只需要`锁住这个链表头节点（红黑树的根节点），就不会影响其他的哈希桶数组元素的读写`，大大提高了并发度。

JDK1.7 与 JDK1.8 中 ConcurrentHashMap 的区别？

- 数据结构：取消了 Segment 分段锁的数据结构，取而代之的是数组+链表+红黑树的结构。
- 保证线程安全机制：JDK1.7 采用 Segment 的分段锁机制实现线程安全，其中 Segment 继承自 ReentrantLock 。JDK1.8 采用 CAS+synchronized 保证线程安全。
- 锁的粒度：JDK1.7 是对需要进行数据操作的 Segment 加锁，JDK1.8 调整为对每个数组元素加锁（Node）。
- 链表转化为红黑树：定位节点的 hash 算法简化会带来弊端，hash 冲突加剧，因此在链表节点数量大于 8（且数据总量大于等于 64）时，会将链表转化为红黑树进行存储。
- 查询时间复杂度：从 JDK1.7 的遍历链表 O(n)， JDK1.8 变成遍历红黑树 O(logN)。
