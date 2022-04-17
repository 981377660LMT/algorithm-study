from typing import List


class Solution:
    def perfectMenu(
        self,
        materials: List[int],
        cookbooks: List[List[int]],
        attribute: List[List[int]],
        limit: int,
    ) -> int:
        def dfs(index, a, b, c, d, e, baofu, meiwei) -> None:
            nonlocal res

            if index == len(cookbooks):
                if baofu >= limit:
                    res = max(res, meiwei)
                return

            # jump
            dfs(index + 1, a, b, c, d, e, baofu, meiwei)

            # not to jump
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

