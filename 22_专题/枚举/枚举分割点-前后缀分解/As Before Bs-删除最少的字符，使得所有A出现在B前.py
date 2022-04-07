# 删除最少的字符，使得所有A出现在B前
# 对每个位置，看左边的B和右边的A，这些都要删除

# 枚举前后缀
class Solution:
    def solve(self, s: str):
        rightA = s.count('A')
        leftB = 0

        res = rightA
        for char in s:
            if char == 'A':
                rightA -= 1
            else:
                leftB += 1
            res = min(res, leftB + rightA)

        return res


print(Solution().solve(s="AABBAB"))
# 1
# We can remove the last A to get AABBB.
