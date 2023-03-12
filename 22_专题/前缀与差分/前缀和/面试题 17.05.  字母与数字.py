from collections import defaultdict
from typing import List

# !给定一个放有字母和数字的数组，找到最长的子数组，且包含的字母和数字的个数相同。
# 返回该子数组，若存在多个最长子数组，返回左端点下标值最小的子数组。若不存在这样的数组，返回一个空数组。


# https://leetcode.cn/problems/find-longest-subarray-lcci/solution/tao-lu-qian-zhui-he-ha-xi-biao-xiao-chu-3mb11/
# 消除分支的技巧:
# !对于任意小写/大写英文字母字符，其 ASCII 码的二进制都形如 01xxxxxx；
# !对于任意数字字符，其 ASCII 码的二进制都形如 0011xxxx。
# (c >> 6 & 1) * 2 - 1 # 1(字母) or -1(数字)


class Solution:
    def findLongestSubarray(self, array: List[str]) -> List[str]:
        end, max_ = -1, -1
        preSum = defaultdict(int, {0: -1})
        curSum = 0
        for i, c in enumerate(array):
            curSum += 1 if c.isdigit() else -1
            if curSum in preSum:
                len_ = i - preSum[curSum]
                if len_ > max_:
                    end, max_ = i, len_
            else:
                preSum[curSum] = i
        return [] if max_ == -1 else array[end + 1 - max_ : end + 1]


assert Solution().findLongestSubarray(["A", "1"]) == ["A", "1"]
