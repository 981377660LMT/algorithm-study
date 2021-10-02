class Solution:
    def isStrobogrammatic(self, num: str) -> bool:
        return ''.join(dict(zip('01689', '01986')).get(x, '') for x in num) == num[::-1]


print(Solution().isStrobogrammatic("609"))
