from typing import List, Tuple


Interval = Tuple[int, int]


class Solution:
    def solve(self, l0: List[Interval], l1: List[Interval]) -> List[Interval]:
        """双指针求两个区间数组的交集"""
        n1, n2 = len(l0), len(l1)
        res: List[Interval] = []
        left, right = 0, 0
        while left < n1 and right < n2:
            s1, e1, s2, e2 = *l0[left], *l1[right]
            # 相交
            if s1 <= e2 <= e1 or s2 <= e1 <= e2:
                # 尽量往内缩
                res.append((max(s1, s2), min(e1, e2)))
            if e1 < e2:
                left += 1
            else:
                right += 1
        return res

