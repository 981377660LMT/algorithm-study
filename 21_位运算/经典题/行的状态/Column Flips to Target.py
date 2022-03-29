# 可以互换行
# 问最少需要几次反转列的操作 使得原矩阵等于目标矩阵

# n, m ≤ 100
from typing import List


class Solution:
    def solve(self, matrix, target):
        # 讨论最后哪一行成为目标的第一行，从而异或得出flip的状态
        def compress(row: List[int]) -> int:
            return int(''.join(map(str, row)), 2)

        cur = list(map(compress, matrix))
        target = sorted(map(compress, target))

        res = int(1e20)
        for cand in cur:
            flip = cand ^ target[0]
            trans = sorted(row ^ flip for row in cur)
            if trans == target:
                res = min(res, bin(flip).count('1'))

        return res if res < int(1e20) else -1


print(Solution().solve(matrix=[[0, 0], [1, 0], [1, 1]], target=[[0, 1], [1, 0], [1, 1]]))

