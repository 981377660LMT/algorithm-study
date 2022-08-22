# 三元组一般是`定一移二`的做法
# 或者枚举贡献
# 网易笔试t4
# !求三元组i<j<k的个数 nums[i]==nums[k]>nums[j]
# n<=1e5


# 离线查询区间内比key小的数的个数 -> 按照key排序
# !将所有查询按照key从小到大排序，把原来的序号挨个放入SortedList/树状数组
# !然后对每个区间,二分索引/查找索引间数字的个数，就是区间内比key小的数字的个数
# ps:不过这种方法似乎并不能推广到任意的离线查询 [1,5] [3,8] 这种区间有重叠就不行


from typing import List
from collections import defaultdict


class BIT:
    """单点修改的树状数组"""

    def __init__(self, n: int):
        self.size = n
        self.tree = defaultdict(int)

    def add(self, index: int, delta: int) -> None:
        if index <= 0:
            raise ValueError("index 必须是正整数")
        while index <= self.size:
            self.tree[index] += delta
            index += index & -index

    def query(self, index: int) -> int:
        if index > self.size:
            index = self.size
        res = 0
        while index > 0:
            res += self.tree[index]
            index -= index & -index
        return res

    def queryRange(self, left: int, right: int) -> int:
        return self.query(right) - self.query(left - 1)


class Solution:
    def goodTriplets(self, nums: List[int]) -> int:
        """求三元组i<j<k的个数 nums[i]==nums[k]>nums[j]"""
        indexMap = defaultdict(list)
        for i, num in enumerate(nums):
            indexMap[num].append(i + 1)

        res = 0
        bit = BIT(int(1e9 + 10))
        for key in sorted(indexMap):
            indexes = indexMap[key]
            bit.add(indexes[0], 1)
            for i, (pre, cur) in enumerate(zip(indexes, indexes[1:])):
                left, right = i + 1, len(indexes) - (i + 1)
                # !求区间[pre+1:cur]中小于key的数的个数
                # !因为已经排序了 只需查询索引区间内之前有多少个数
                count = bit.queryRange(pre + 1, cur - 1)
                res += count * left * right
                bit.add(cur, 1)
        return res


print(Solution().goodTriplets(nums=[3, 1, 3, 4, 3, 4]))  # 3
print(Solution().goodTriplets(nums=[3, 4, 2, 5, 6, 1, 2, 3]))  # 4
