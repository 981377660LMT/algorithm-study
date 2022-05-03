from collections import defaultdict
from typing import List, Optional, Tuple

MOD = int(1e9 + 7)
INF = int(1e20)

# 1 <= nums.length <= 200


class Solution:
    def countDistinct(self, nums: List[int], k: int, p: int) -> int:
        """按照起点枚举子数组
        
        tuple + set 去重
        O(n^3)
        """
        res = set()
        for start in range(len(nums)):
            count = 0
            for end in range(start, len(nums)):
                if nums[end] % p == 0:
                    count += 1
                if count <= k:
                    res.add(tuple((nums[start : end + 1])))
                else:
                    break

        return len(res)

    def countDistinct2(self, nums: List[int], k: int, p: int) -> int:
        """按照起点枚举子数组
        
        对于每个满足题意的子数组，我们将它加入字典树。由于每产生一个不同的子数组，
        必然将会在字典树中插入一个节点，故最终答案 = 字典树中新插入的节点数量。
        O(n^2)
        """
        trie = dict()
        res = 0
        for start in range(len(nums)):
            root, count = trie, 0
            for end in range(start, len(nums)):
                if nums[end] % p == 0:
                    count += 1

                if count <= k:
                    if nums[end] not in root:
                        root[nums[end]] = {}
                        res += 1
                    root = root[nums[end]]
                else:
                    break

        return res


print(Solution().countDistinct(nums=[2, 3, 3, 2, 2], k=2, p=2))

# [2]、[2,3]、[2,3,3]、[2,3,3,2]、[3]、[3,3]、[3,3,2]、[3,3,2,2]、[3,2]、[3,2,2] 和 [2,2] 。
