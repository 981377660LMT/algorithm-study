class Solution:
    def isStrobogrammatic(self, num: str) -> bool:
        map = dict(zip('01689', '01986'))
        l, r = 0, len(num) - 1
        while l <= r:
            if num[l] not in map or num[r] not in map or map[num[l]] != num[r]:
                return False
            l += 1
            r -= 1
        return True


print(Solution().isStrobogrammatic("609"))
