# 1588. 所有奇数长度子数组的和

from typing import List


class Solution:
    def sumOddLengthSubarrays(self, arr: List[int]) -> int:
        """看每个元素会在多少个长度为奇数的数组中出现过。

        计算一个数字在多少个奇数长度的数组中出现过
        左侧有i+1个选择 其中偶数个数字的选择为 (i+1)+1//2 奇数个数字的选择为 (i+1)//2
        右侧有n-i个选择 其中偶数个数字的选择为 (n-i)+1//2 奇数个数字的选择为 (n-i)//2
        """
        res = 0
        for i, num in enumerate(arr):
            left, right = i + 1, len(arr) - i
            leftEven, leftOdd = (left + 1) // 2, left // 2
            rightEven, rightOdd = (right + 1) // 2, right // 2
            res += num * leftEven * rightEven + num * leftOdd * rightOdd
        return res
