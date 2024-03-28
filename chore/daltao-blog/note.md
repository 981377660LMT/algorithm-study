daltao-blog

https://taodaling.github.io/blog/categories/

1. [统计学习方法](https://taodaling.github.io/blog/2022/08/18/%E7%BB%9F%E8%AE%A1%E5%AD%A6%E4%B9%A0%E6%96%B9%E6%B3%95/)
2. [一些可能有点用的算法](https://taodaling.github.io/blog/2022/12/17/%E4%B8%80%E4%BA%9B%E5%8F%AF%E8%83%BD%E6%9C%89%E7%82%B9%E7%94%A8%E7%9A%84%E7%AE%97%E6%B3%95/)

   - 组合数随机选取问题
     `random.sample(arr, k): 按照 k 的大小采取两种解法-> k大，洗牌算法；k小，维护大小为k的集合.`
     时间复杂度`O(k)`
   - 红包算法(微信红包)
     n 元红包分给 k 个人
     先将金额转化为分(`*100`)，然后变成小球与隔板问题：`有N+K−1个不同的球，从中选择K−1个不同的球，每种选法的概率相等。`
     这个问题就是 `random.sample` 问题.
   - 桌游 rating 算法
     - `Elo rating` 算法
     - ∑wixi (wi = 1/(i+T)；xi:胜利为 1，失败为-1) 对所有比赛进行加权求和
   - 数值信息压缩问题
     cost(l,r)表示将灰度 l 到 r 映射到同一颜色的最少惩罚值
     dp[i][j]表示前 j 小的颜色用 i 个不同的颜色表示的最小惩罚值
     可以用决策单调性来优化

3. [实现一个现代编译器](https://taodaling.github.io/blog/2022/04/09/%E5%AE%9E%E7%8E%B0%E4%B8%80%E4%B8%AA%E7%8E%B0%E4%BB%A3%E7%BC%96%E8%AF%91%E5%99%A8/)
4. [Rust 学习笔记](https://taodaling.github.io/blog/2021/11/21/rust%E5%AD%A6%E4%B9%A0%E7%AC%94%E8%AE%B0/)
5. 环上匹配问题、环上连通问题 (atcoder 某次的 d 题?)
   https://taodaling.github.io/blog/2021/07/26/%E4%B8%80%E4%BA%9B%E5%9C%86%E7%8E%AF%E9%97%AE%E9%A2%98/
6. **一类撤销问题**
   https://taodaling.github.io/blog/2020/10/11/%E4%B8%80%E7%B1%BB%E6%92%A4%E9%94%80%E9%97%AE%E9%A2%98/
7. 正则表达式
   https://taodaling.github.io/blog/2020/08/04/%E6%AD%A3%E5%88%99%E8%A1%A8%E8%BE%BE%E5%BC%8F/
8. Java 的基础运算性能
   https://taodaling.github.io/blog/2020/09/22/java-competitive-program/

   执行操作 1e8 次

   ```
   plusMultiway finished in 38ms (循环展开实现的加法)
   empty finished in 46ms
   and finished in 47ms
   subtract finished in 48ms
   plus finished in 50ms
   longPlus finished in 80ms
   mul finished in 106ms
   longMul finished in 109ms
   invoke finished in 149ms
   choose finished in 187ms
   doubleMul finished in 279ms
   doublePlus finished in 284ms
   doubleDiv finished in 462ms
   div finished in 896ms
   mod finished in 919ms
   ```

   empty 是循环体带来的时间，所有时间减去 empty 的时间就是其方法体的时间。
   可以发现在 java 中加减和位运算非常快。
   长整数会使得运算时间翻倍，但是依旧在可以接受的范围中。
   函数调用和分支跳转更慢，但是依旧不会影响你的程序的总体时间。所以爱咋用就咋用。
   浮点数的运算时间非常慢，比函数调用还慢。但是浮点数的除法却比整数除法快很多，可能是算法不同。
   **整数除法和取模异常慢**，所以尽量避免。在取模运算的时候用 long 来存储累积和，只有在发生除法的时候才进行必要的取模。

9. 浮点数加法精度误差
   https://taodaling.github.io/blog/2020/06/16/%E6%B5%AE%E7%82%B9%E6%95%B0%E5%8A%A0%E6%B3%95%E7%B2%BE%E5%BA%A6%E8%AF%AF%E5%B7%AE/
   https://oi-wiki.org/misc/kahan-summation/
   在浮点加法计算中，交换律（commutativity）成立，但结合律（associativity）不成立。
   计算机在执行浮点数加法的时候，由于 double 类型只有 52 个有效位，因此精度是有限的（不到 16 位十进制有效数字），在执行加法运算的时候，较低的位会被省略掉，比如 1e10+(1+1e−10)，其结果为 1e10+1，而 1e−10 就被省略了。虽然看起来的省略项微不足道，但是如果有许多这样的项合在一起就变成了明显的精度误差了。
   `kahan summation` 算法可以用于一些计算浮点数的和，同时其造成的误差是 O(1)，即与加总的次数无关
   其具体代码如下：

   ```java

    public static double sum(double[] arr, int l, int r) {
        double sum = 0;
        double err = 0;  // 通过一个单独变量用来累积误差
        for (int i = l; i < r; i++) {
            double x = arr[i] - err;
            double t = sum + x;
            err = (t - sum) - x;
            sum = t;
        }
        return sum;
    }
   ```

   java 中利用 DoubleStream#sum 来计算浮点数和

10. Binary Lifting
    https://taodaling.github.io/blog/2020/03/18/binary-lifting/

    - 倍增技术：如果仅仅为了二分的话，倍增可以做到 O(n)的预处理空间复杂度和处理时间复杂度，参考 https://codeforces.com/blog/entry/74847
    - 倍增结构：
      倍增技术的强大是基于一个很简单的倍增结构。这个结构实际上还有额外的用处。
      记 jump(u,i)，它表示 u 沿着出边移动 2^i 步所在的位置，link(u,i)表示从 u 出发，通过少于 i 次的 next 转移能抵达的所有结点的集合。将 jump(u,i)视作一个结点，它覆盖了所有从 link(u,2^i)上的结点，称 i 为这个结点的高度。
      每一个倍增结点可以包含 2^i 个结点，其中 i 为倍增的次数。
      (u,v) 对应的区间可以转化为 O(logn)个倍增结点.

      现在有这样的一些问题，我们需要处理若干个请求，每个请求要求修改路径 link(u,l)上的所有结点。在所有请求完成后，要求输出所有结点的权值。
      我们实际上可以发现 u,v 对应的区间可以截断为 O(log2n)个倍增结构上的结点，我们只需要在这些结点上打上标记就可以了。并且考虑到标记只需要从高度较大的结点下推到高度较小的结点，因此在最后阶段我们可以从高到低处理结点。
      当然这个问题平平无奇，LCT 数据结构也可以做到 O(nlog2n)，但是倍增结构可以处理存在环的情况，而 LCT 就需要比较复杂的特殊处理。

      1. 并行并查集问题：给定 n 个结点，之后有 q 个请求。每个请求给定两端等长区间(a,b)和(c,d) ，表示对于所有 0≤i≤b−a，结点 a+i 和 c+i 在同一个连通块中。接下来要求计算最多可能存在多少个连通块。
         https://www.luogu.com.cn/problem/solution/P3295
      2. 树上最大流问题：现在有 n 个人，以及一颗大小为 m 的树。第 i 个人可以居住在 ui 和 vi 之间的路径上的任意一个顶点中，且一个顶点最多居住一个人。现在希望让尽可能多人居住在树上，问最多有多少人可以居住在树上。
      3. 倍增优化建图
         https://www.luogu.com.cn/problem/P5344
         类比线段树拆区间，倍增结构拆倍增链上的一段路径

      TODO: **DivideFunctionalGraph**

11. 二分中的相对和绝对误差
    https://taodaling.github.io/blog/2019/09/12/%E4%BA%8C%E5%88%86/
    现在很多输出浮点数的题目都会提供两种 AC 条件，一种是输出与真实结果的绝对误差不超过阈值，一种是输出与真实结果的相对误差不超过阈值。
    相对误差在输出结果很大的时候会发挥巨大的作用。众所周知双精度浮点型共 64 位，其中 1 位用于表示符号，11 位表示指数，其余的 52 位用于表示有效数字。简单换算就可以知道双精度浮点型可以精确表示大概 15 位十进制整数。
    现在考虑一个问题，最终结果为 10^8 ，但是要求绝对误差小于 10^−8，这现实吗。事实上尾部的数值由于有效数值不足会被舍去。这就会导致二分的时候，(l+r)/2 可能会等于 l 或等于 r，从而导致二分进入死循环。
    但是有了相对误差，情况就会大为不同，当输出为 10^8 时，我们可以不需要保留任意小数。
12. 缝隙二分
    二分还能作用在非单调函数上查找缝隙。
    所谓的缝隙是指这样一个整数 x，满足 check(x−1)≠check(x)。
    要执行缝隙二分的前提是，一开始给定的 l 和 r 满足 check(l)≠check(r)。

    在执行缝隙而二分的过程，利用 l 和 r 算出 mid=(l+r)/2 后，如果 check(mid)=check(l)，那么就令 l=mid，否则令 r=mid。
    容易发现这样做始终能保证 check(l)≠check(r)。
    由于区间在不断缩小，因此最终一定能找到一个缝隙。当然我们无法确定找到的是哪个缝隙，但是这不重要。

    https://atcoder.jp/contests/ddcc2020-qual/tasks/ddcc2020_qual_e

13. trie 合并
    https://www.luogu.com.cn/problem/CF778C
14. lca
    - 判断某个顶点 x 是否落在 u,v 的唯一路径上：首先 x 一定处在 lca(u,v)子树内，且 x 一定是 u 或 v 的祖先。
    - 判断 x1,y1 的路径与 x2,y2 的路径是否有交点：
      两条路径假如有交点，记任意交点为 z，很显然 z 一定是 x1 或 y1 的祖先，且 z 还是 x2 或 y2 的祖先。
      因此可以推出 lca(x2,y2)一定是 x1 或 y1 的祖先，且 lca(x1,y1)还是 x2 或 y2 的祖先。
