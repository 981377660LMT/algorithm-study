# 如果盒子 x 需要放置在盒子 y 的顶部，那么盒子 y 竖直的四个侧面都 必须 与另一个盒子或墙相邻。
# 给你一个整数 n ，返回接触地面的盒子的 最少 可能数量。

# 1 <= n <= 10^9   =>  数学方法 根号n或者O(1)
# we find the the largest x such that x*(x+1)*(x+2)//6 <= n as the starting point for which there will be x*(x+1)//2 blocks at the bottom.

# https://leetcode-cn.com/problems/building-boxes/solution/cpython3-tan-xin-deng-chai-shu-lie-qiu-h-xuba/
# cur = 1 + 3 + 6 ... + i + (1 + 2 + 3 + ... + j)
# 我们现在尝试找到最小的j使得cur >= n。

# 复杂度根号n


class Solution:
    def minimumBoxes(self, n: int) -> int:
        base = 0
        count = 0
        # 如果带下层不越界(悲观锁的味道)
        while count + base * (base + 1) // 2 <= n:
            count += base * (base + 1) // 2
            base += 1
        base -= 1

        bottomCount = base * (base + 1) // 2

        add = 0
        while count + add * (add + 1) // 2 < n:
            add += 1

        return bottomCount + add


print(Solution().minimumBoxes(4))
# 输出：3
# 解释：上图是 3 个盒子的摆放位置。
# 这些盒子放在房间的一角，对应左侧位置。
