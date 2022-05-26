from typing import List


class Solution:
    def missingRolls(self, rolls: List[int], mean: int, n: int) -> List[int]:
        total = mean * (n + len(rolls)) - sum(rolls)
        div, mod = divmod(total, n)

        # 多余的数，前面的人多分1个 => 爬虫任务分组
        if n <= total <= 6 * n:
            return [div + int(i < mod) for i in range(n)]
        return []


print(Solution().missingRolls(rolls=[1, 5, 6], mean=3, n=4))
# 输出：[2,3,2,2]
# 解释：所有 n + m 次投掷的平均值是 (1 + 5 + 6 + 2 + 3 + 2 + 2) / 7 = 3 。
