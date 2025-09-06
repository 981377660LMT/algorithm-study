# 给你一个长度为 n 的整数数组 order 和一个整数数组 friends。
# order 包含从 1 到 n 的每个整数，且 恰好出现一次 ，表示比赛中参赛者按照 完成顺序 的 ID。
# friends 包含你朋友们的 ID，按照 严格递增 的顺序排列。friends 中的每个 ID 都保证出现在 order 数组中。
# 请返回一个数组，包含你朋友们的 ID，按照他们的 完成顺序 排列。

from typing import List


class Solution:
    def recoverOrder(self, order: List[int], friends: List[int]) -> List[int]:
        s = set(friends)
        return [x for x in order if x in s]
