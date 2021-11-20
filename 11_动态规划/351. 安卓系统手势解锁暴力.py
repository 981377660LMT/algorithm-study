# @param {number} m
# @param {number} n
# 1 <= m, n <= 9
# @return {number}
# 那么请你统计一下有多少种 不同且有效的解锁模式 ，是 至少 需要经过 m 个点，但是 不超过 n 个点的。
# 此题类似于哈密尔顿路径的解法:状态压缩
# 不压缩n! 压缩n^2*2^n
# 数据较小，可以直接生成所有可能的排列，判断是否有效即可。判断有效性可以利用二维坐标。

from itertools import permutations
from typing import Tuple


class Solution:
    def numberOfPatterns(self, m: int, n: int) -> int:
        # 跳过了中点
        def check(state: Tuple[int, ...]) -> bool:
            visited = {state[0]}
            for i in range(1, len(state)):
                x1, y1 = divmod(state[i], 3)
                x2, y2 = divmod(state[i - 1], 3)
                mid = ((x1 + x2) // 2) * 3 + (y1 + y2) // 2
                if ((x1 + x2) & 1 == 0) and ((y1 + y2) & 1 == 0) and (mid) not in visited:
                    return False
                visited.add(state[i])
            return True

        return sum(
            check(state) for count in range(m, n + 1) for state in permutations(range(9), count)
        )


print(Solution().numberOfPatterns(1, 1))
print(Solution().numberOfPatterns(1, 2))

