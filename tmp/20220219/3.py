from typing import List, Tuple

MOD = int(1e9 + 7)


class Solution:
    def maximumEvenSplit(self, finalSum: int) -> List[int]:
        init = finalSum
        """拆分成若干个 互不相同 的偶整数之和，且拆分出来的偶整数数目 最多"""
        if finalSum % 2 == 1:
            return []
        if finalSum == 2:
            return [2]
        if finalSum == 4:
            return [4]
        if finalSum == 6:
            return [2, 4]
        cur = 2
        res = []
        while finalSum >= cur:
            finalSum -= cur
            res.append(cur)
            if finalSum == cur:
                res[-1] *= 2
                break
            cur += 2
        # print(finalSum, cur)
        if sum(res) != init:
            res[-1] = init - sum(res[:-1])
        return res


print(Solution().maximumEvenSplit(12))
print(Solution().maximumEvenSplit(7))
print(Solution().maximumEvenSplit(28))
print(Solution().maximumEvenSplit(8))
print(Solution().maximumEvenSplit(10))
print(Solution().maximumEvenSplit(14))
