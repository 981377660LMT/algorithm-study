from typing import List
from collections import defaultdict

# 2 <= flowers.length <= 105
# -104 <= flowers[i] <= 104
# 一个花园满足下列条件时，该花园是有效的。

# 花园中至少包含两朵花。
# 第一朵花和最后一朵花的美观度`相同`。

# 返回你去除了任意朵花（也可以不去除任意一朵）之后形成的有效花园中最大可能的美观度。


# 1. 为了取出端点相同的这一段子数组，我们需要用adjMap保存同种花的index
# 2. 子数组最大和问题
# 3. 计算前缀和的时候忽略掉负数，但在最后求和的时候要把两端的负数加上（因为中间的负数花都可以移除，两端不行）
class Solution:
    def maximumBeauty(self, flowers: List[int]) -> int:
        indexMap = defaultdict(list)
        for i, flower in enumerate(flowers):
            indexMap[flower].append(i)

        preSum = [0]
        for flower in flowers:
            # 处理前缀和的时候忽略负数(剪花)
            preSum.append(preSum[-1] + max(0, flower))

        res = -0x7FFFFFFF
        for indexes in indexMap.values():
            if len(indexes) <= 1:
                continue
            first, last = indexes[0], indexes[-1]
            curSum = flowers[first] + flowers[last] + preSum[last] - preSum[first + 1]
            res = max(res, curSum)

        return res


print(Solution().maximumBeauty(flowers=[1, 2, 3, 1, 2]))
# 输出: 8
# 解释: 你可以修整为有效花园 [2,3,1,2] 来达到总美观度 2 + 3 + 1 + 2 = 8。
