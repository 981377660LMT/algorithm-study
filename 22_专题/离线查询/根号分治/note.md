https://blog.hamayanhamayan.com/entry/2017/04/12/180257
https://drken1215.hatenablog.com/archive/category/%E5%B9%B3%E6%96%B9%E5%88%86%E5%89%B2
https://drken1215.hatenablog.com/archive/category/%E3%82%AF%E3%82%A8%E3%83%AA%E3%81%AE%E5%B9%B3%E6%96%B9%E5%88%86%E5%89%B2

---

[分块的一些技巧](https://dpair.gitee.io/articles/block/)

[浅谈根号算法](https://ddosvoid.github.io/2020/10/18/%E6%B5%85%E8%B0%88%E6%A0%B9%E5%8F%B7%E7%AE%97%E6%B3%95/)

[浅 谈 分 块](https://www.luogu.com.cn/blog/220037/Sqrt1)

1. 强制在线，只有询问
   对于这种题目，分块算法一般的思想是预处理出块 i 到块 j​ 的答案，边角中间的块对边角的影响用其它数据结构来存储

   - [蒲公英](https://ddosvoid.github.io/2020/10/17/Luogu-P4168-Violet-%E8%92%B2%E5%85%AC%E8%8B%B1/)
   - [区间逆序对](https://ddosvoid.github.io/2020/10/18/bzoj-3744-Gty%E7%9A%84%E5%A6%B9%E5%AD%90%E5%BA%8F%E5%88%97/)

2. 由只有询问的题目拓展到包含修改
   这类题目一般在只有询问时可以做到提前处理好答案数组，这时候分块的作用就是使得在修改的时候可以在块内进行重构

   - [弹飞绵羊](https://ddosvoid.github.io/2020/10/17/Luogu-P3203-HNOI2010-%E5%BC%B9%E9%A3%9E%E7%BB%B5%E7%BE%8A/)
   - [哈希冲突](https://ddosvoid.github.io/2020/10/17/Luogu-P3396-%E5%93%88%E5%B8%8C%E5%86%B2%E7%AA%81/)

3. 对操作序列分块
   `对于根号个修改一同处理(batching)`它们对询问的贡献，对于每个询问块内的修改暴力算，块外的之前已经预处理好了

   - [CF 342E Xenia and Tree](https://ddosvoid.github.io/2021/04/21/CF-342E-Xenia-and-Tree/)
   - [记录修改操作，每有根号个修改就新建一个版本，查询时倒序查询修改记录](https://usaco.guide/plat/sqrt?lang=py)

4. 更加一般的数据结构题目(SqrtDecomposition 数据结构)
   思想类似于线段树，只不过有些东西线段树没法维护，所以只能用分块来操作

   - [CF 455D Serega and Fun](https://ddosvoid.github.io/2021/05/04/CF-455D-Serega-and-Fun/)

5. 自然根号
   - 有若干数的和为 n，则不同的数最多有 O(sqrt(n))个
   - ∑ai = n，现在有 O(n) 个**二元组**，每个二元组的代价是 `O(min(ai,aj))`，那么总代价是 `O(n*sqrt(n))`
6. 根号平衡
   不同种类的操作的个数不同的情况下来保证复杂度
   - O(1)区间查询 O(sqrt(n))单点修改
   - O(sqrt(n))区间查询 O(1)单点修改
   - O(sqrt(n))插入一个数，O(1)查询小于一个数的个数
     值域分块，维护前缀和即可
7. 根号分治
   首先对于这种题目，我们考虑设置一个阈值 limit 。
   对于> limit 的数据，会有一些性质，我们根据这个性质进行处理。
   对于< limit 的数据同样也会有一些性质，我们也根据这个性质特殊处理。
   而且这两部分若缺少了这些性质复杂度就难以保证，但把它分开了复杂度就对了。我们 limit 一般取 √n,但实际情况实际考虑。

---

https://github.com/EndlessCheng/codeforces-go/blob/bd8ed7c9523602007483193b3855fe6204f5349a/copypasta/sqrt_decomposition.go#L1
