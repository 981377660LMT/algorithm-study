https://www.kancloud.cn/aceld/golang/1958308
Golang 三色标记+混合写屏障 GC 模式全分析

1. Go V1.3 之前的标记-清除(mark and sweep)算法
   有一点需要额外注意：mark and sweep 算法在执行的时候，需要程序暂停！即 **STW(stop the world)**CPU 不执行用户代码，全部用于垃圾回收，
   缺点：

   1. STW，stop the world；**让程序暂停**，程序出现卡顿 (重要问题)；
   2. 标记需要扫描整个 heap；
   3. 清除数据会产生 heap 碎片。

      开 stw=>标记=>清除=>关 stw
      优化成
      开 stw=>标记=>关 stw=>清除
      **但 单一的 mark and sweep 还是太慢**

2. Go V1.5 的三色并发标记法(注意是 bfs)
   GC 过程和其他用户 goroutine 可并发运行，但需要一定时间的 STW(stop the world)
   **综合来说就是树遍历中实际存在三种状态：已经遍历完的，遍历中的，没遍历的**
   为了在 GC 过程中保证数据的安全，我们在开始三色标记之前就会加上 STW，在扫描确定黑白对象之后再放开 STW。但是很明显这样的 GC 扫描的性能实在是太低了。

3. 没有 STW 的三色标记法
   如何三色标记过程不启动 STW，那么在 GC 扫描过程中，任意的对象均可能发生读写操作
   原来灰色 G 指向白色 W
   之后 G 断开指向 W
   另一个黑色 B 指向了 W
   此时 W 之后的对象都会被误判成垃圾被回收

   可以看出，有两种情况，在三色标记法中，是不希望被发生的。

   **条件 1**: 一个白色对象被黑色对象引用**(白色被挂在黑色下)**
   **条件 2**: 灰色对象与它之间的可达关系的白色对象遭到破坏**(灰色同时丢了该白色)**
   如果当以上两个条件同时满足时，就会出现对象丢失现象!
   为了防止这种现象的发生，最简单的方式就是 STW，直接禁止掉其他用户程序对对象引用关系的干扰，但是 STW 的过程有明显的资源浪费，对所有的用户程序都有很大影响。
   那么是否可以在保证对象不丢失的情况下合理的尽可能的提高 GC 效率，减少 STW 时间呢？答案是可以的，**我们只要使用一种机制，尝试去破坏上面的两个必要条件就可以了。**

4. 屏障机制(hook/handler/回调)
   `注意，所有的屏障机制都是只在堆上启用的`

   1. **强三色不变式**：不存在黑色对象引用到白色对象的指针，强制性的不允许黑色对象引用白色对象，这样就不会出现有白色对象被误删的情况。
   2. **弱三色不变式**：黑色对象可以引用白色对象，但是所有被黑色对象引用的白色对象都要处于灰色保护状态。(等待遍历状态)

   为了遵循上述的两个方式，GC 算法演进到两种屏障方式，他们“插入写屏障”, “删除写屏障”。

   1. 插入写屏障(堆空间)：在 A 对象引用 B 对象的时候，B 对象被标记为灰色。(将 B 挂在 A 下游，B 必须被标记为灰色)
      满足: 强三色不变式. (不存在黑色对象引用白色对象的情况了， 因为白色会强制变成灰色)

   ```Go
   添加下游对象(当前下游对象slot, 新下游对象ptr) {
        //1
        先标记灰色(新下游对象ptr)   // 写屏障

        //2
        当前下游对象slot = 新下游对象ptr
   }
   ```

   注意：黑色对象的内存槽有两种位置, 栈和堆. 栈空间的特点是容量小,但是要求相应速度快,因为函数调用弹出频繁使用, **所以“插入屏障”机制,在栈空间的对象操作中不使用. 而仅仅使用在堆空间对象的操作中.**
   但是如果栈不添加,当全部三色标记扫描之后,栈上有可能依然存在白色对象被引用的情况(如上图的对象 9). **所以要对栈重新进行三色标记扫描, 但这次为了对象不丢失, 要对本次标记扫描启动 STW 暂停. 直到栈空间的三色标记结束.**
   最后将栈和堆空间 扫描剩余的全部 白色节点清除. 这次 STW 大约的时间在 10~100ms 间.

   2. 删除写屏障:被删除的对象，如果自身为灰色或者白色，那么被标记为灰色。满足: 弱三色不变式. (保护灰色对象到白色对象的路径不会断)

   ```Go
   添加下游对象(当前下游对象slot， 新下游对象ptr) {
        //1
        if (当前下游对象slot是灰色 || 当前下游对象slot是白色) {
                标记灰色(当前下游对象slot)     //slot为被删除对象， 标记为灰色
        }

        //2
        当前下游对象slot = 新下游对象ptr
   }
   ```

   缺点：这种方式的回收精度低，一个对象即使被删除了最后一个指向它的指针也依旧可以活过这一轮，在下一轮 GC 中被清理掉。

5. Go V1.8 的混合写屏障(hybrid write barrier)机制
   插入写屏障和删除写屏障的短板：
   插入写屏障：**结束时需要 STW 来重新扫描栈**，标记栈上引用的白色对象的存活；
   删除写屏障：回收精度低，**一个对象即使被删除了最后一个指向它的指针也依旧可以活过这一轮**，在下一轮 GC 中被清理掉。

   Go V1.8 版本引入了混合写屏障机制（hybrid write barrier），避免了对栈 re-scan 的过程，极大的减少了 STW 的时间。结合了两者的优点。

   混合写屏障规则
   具体操作:

   1. GC 开始将栈上的可达对象全部扫描并标记为黑色(之后不再进行第二次重复扫描，无需 STW)，
   2. GC 期间，任何在栈上创建的新对象，均为黑色。
   3. 被删除的对象标记为灰色。
   4. 被添加的对象标记为灰色。

   这里我们注意， **屏障技术是不在栈上应用的**，因为要保证栈的运行效率。

   ```Go
   添加下游对象(当前下游对象slot, 新下游对象ptr) {
       //1
       标记灰色(当前下游对象slot)    //只要当前下游对象被移走，就标记灰色

       //2
       标记灰色(新下游对象ptr)

       //3
       当前下游对象slot = 新下游对象ptr
   }
   ```

   ​ Golang 中的混合写屏障满足弱三色不变式，结合了删除写屏障和插入写屏障的优点，**只需要在开始时并发扫描各个 goroutine 的栈，使其变黑并一直保持**，这个过程不需要 STW，而标记结束后，**因为栈在扫描后始终是黑色的，也无需再进行 re-scan 操作了**，减少了 STW 的时间。

6. 总结
   GoV1.3- 普通标记清除法，整体过程需要启动 STW，效率极低。
   GoV1.5- 三色标记法， 堆空间启动写屏障，栈空间不启动，全部扫描之后，需要重新扫描一次栈(需要 STW)，效率普通
   GoV1.8-三色标记法，混合写屏障机制， 栈空间不启动，堆空间启动。在开始时并发扫描各个 goroutine 的栈，使可达对象变黑并一直保持。整个过程几乎不需要 STW，效率较高。
