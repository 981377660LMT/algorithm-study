from typing import List

# 'D' 表示两个数字间的递减关系，'I' 表示两个数字间的递增关系
# 现在你的任务是找到具有最小字典序的 [1, 2, ... n] 的排列，使其能代表输入的 秘密签名。

# 按 'I' 分段，每一段尽量字典序最小即可。
class Solution:
    def findPermutation(self, s: str) -> List[int]:
        res = []
        curMin = 1
        for sub in s.split('I'):
            count = len(sub)
            res.extend(range(curMin + count, curMin - 1, -1))
            curMin += count + 1
        return res


print(Solution().findPermutation("DI"))
