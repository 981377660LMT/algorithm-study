from typing import List, Tuple, Optional
from collections import defaultdict, Counter
from sortedcontainers import SortedList


MOD = int(1e9 + 7)
INF = int(1e20)

# 给你一个下标从 0 开始的数组 nums ，它包含 非负 整数，且全部为 2 的幂，同时给你一个整数 target 。

# 一次操作中，你必须对数组做以下修改：

# 选择数组中一个元素 nums[i] ，满足 nums[i] > 1 。
# 将 nums[i] 从数组中删除。
# 在 nums 的 末尾 添加 两个 数，值都为 nums[i] / 2 。
# 你的目标是让 nums 的一个 子序列 的元素和等于 target ，请你返回达成这一目标的 最少操作次数 。如果无法得到这样的子序列，请你返回 -1 。


# 数组中一个 子序列 是通过删除原数组中一些元素，并且不改变剩余元素顺序得到的剩余数组。


class Solution:
    def minOperations(self, nums: List[int], target: int) -> int:
        if sum(nums) < target:
            return -1

        targetCounter = [0] * 40
        for i in range(40):
            if target & (1 << i):
                targetCounter[i] += 1

        curCounter = [0] * 40
        for v in nums:
            for i in range(40):
                if v & (1 << i):
                    curCounter[i] += 1

        def moveTo(i: int) -> int:
            """将>2^i的数移动到2^i,返回移动次数"""
            nextPos = -1
            for j in range(i + 1, 40):
                if curCounter[j]:
                    nextPos = j
                    break
            curPos = nextPos
            while curPos > i:
                curCounter[curPos] -= 1
                curCounter[curPos - 1] += 2
                curPos -= 1
            return nextPos - i

        # 填平补齐
        res = 0
        moveToLeft = False
        for i in range(39):
            moveToLeft |= curCounter[i] < targetCounter[i]
            if curCounter[i]:
                if targetCounter[i]:
                    over = curCounter[i] - targetCounter[i]
                    cand = over // 2
                    curCounter[i + 1] += cand
                    curCounter[i] -= cand * 2
                    if curCounter[i] > 0:
                        moveToLeft = False
                else:
                    over = curCounter[i]
                    cand = over // 2
                    curCounter[i + 1] += cand
                    curCounter[i] -= cand * 2
                    if curCounter[i] > 0:
                        moveToLeft = False

        # 倒着移动
        for i in range(39, -1, -1):
            if targetCounter[i] and not curCounter[i]:
                res += moveTo(i)
            curCounter[i] -= targetCounter[i]
        return res


# # nums = [1,2,8], target = 7
# print(Solution().minOperations([1, 2, 8], 7))
# # nums = [1,32,1,2], target = 12
# print(Solution().minOperations([1, 32, 1, 2], 12))
# # [16,128,32]
# # 1
# print(Solution().minOperations([16, 128, 32], 1))
# # yuqi1:4
# print(
#     Solution().minOperations(
#         [1, 1, 1, 1, 1, 1],
#         5,
#     )
# )
# [16,64,4,128]
# 6
# 预期3：
# print(Solution().minOperations([16, 64, 4, 128], 6))
# [8,1024,8388608,4,8,2097152,1024,1024,128,1073741824,4,4096,4,4,524288,65536,33554432,2097152,65536,65536,128,4,4,8,268435456,256,268435456,65536,33554432,4096,1073741824,1073741824,524288,8388608,33554432,4096,33554432,1024,1073741824,4,8,2097152,4]
# 43
# print(
#     Solution().minOperations(
#         [
#             8,
#             1024,
#             8388608,
#             4,
#             8,
#             2097152,
#             1024,
#             1024,
#             128,
#             1073741824,
#             4,
#             4096,
#             4,
#             4,
#             524288,
#             65536,
#             33554432,
#             2097152,
#             65536,
#             65536,
#             128,
#             4,
#             4,
#             8,
#             268435456,
#             256,
#             268435456,
#             65536,
#             33554432,
#             4096,
#             1073741824,
#             1073741824,
#             524288,
#             8388608,
#             33554432,
#             4096,
#             33554432,
#             1024,
#             1073741824,
#             4,
#             8,
#             2097152,
#             4,
#         ],
#         43,
#     )
# )
# [32,1024,131072,1073741824,4096,4096,262144,4096,4096,8192,4096,131072,256,262144,262144,262144,1073741824,4194304,131072,8192,4194304,16,4096,524288,1073741824,1073741824,16,1,1048576,64,4096,8,1073741824,4194304,134217728,512,65536,4096,1048576,4096,1073741824,2,8192,4096,256]
# 45
print(
    Solution().minOperations(
        [
            32,
            1024,
            131072,
            1073741824,
            4096,
            4096,
            262144,
            4096,
            4096,
            8192,
            4096,
            131072,
            256,
            262144,
            262144,
            262144,
            1073741824,
            4194304,
            131072,
            8192,
            4194304,
            16,
            4096,
            524288,
            1073741824,
            1073741824,
            16,
            1,
            1048576,
            64,
            4096,
            8,
            1073741824,
            4194304,
            134217728,
            512,
            65536,
            4096,
            1048576,
            4096,
            1073741824,
            2,
            8192,
            4096,
            256,
        ],
        45,
    )
)
