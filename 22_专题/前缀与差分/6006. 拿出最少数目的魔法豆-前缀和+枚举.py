from typing import List, Tuple


class Solution:
    def minimumRemoval(self, beans: List[int]) -> int:
        """从每个袋子中 拿出 一些豆子（也可以 不拿出），使得剩下的 非空 袋子中（即 至少 还有 一颗 魔法豆的袋子）魔法豆的数目 相等"""
        beans.sort()
        res = sum(beans[i] - beans[0] for i in range(len(beans)))
        cur = res
        for i in range(1, len(beans)):
            cur += beans[i - 1]
            cur -= (beans[i] - beans[i - 1]) * (len(beans) - i)
            res = min(res, cur)

        return res


print(Solution().minimumRemoval(beans=[4, 1, 6, 5]))
print(Solution().minimumRemoval(beans=[2, 2, 3, 10]))
