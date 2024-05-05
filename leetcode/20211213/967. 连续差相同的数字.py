from typing import List

# 返回所有长度为 n 且满足其每两个连续位上的数字之间的差的绝对值为 k 的 非负整数 。
# 2 <= n <= 9
# 0 <= k <= 9
class Solution:
    def numsSameConsecDiff(self, n: int, k: int) -> List[int]:
        res = []

        def bt(remain: int, pre: int, cur: int) -> None:
            if remain == 0:
                res.append(cur)
                return

            if k == 0:
                bt(remain - 1, pre, cur * 10 + pre)
            else:
                for delta in (k, -k):
                    choice = pre + delta
                    if choice < 0 or choice >= 10:
                        continue
                    bt(remain - 1, choice, cur * 10 + choice)

        for start in range(1, 10):
            bt(n - 1, start, start)

        return res


print(Solution().numsSameConsecDiff(n=2, k=1))
# 输出：[10,12,21,23,32,34,43,45,54,56,65,67,76,78,87,89,98]
print(Solution().numsSameConsecDiff(n=3, k=7))
