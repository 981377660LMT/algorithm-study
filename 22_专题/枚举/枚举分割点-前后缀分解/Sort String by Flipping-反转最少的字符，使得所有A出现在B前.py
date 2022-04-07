# 反转最少的字符，使得所有A出现在B前
# As Before Bs-删除最少的字符，使得所有A出现在B前  一样的
class Solution:
    def solve(self, s: str) -> int:
        rightA = s.count('x')
        leftB = 0

        res = rightA
        for char in s:
            if char == 'x':
                rightA -= 1
            else:
                leftB += 1
            res = min(res, leftB + rightA)

        return res


print(Solution().solve(s="xyxxxyxyy"))
# 2
# It suffices to flip the second and sixth characters.
