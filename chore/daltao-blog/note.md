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
