from typing import List, Tuple

MOD = int(1e9 + 7)

# 每一步 操作 中，如果 num1 >= num2 ，你必须用 num1 减 num2 ；否则，你必须用 num2 减 num1
class Solution:
    def countOperations(self, num1: int, num2: int) -> int:
        res = 0
        # O(log) 辗转相除法，直到出现0
        while num1 and num2:
            res += num2 // num1
            num1, num2 = num2 % num1, num1
        return res


print(Solution().countOperations(1, 2))
print(Solution().countOperations(10, 10))
print(Solution().countOperations(2, 3))
