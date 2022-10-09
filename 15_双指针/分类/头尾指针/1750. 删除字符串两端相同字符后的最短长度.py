"""删除字符串两端相同字符后的最短长度

left<right 等价于 len(queue)>=2
left<=right 等价于 len(queue)>=1
"""


from collections import deque


class Solution:
    def minimumLength(self, s: str) -> int:
        left, right = 0, len(s) - 1
        while left < right and s[left] == s[right]:  # 注意'cbc'的情况
            cur = s[left]
            while left <= right and s[left] == cur:
                left += 1
            while left <= right and s[right] == cur:
                right -= 1
        return right - left + 1

    def minimumLength2(self, s: str) -> int:
        """deque不容易错"""
        queue = deque(s)
        while len(queue) >= 2 and queue[0] == queue[-1]:
            same = queue[0]
            while queue and queue[0] == same:
                queue.popleft()
            while queue and queue[-1] == same:
                queue.pop()
        return len(queue)


print(Solution().minimumLength(s="aabccabba"))
print(Solution().minimumLength(s="cabaabac"))
print(Solution().minimumLength(s="bbbbbbbbbbbbbbbbbbbbbbbbbbbabbbbbbbbbbbbbbbccbcbcbccbbabbb"))
print(Solution().minimumLength(s="cbc"))
print(Solution().minimumLength2(s="cbc"))

# 输出：3
# 解释：最优操作序列为：
# - 选择前缀 "aa" 和后缀 "a" 并删除它们，得到 s = "bccabb" 。
# - 选择前缀 "b" 和后缀 "bb" 并删除它们，得到 s = "cca" 。
