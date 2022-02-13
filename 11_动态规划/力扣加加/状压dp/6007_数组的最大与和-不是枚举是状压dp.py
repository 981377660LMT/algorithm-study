from functools import lru_cache
from typing import List, Tuple


def getDigit(num: int, /, *, index: int, radix: int):
    """返回 `radix` 进制下 `num` 的 `index` 位的数字，`index` 最低位(最右)为 0 """
    assert radix >= 2 and index >= 0
    prefix = num // pow(radix, index)
    return prefix % radix


# 1 <= numSlots <= 9
# 1 <= n <= 2 * numSlots
# 1 <= nums[i] <= 15
# 不是枚举子集，肯定超时；
# 需要状压dfs

# 1066. 校园自行车分配 II
# 1879. 两个数组最小的异或值之和

# 这道题尝试了枚举子集和回溯，没有想状压dp/dfs 以后看到数据量小的题应该尝试状压dfs
# 状压dfs时间复杂度是m*3^n，n<=9完全可行
# 周赛的时候只往枚举子集和回溯上面想，没有考虑状压dp/dfs，感觉挺遗憾的...
# 直接套状压dp/dfs的模板即可，时间复杂度即为dfs可能的状态数O(n*3^m)


# `带有两个维度的问题，考虑清除谁做index，state到底是哪个维度，复杂度合不合理`
# 需要全部看完的做index 另一个维度做state


class Solution(object):
    def maximumANDSum(self, nums: List[int], numSlots: int) -> int:
        n = len(nums)

        @lru_cache(None)
        def dfs(index: int, state: int) -> int:
            if index == n:
                return 0

            res = 0
            for pos in range(numSlots):
                mod = getDigit(state, index=pos, radix=3)
                if mod == 2:
                    continue
                res = max(res, (nums[index] & (pos + 1)) + dfs(index + 1, state + pow(3, pos)))
            return res

        return dfs(0, 0)


print(Solution().maximumANDSum(nums=[1, 2, 3, 4, 5, 6], numSlots=3))
print(Solution().maximumANDSum(nums=[1, 3, 10, 4, 7, 1], numSlots=9))
