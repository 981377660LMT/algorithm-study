# https://leetcode.cn/problems/intersection-of-two-arrays-ii/description/
# 进阶：
#
# - 如果给定的数组已经排好序呢？你将如何优化你的算法？-> 双指针
# - 如果 nums1 的大小比 nums2 小，哪种方法更优？-> 长度小的数组转成哈希表
# - 如果 nums2 的元素存储在磁盘上，内存是有限的，并且你不能一次加载所有的元素到内存中，你该怎么办？
#   用一个小型缓冲区（buffer）一边读数据一边遍历数据。这等价于问 nums2是一个流（Stream）的情况要怎么做。
#   由于我们写的是一次遍历的代码，所以已经符合这个要求。


from collections import Counter
from typing import List


class Solution:
    def intersect(self, nums1: List[int], nums2: List[int]) -> List[int]:
        if len(nums1) > len(nums2):
            nums1, nums2 = nums2, nums1
        res = []
        counter = Counter(nums1)
        for x in nums2:
            if counter[x] > 0:
                counter[x] -= 1
                res.append(x)
        return res
