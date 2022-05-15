from collections import deque, Counter, defaultdict

# 请你通过若干次以下操作将字符串 s 转化成字符串 t ：
# 选择 s 中一个 非空 子字符串并将它包含的字符就地 `升序 排序`。

# 如果可以将字符串 s 变成 t ，返回 true 。否则，返回 false 。
# 1 <= s.length <= 1e5


class Solution:
    def isTransformable(self, s: str, t: str) -> bool:
        """邻接表，检查原来的每一种逆序对是否增加        
        操作前后，任何一种逆序对的不同的二元组(i,j)的数目都只能减少而不能增加
        时间复杂度O(n*10)
        """
        if Counter(s) != Counter(t):
            return False
        nums1 = list([int(char) for char in s])
        nums2 = list([int(char) for char in t])

        indexMap = defaultdict(deque)
        for i, num in enumerate(nums1):
            indexMap[num].append(i)

        for num in nums2:
            first = indexMap[num].popleft()
            for smaller in range(num):
                # 例如 3...8 变成了 8...3 逆序对增加了
                if indexMap[smaller] and indexMap[smaller][0] < first:
                    return False

        return True


print(Solution().isTransformable(s="84532", t="34852"))
# 输出：true
# 解释：你可以按以下操作将 s 转变为 t ：
# "84532" （从下标 2 到下标 3）-> "84352"
# "84352" （从下标 0 到下标 2） -> "34852"
