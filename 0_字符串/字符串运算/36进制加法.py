import string
from typing import List


BASE36 = string.digits + string.ascii_lowercase


class Solution:
    def add(self, arr1: List[str], arr2: List[str]) -> List[str]:
        def toInt(s: str) -> int:
            if "0" <= s <= "9":
                return int(s)
            return ord(s) - 97 + 10

        def toChar(n: int) -> str:
            return BASE36[n]

        res = []
        carry = 0
        i, j = len(arr1) - 1, len(arr2) - 1

        # carry逻辑放在这里面更好
        while i >= 0 or j >= 0 or carry:
            d1 = 0 if i < 0 else toInt(arr1[i])
            d2 = 0 if j < 0 else toInt(arr2[j])
            add = d1 + d2 + carry
            div, mod = divmod(add, 36)
            res.append(toChar(mod))
            carry = div
            i, j = i - 1, j - 1

        return res[::-1]


# '1b'+'2x'
print(Solution().add(["1", "b"], ["2", "x"]))
