from typing import List


class Solution:
    def addNegabinary(self, arr1: List[int], arr2: List[int]) -> List[int]:
        res = []
        carry = 0
        i, j = len(arr1) - 1, len(arr2) - 1
        while i >= 0 or j >= 0 or carry:
            d1 = 0 if i < 0 else arr1[i]
            d2 = 0 if j < 0 else arr2[j]
            add = d1 + d2 + carry
            div, mod = divmod(add, 2)
            res.append(mod)
            # 注意carry取负 与正二进制相加的区别就在这里
            carry = -div
            i, j = i - 1, j - 1

        # 去除前导0
        while len(res) > 1 and res[-1] == 0:
            res.pop()

        return res[::-1]


print(Solution().addNegabinary(arr1=[1, 1, 1, 1, 1], arr2=[1, 0, 1]))
# 输出：[1,0,0,0,0]
# 解释：arr1 表示 11，arr2 表示 5，输出表示 16 。

