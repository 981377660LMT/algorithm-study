from typing import List


class Solution:
    def addBinary1(self, arr1: List[int], arr2: List[int]) -> List[int]:
        res = []
        carry = 0
        i, j = len(arr1) - 1, len(arr2) - 1
        while i >= 0 or j >= 0:
            d1 = 0 if i < 0 else arr1[i]
            d2 = 0 if j < 0 else arr2[j]
            add = d1 + d2 + carry
            div, mod = divmod(add, 2)
            res.append(mod)
            carry = div
            i, j = i - 1, j - 1

        if carry > 0:
            res.append(carry)
        return res[::-1]

    def addBinary(self, arr1: List[int], arr2: List[int]) -> List[int]:
        res = []
        carry = 0
        i, j = len(arr1) - 1, len(arr2) - 1

        # carry逻辑放在这里面更好
        while i >= 0 or j >= 0 or carry:
            d1 = 0 if i < 0 else arr1[i]
            d2 = 0 if j < 0 else arr2[j]
            add = d1 + d2 + carry
            div, mod = divmod(add, 2)
            res.append(mod)
            carry = div
            i, j = i - 1, j - 1

        # if carry > 0:
        #     res.append(carry)
        return res[::-1]


print(Solution().addBinary(arr1=[1, 1, 1, 1, 1], arr2=[1, 0, 1]))
