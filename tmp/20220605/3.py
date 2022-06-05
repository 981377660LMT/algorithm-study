from collections import defaultdict
from typing import List, Optional, Tuple

MOD = int(1e9 + 7)
INF = int(1e20)

# !注意到每次操作前后数组里的元素都是互不相同的，
# 因此用一个 hash map 维护每种元素在哪个位置。
# 这样每次操作即可 O(1) 查找需要替换的位置
# 其实有序容器找位置也可以logn
# !之前一直想着index找位置是O(n) 没想到用哈希表存元素的位置

# !一直想着怎么对顺序进行变形
# !没想到直接找元素然后改


class Solution:
    def arrayChange(self, nums: List[int], operations: List[List[int]]) -> List[int]:
        pre, next = dict(), dict()
        for a, b in operations:
            pre[b] = pre.get(a, a)
            next[pre[b]] = b
        return [next.get(v, v) for v in nums]


# [91,93,94,95,96,97,98,99,100,101,1,2,3,4,5,6,7,8,9,10,11,12,13,14,15,16,17,18,19,20,21,22,23,24,25,26,27,28,29,30,31,32,33,34,35,36,37,38,39,40,41,42,43,44,45,46,47,48,49,50,51,52,53,54,55,56,57,58,59,60,61,62,63,64,65,66,67,68,69,70,71,72,73,74,75,76,77,78,79,80,81,82,83,84,85,86,87,88,89,90]
