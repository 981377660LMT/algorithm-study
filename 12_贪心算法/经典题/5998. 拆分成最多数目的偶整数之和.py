from typing import List

# 给你一个整数 finalSum 。请你将它拆分成若干个 互不相同 的偶整数之和，且拆分出来的偶整数数目 最多 。
class Solution:
    def maximumEvenSplit(self, finalSum: int) -> List[int]:
        """从小到大使用所有的偶数，最后剩余部分不够大时，添加到上一个数中即可。"""
        if finalSum % 2 != 0:
            return []
        res = []
        cur = 2
        while finalSum > 0:
            if finalSum < cur:
                res[-1] += finalSum
                break
            else:
                res.append(cur)
                finalSum -= cur
                cur += 2
        return res

