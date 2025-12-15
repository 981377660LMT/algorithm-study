from typing import List


class Solution:
    def minMoves(self, balance: List[int]) -> int:
        if sum(balance) < 0:
            return -1
        negIndex = next((i for i, x in enumerate(balance) if x < 0), None)
        if negIndex is None:
            return 0
        neg = -balance[negIndex]
        cands = []
        for i, v in enumerate(balance):
            if v > 0:
                d = abs(i - negIndex)
                d = min(d, len(balance) - d)
                cands.append((d, v))
        cands.sort()
        res = 0
        for d, v in cands:
            if neg <= 0:
                break
            cur = min(neg, v)
            res += cur * d
            neg -= cur
        return res
