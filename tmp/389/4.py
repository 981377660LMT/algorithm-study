from typing import List, Tuple, Optional
from collections import defaultdict, Counter
from sortedcontainers import SortedList


MOD = int(1e9 + 7)
INF = int(1e20)

# 给你一个下标从 0 开始的二进制数组 nums，其长度为 n ；另给你一个 正整数 k 以及一个 非负整数 maxChanges 。

# 灵茶山艾府在玩一个游戏，游戏的目标是让灵茶山艾府使用 最少 数量的 行动 次数从 nums 中拾起 k 个 1 。游戏开始时，灵茶山艾府可以选择数组 [0, n - 1] 范围内的任何索引index 站立。如果 nums[index] == 1 ，灵茶山艾府就会拾起一个 1 ，并且 nums[index] 变成0（这 不算 作一次行动）。之后，灵茶山艾府可以执行 任意数量 的 行动（包括零次），在每次行动中灵茶山艾府必须 恰好 执行以下动作之一：


# 选择任意一个下标 j != index 且满足 nums[j] == 0 ，然后将 nums[j] 设置为 1 。这个动作最多可以执行 maxChanges 次。
# 选择任意两个相邻的下标 x 和 y（|x - y| == 1）且满足 nums[x] == 1, nums[y] == 0 ，然后交换它们的值（将 nums[y] = 1 和 nums[x] = 0）。如果 y == index，在这次行动后灵茶山艾府拾起一个 1 ，并且 nums[y] 变成 0 。
# 返回灵茶山艾府拾起 恰好 k 个 1 所需的 最少 行动次数。


def min2(a: int, b: int) -> int:
    return a if a < b else b


class Solution:
    def minimumMoves(self, nums: List[int], k: int, maxChanges: int) -> int:
        def cal1(startIndex: int) -> Optional[int]:
            # 枚举起点，只把maxChanges用完就可以成功的情况
            ones = 0
            step = 0
            if nums[startIndex] == 1:
                ones += 1
            if ones >= k:
                return 0
            if startIndex - 1 >= 0 and nums[startIndex - 1] == 1:
                ones += 1
                step += 1
                if ones >= k:
                    return step
            if startIndex + 1 < n and nums[startIndex + 1] == 1:
                ones += 1
                step += 1
                if ones >= k:
                    return step
            diff = k - ones
            if diff > maxChanges:
                return None
            return step + diff * 2

        # 枚举起点，把maxChanges用完还不能成功
        def cal2(startIndex: int) -> int:
            # 先把maxChanges用完
            step = maxChanges * 2
            diff = k - maxChanges
            # diff个1到startIndex(包括startIndex自己的1)需要的最小总移动次数

        n = len(nums)
        res = INF
        for i in range(n):
            tmp = cal1(i)
            if tmp is not None:
                res = min2(res, tmp)
        if res != INF:
            return res
        for i in range(n):
            tmp = cal2(i)
            res = min2(res, tmp)
        return res
