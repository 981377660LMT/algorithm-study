# A 相对于 B 的优势可以用满足 A[i] > B[i] 的索引 i 的数目来描述。
# 返回 A 的任意排列，使其相对于 B 的优势最大化。
# !田忌赛马
# 每次在A中寻找大于B[i]的最小值，若没有，则返回A中的最小值


from typing import List
from sortedcontainers import SortedList


class Solution:
    def advantageCount(self, A: List[int], B: List[int]) -> List[int]:
        sl = SortedList(A)
        res = []
        for num in B:
            pos = sl.bisect_right(num)
            if pos < len(sl):
                choose = sl.pop(pos)
            else:
                choose = sl.pop(0)
            res.append(choose)
        return res


print(Solution().advantageCount(A=[2, 7, 11, 15], B=[1, 10, 4, 11]))
# 输出：[2,11,7,15]
