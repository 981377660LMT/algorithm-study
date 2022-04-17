from functools import lru_cache
from typing import List, Tuple


MOD = int(1e9 + 7)
INF = int(1e20)


class Solution:
    def perfectMenu(
        self,
        materials: List[int],
        cookbooks: List[List[int]],
        attribute: List[List[int]],  # 美味度 x 和饱腹感 y。
        limit: int,
    ) -> int:
        def dfs(index, a, b, c, d, e, baofu, meiwei) -> None:
            nonlocal res
            if baofu >= limit:
                res = max(res, meiwei)
            if index == len(cookbooks):
                return
            # print(index, baofu, meiwei, a, b, c, d, e)

            dfs(index + 1, a, b, c, d, e, baofu, meiwei)

            ca, cb, cc, cd, ce = cookbooks[index]
            isOk = a >= ca and b >= cb and c >= cc and d >= cd and e >= ce
            if isOk:
                dfs(
                    index + 1,
                    a - ca,
                    b - cb,
                    c - cc,
                    d - cd,
                    e - ce,
                    baofu + attribute[index][1],
                    meiwei + attribute[index][0],
                )

        a, b, c, d, e = materials
        res = -1
        dfs(0, a, b, c, d, e, 0, 0)
        return res


print(
    Solution().perfectMenu(
        materials=[3, 2, 4, 1, 2],
        cookbooks=[[1, 1, 0, 1, 2], [2, 1, 4, 0, 0], [3, 2, 4, 1, 0]],
        attribute=[[3, 2], [2, 4], [7, 6]],
        limit=5,
    )
)
