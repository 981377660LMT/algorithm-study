from typing import List

# 1 <= arr.length <= 10^5

# 题目中的中位数与平时的中位数不同，题目中指的就是位于 ((n - 1) / 2) 的元素。
class Solution:
    def getStrongest(self, arr: List[int], k: int) -> List[int]:
        mid = sorted(arr)[(len(arr) - 1) // 2]
        return sorted(arr, reverse=True, key=lambda x: (abs(x - mid), x))[:k]


print(Solution().getStrongest(arr=[1, 2, 3, 4, 5], k=2))
# 输出：[5,1]
# 解释：中位数为 3，按从强到弱顺序排序后，数组变为 [5,1,4,2,3]。最强的两个元素是 [5, 1]。[1, 5] 也是正确答案。
# 注意，尽管 |5 - 3| == |1 - 3| ，但是 5 比 1 更强，因为 5 > 1 。

