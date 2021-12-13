from typing import List

# 我们定义「顺次数」为：每一位上的数字都比前一位上的数字大 1 的整数。
# 请你返回由 [low, high] 范围内所有顺次数组成的 有序 列表（从小到大排序）。


class Solution:
    def sequentialDigits(self, low: int, high: int) -> List[int]:
        def gen(len: int) -> List[int]:
            res = []
            for i in range(1, 11 - len):
                charArray = list(map(str, list(range(i, i + len))))
                res.append(int(''.join(charArray)))
            return res

        res = []
        for i in range(len(str(low)), len(str(high)) + 1):
            res.extend(gen(i))

        return list(filter(lambda x: low <= x <= high, res))


print(Solution().sequentialDigits(low=100, high=300))
# 输出：[123,234]
