from typing import List

# 给你两个整数 n 和 start。你的任务是返回任意 (0,1,2,,...,2^n-1) 的排列 p
class Solution:
    def circularPermutation(self, n: int, start: int) -> List[int]:
        return [start ^ i ^ (i >> 1) for i in range(2 ** n)]


print(Solution().circularPermutation(n=3, start=2))
# 输出：[2,6,7,5,4,0,1,3]
# 解释：这个排列的二进制表示是 (010,110,111,101,100,000,001,011)
