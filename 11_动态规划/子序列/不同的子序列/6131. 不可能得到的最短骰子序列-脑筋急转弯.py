# 不可能得到的最短骰子序列
from typing import List


# 给你一个长度为 n 的整数数组 rolls 和一个整数 k 。
# 你扔一个 k 面的骰子 n 次，骰子的每个面分别是 1 到 k ，其中第 i 次扔得到的数字是 rolls[i] 。
# 请你返回 无法 从 rolls 中得到的 最短 骰子子序列的长度。

# !满足长度为 1,2,3,4... 的序列全集需要满足什么条件
class Solution:
    def shortestSequence(self, rolls: List[int], k: int) -> int:
        """每一段每一段地考察 因为要让所有子序列都在原序列出现 所以各个开头(大家)要齐头并进"""
        visited = set()
        res = 0
        for char in rolls:
            visited.add(char)
            if len(visited) == k:  # 多凑齐了一个长度
                res += 1
                # visited.clear()
                visited = set()
        return res + 1


print(Solution().shortestSequence(rolls=[4, 2, 1, 2, 3, 3, 2, 4, 1], k=4))
print(Solution().shortestSequence(rolls=[1, 1, 2, 2], k=2))
print(Solution().shortestSequence(rolls=[1, 1, 3, 2, 2, 2, 3, 3], k=4))

# !如果这个题要求构造一个得不到的最短的子序列 那么做法也类似，最后查一下对应ans+1这组的哈希集合缺什么就可以了
def shortestSequence2(rolls: List[int], k: int) -> List[int]:
    visited = set()
    res = []
    for char in rolls:
        visited.add(char)
        if len(visited) == k:  # 多凑齐了一个长度
            res.append(char)
            visited.clear()

    for i in range(1, k + 1):
        if i not in visited:
            res.append(i)
            return res
    return res
