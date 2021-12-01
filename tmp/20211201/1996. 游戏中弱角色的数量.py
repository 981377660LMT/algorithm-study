from typing import List

# 2 <= properties.length <= 105
# 如果存在一个其他角色的攻击和防御等级 都`严格高于` 该角色的攻击和防御等级，则认为该角色为 弱角色
# 一个降序，一个升序排列即可
class Solution:
    def numberOfWeakCharacters(self, properties: List[List[int]]) -> int:
        properties.sort(key=lambda x: (x[0], -x[1]))
        print(properties)
        res = 0
        pre = -0x7FFFFFFF

        # 这里从后往前遍历,因为要获取最大的防御值
        for _, cur in reversed(properties):
            if pre > cur:
                res += 1
            pre = max(pre, cur)
        return res


print(Solution().numberOfWeakCharacters([[1, 5], [10, 4], [4, 3]]))
# 输出：1
# 解释：第三个角色是弱角色，因为第二个角色的攻击和防御严格大于该角色。
