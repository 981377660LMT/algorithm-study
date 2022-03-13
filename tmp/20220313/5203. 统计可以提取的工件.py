from typing import List, Tuple


MOD = int(1e9 + 7)


class Solution:
    def digArtifacts(self, n: int, artifacts: List[List[int]], dig: List[List[int]]) -> int:
        visited = set(tuple(d) for d in dig)
        res = 0
        for r1, c1, r2, c2 in artifacts:
            ok = True
            for r in range(r1, r2 + 1):
                for c in range(c1, c2 + 1):
                    if (r, c) not in visited:
                        ok = False
            if ok:
                res += 1
        return res

