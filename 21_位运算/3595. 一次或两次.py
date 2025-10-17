# 3595. 一次或两次
# https://leetcode.cn/problems/once-twice/description/?envType=problem-list-v2&envId=sZVESpvF
# 给定一个整数数组 nums。在这个数组中：
# 有一个元素出现了 恰好 1 次。
# 有一个元素出现了 恰好 2 次。
# 其它所有元素都出现了 恰好 3 次。
# 返回一个长度为 2 的整数数组，其中第一个元素是只出现 1 次 的那个元素，第二个元素是只出现 2 次 的那个元素。
# 你的解决方案必须在 O(n) 时间 与 O(1) 空间中运行。

from typing import List


class Solution:
    def onceTwice(self, nums: List[int]) -> List[int]:
        a, b, ma, mb, na, nb = 0, 0, 0, 0, 0, 0
        for c in nums:
            a, b = a ^ (c & (~b)), b ^ (c & (a | b))
        for c in nums:
            if c & a == a and c & b == 0:
                ma, mb = ma ^ (c & (~mb)), mb ^ (c & (ma | mb))
            if c & b == b and c & a == 0:
                na, nb = na ^ (c & (~nb)), nb ^ (c & (na | nb))
        return [ma, nb]
