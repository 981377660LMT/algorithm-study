# n ≤ 10^5
# 牛牛有n件带水的衣服，干燥衣服有两种方式。
# 一、是用烘干机，可以每分钟烤干衣服的k滴水。
# 二、是自然烘干，每分钟衣服会自然烘干1滴水。
# 烘干机比较小，每次只能放进一件衣服。
# 注意，使用烘干机的时候，其他衣服仍然可以保持自然烘干状态，现在牛牛想知道最少要多少时间可以把衣服全烘干。

from math import ceil
from typing import List


class Solution:
    def dryClothes(self, n: int, nums: List[int], k: int):
        def check(mid: int) -> bool:
            time = 0
            for num in nums:
                if num <= mid:
                    continue
                # 列不等式：机器烘干的水滴数为t*k，自然烘干时间为mid-t，自然烘干的水滴数为m-t。
                # 烘干水滴数要>=num
                time += ceil((num - mid) / (k - 1))
            return time <= mid

        # write code here
        left, right = 1, max(nums)
        while left <= right:
            mid = (left + right) >> 1
            if check(mid):
                right = mid - 1
            else:
                left = mid + 1
        return left


print(1)
