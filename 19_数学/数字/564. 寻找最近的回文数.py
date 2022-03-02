# 小于等于10的数，返回n-1
# 10的幂，返回n-1
# 11，这个数字比较特殊，返回9
# https://leetcode-cn.com/problems/find-the-closest-palindrome/comments/467740


class Solution:
    def nearestPalindromic(self, n: str) -> str:
        if int(n) < 10 or int(n[::-1]) == 1:
            return str(int(n) - 1)

        if n == '11':
            return '9'

        half = (len(n) + 1) >> 1
        a, b = n[:half], n[half:]
        temp = [str(int(a) - 1), a, str(int(a) + 1)]
        temp = [i + i[len(b) - 1 :: -1] for i in temp]
        print(temp)
        return min(temp, key=lambda x: abs(int(x) - int(n)))


print(Solution().nearestPalindromic('12422'))
