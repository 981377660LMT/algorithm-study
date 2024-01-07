# 100154. 执行操作后的最大分割数量
# https://leetcode.cn/problems/maximize-the-number-of-partitions-after-operations/description/
# https://leetcode.cn/circle/discuss/HKYh6a/

# 给你一个下标从 0 开始的字符串 s 和一个整数 k。
# 你需要执行以下分割操作，直到字符串 s 变为 空：
# 选择 s 的最长前缀，该前缀最多包含 k 个 不同 字符。
# 删除 这个前缀，并将分割数量加一。如果有剩余字符，它们在 s 中保持原来的顺序。
# 执行操作之 前 ，你可以将 s 中 至多一处 下标的对应字符更改为另一个小写英文字母。
# 在最优选择情形下改变至多一处下标对应字符后，用整数表示并返回操作结束时得到的最大分割数量。

# 1 <= s.length <= 1e4
# s 只包含小写英文字母。
# 1 <= k <= 26
# O(2*26*26*n)


# 通过一个状压mask来记录每一段出现不同字符数，
# 当mask中的字符个数超过k个时，
# 就需要进行分割，然后在每一段枚举修改与不修改即可.
# !dfs(index, changed, mask) 表示当前在第index个位置，是否修改过，mask为当前看过的字符集合
# !原理同logTrick，mask 单调不减，最多26种状态.

from functools import lru_cache


def max2(a: int, b: int) -> int:
    return a if a > b else b


class Solution:
    def maxPartitionsAfterOperations(self, s: str, k: int) -> int:
        @lru_cache(None)
        def dfs(index: int, changed: bool, mask: int) -> int:
            if index == n:
                return 0

            cur = nums[index]
            newMask = mask | (1 << cur)
            res = 0

            # !不修改
            if newMask.bit_count() > k:
                res = 1 + dfs(index + 1, changed, 1 << cur)
            else:
                res = dfs(index + 1, changed, newMask)

            # !修改
            if not changed:
                for v in range(26):
                    if v == cur:
                        continue
                    newMask = mask | (1 << v)
                    if newMask.bit_count() > k:
                        res = max2(res, 1 + dfs(index + 1, True, 1 << v))
                    else:
                        res = max2(res, dfs(index + 1, True, newMask))

            return res

        nums = [ord(c) - 97 for c in s]
        n = len(nums)
        res = dfs(0, False, 0) + 1
        dfs.cache_clear()
        return res
