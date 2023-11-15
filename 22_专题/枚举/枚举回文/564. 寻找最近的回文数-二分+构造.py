# 564. 寻找最近的回文数
# https://leetcode-cn.com/problems/find-the-closest-palindrome/comments/467740
# 小于等于10的数，返回n-1
# 10的幂，返回n-1
# 11，这个数字比较特殊，返回9

from enumeratePalindrome import getPalindromeByHalf


INF = int(1e20)


class Solution:
    def nearestPalindromic(self, n: str) -> str:
        def bisect(length: int) -> None:
            nonlocal res, minDiff

            left = 10 ** ((length - 1) >> 1)
            right = left * 10 - 1

            while left <= right:
                mid = (left + right) >> 1
                palindrome = int(getPalindromeByHalf(mid, length % 2 == 0))
                diff = abs(palindrome - target)
                if palindrome != target:
                    if diff < minDiff:
                        minDiff = abs(palindrome - target)
                        res = palindrome
                    elif diff == minDiff and palindrome < res:
                        res = palindrome

                if palindrome >= target:
                    right = mid - 1
                else:
                    left = mid + 1

        if int(n) <= 9:
            return str(int(n) - 1)

        target = int(n)
        minDiff, res = INF, INF
        bisect(len(n))
        bisect(len(n) + 1)
        bisect(len(n) - 1)

        return str(res)


print(Solution().nearestPalindromic("12422"))
print(Solution().nearestPalindromic("123"))
print(Solution().nearestPalindromic("999"))
print(Solution().nearestPalindromic("9"))
