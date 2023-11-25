# 2927. 给小朋友们分糖果 III
# https://leetcode.cn/problems/distribute-candies-among-children-iii/description/

# 你被给定两个正整数 n 和 limit。
# 返回 在每个孩子得到不超过 limit 个糖果的情况下，将 n 个糖果分发给 3 个孩子的 总方法数。
# n<=1e8 limit<=1e8


# !所有方案减去不合法的方案.


def C2(n: int) -> int:
    if n <= 1:
        return 0
    return n * (n - 1) // 2


class Solution:
    def distributeCandies(self, n: int, limit: int) -> int:
        res0 = C2(n + 2)  # 总方案数,x1+x2+x3=n的非负整数解的个数(隔板法)
        res1 = 3 * C2(n + 2 - (limit + 1))  # 至少一个小朋友分到的糖果超过 limit 的方案数
        res2 = 3 * C2(n + 2 - 2 * (limit + 1))  # 至少两个小朋友分到的糖果超过 limit
        res3 = C2(n + 2 - 3 * (limit + 1))  # 三个小朋友分到的糖果都超过 limit
        return res0 - res1 + res2 - res3
