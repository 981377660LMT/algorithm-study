# 866. 回文素数-折半构造回文数
# https://leetcode.cn/problems/prime-palindrome/


from enumeratePalindrome import emumeratePalindrome


def isPrime(n: int) -> bool:
    return n >= 2 and all(n % i for i in range(2, int(n**0.5) + 1))


class Solution:
    def primePalindrome(self, n: int) -> int:
        """
        求出大于或等于 N 的最小回文素数。
        1 <= N <= 10^8
        """

        for p in emumeratePalindrome(1, 9):
            p = int(p)
            if p < n:
                continue
            if isPrime(p):
                return p

        return -1


if __name__ == "__main__":
    for cand in emumeratePalindrome(7, 8):  # 生成回文素数
        if isPrime(int(cand)):
            print(cand)
# 10301
# 10501
# 10601
# 11311
# 11411
# 12421
# 12721
# 12821
# 13331
