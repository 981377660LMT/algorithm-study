# 如果一个人的右侧k个`都比他矮`，那么它可以看到电影
# 问那些人可以看电影
class Solution:
    def solve(self, heights, k):
        stack = []
        for i, h in enumerate(heights):
            while stack and i - stack[-1] <= k and h >= heights[stack[-1]]:
                stack.pop()
            stack.append(i)
        return stack


print(Solution().solve(heights=[9, 8, 7, 7, 4, 9], k=2))
