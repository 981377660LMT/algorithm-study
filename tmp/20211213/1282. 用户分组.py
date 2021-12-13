from typing import List
from collections import defaultdict

# 给你一个长度为 n 的数组 groupSizes，其中包含每位用户所处的用户组的大小，请你返回用户分组情况
class Solution:
    def groupThePeople(self, groupSizes: List[int]) -> List[List[int]]:
        res = []
        dic = defaultdict(list)
        for id, size in enumerate(groupSizes):
            dic[size].append(id)
            if len(dic[size]) == size:
                res.append(dic.pop(size))
        return res


print(Solution().groupThePeople(groupSizes=[2, 1, 3, 3, 3, 2]))
# 输出：[[1],[0,5],[2,3,4]]
