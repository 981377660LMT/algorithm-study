# 你有两个篮子，每个篮子可以携带任何数量的水果，但你希望每个篮子只携带一种类型的水果。
# !求只包含两种元素的最长子数组
from collections import defaultdict
from typing import List


class Solution:
    def totalFruit(self, fruits: List[int]) -> int:
        res, left, n = 0, 0, len(fruits)
        counter = defaultdict(int)
        for right in range(n):
            counter[fruits[right]] += 1
            while left <= right and len(counter) > 2:
                counter[fruits[left]] -= 1
                if counter[fruits[left]] == 0:
                    del counter[fruits[left]]
                left += 1
            res = max(res, right - left + 1)
        return res


print(Solution().totalFruit([1, 2, 1]))
