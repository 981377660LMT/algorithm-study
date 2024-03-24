from typing import List, Tuple, Optional
from collections import defaultdict, Counter
from sortedcontainers import SortedList


MOD = int(1e9 + 7)
INF = int(1e20)
# 你需要在一个集合里动态记录 ID 的出现频率。给你两个长度都为 n 的整数数组 nums 和 freq ，nums 中每一个元素表示一个 ID ，对应的 freq 中的元素表示这个 ID 在集合中此次操作后需要增加或者减少的数目。

# 增加 ID 的数目：如果 freq[i] 是正数，那么 freq[i] 个 ID 为 nums[i] 的元素在第 i 步操作后会添加到集合中。
# 减少 ID 的数目：如果 freq[i] 是负数，那么 -freq[i] 个 ID 为 nums[i] 的元素在第 i 步操作后会从集合中删除。
# 请你返回一个长度为 n 的数组 ans ，其中 ans[i] 表示第 i 步操作后出现频率最高的 ID 数目 ，如果在某次操作后集合为空，那么 ans[i] 为 0 。


class Solution:
    def mostFrequentIDs(self, nums: List[int], freq: List[int]) -> List[int]:
        sl = SortedList()
        counter = defaultdict(int)
        res = [0] * len(nums)
        for i, (v, f) in enumerate(zip(nums, freq)):
            pre = counter[v]
            counter[v] += f
            cur = counter[v]
            if f > 0:
                if pre > 0:
                    sl.discard(pre)
                sl.add(cur)
            else:
                sl.discard(pre)
                if cur > 0:
                    sl.add(cur)
            if sl:
                res[i] = sl[-1]
        return res
