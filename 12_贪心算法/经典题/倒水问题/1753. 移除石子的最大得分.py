# 每回合你都要从两个 不同的非空堆 中取出一颗石子，并在得分上加 1 分。
# 当存在 两个或更多 的空堆时，游戏停止。
# 给你三个整数 a 、b 和 c ，返回可以得到的 最大分数 。


class Solution:
    def maximumScore(self, a: int, b: int, c: int) -> int:
        """1 <= a, b, c <= 1e5"""
        max_ = max(a, b, c)
        sum_ = a + b + c
        restSum = sum_ - max_
        if max_ > restSum:
            return restSum
        else:
            return sum_ // 2


print(Solution().maximumScore(a=2, b=4, c=6))
# 输出：6
# 解释：石子起始状态是 (2, 4, 6) ，最优的一组操作是：
# - 从第一和第三堆取，石子状态现在是 (1, 4, 5)
# - 从第一和第三堆取，石子状态现在是 (0, 4, 4)
# - 从第二和第三堆取，石子状态现在是 (0, 3, 3)
# - 从第二和第三堆取，石子状态现在是 (0, 2, 2)
# - 从第二和第三堆取，石子状态现在是 (0, 1, 1)
# - 从第二和第三堆取，石子状态现在是 (0, 0, 0)
# 总分：6 分 。
