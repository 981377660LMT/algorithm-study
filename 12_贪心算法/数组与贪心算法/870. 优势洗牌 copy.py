from typing import List
from sortedcontainers import SortedList

# A 相对于 B 的优势可以用满足 A[i] > B[i] 的索引 i 的数目来描述。
# 返回 A 的任意排列，使其相对于 B 的优势最大化。


class Solution:
    def advantageCount(self, A: List[int], B: List[int]) -> List[int]:
        num1 = SortedList(A)
        res = []
        for num in B:
            index = num1.bisect_right(num)
            if index < len(num1):
                choose = num1.pop(index)
            else:
                first = num1[0]
                num1.discard(first)
                choose = first
            res.append(choose)
        return res


print(Solution().advantageCount(A=[2, 7, 11, 15], B=[1, 10, 4, 11]))
# 输出：[2,11,7,15]
