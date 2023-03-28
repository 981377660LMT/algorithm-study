from typing import List
from collections import defaultdict


# 给你一个下标从 0 开始的数组 nums ，它包含若干正整数，
# 表示数轴上你需要摧毁的目标所在的位置。同时给你一个整数 space 。

# 你有一台机器可以摧毁目标。给机器 输入 nums[i] ，
# !这台机器会摧毁所有位置在 nums[i] + c * space 的目标，
# 其中 c 是任意非负整数。你想摧毁 nums 中 尽可能多 的目标。

# 请你返回在摧毁数目最多的前提下，nums[i] 的 最小值 。

# !nums[j] = nums[i] + c * k 等价于 nums[j] % k = nums[i] % k

INF = int(1e20)


class Solution:
    def destroyTargets(self, nums: List[int], space: int) -> int:
        group = defaultdict(list)
        for num in nums:
            group[num % space].append(num)
        max_, res = 0, INF
        for g in group.values():
            if len(g) > max_:
                max_ = len(g)
                res = min(g)
            elif len(g) == max_:
                res = min(res, min(g))
        return res


print(Solution().destroyTargets([1, 2, 3, 4, 5], 3))
print(Solution().destroyTargets([1, 5, 3, 2, 2], 10000))
