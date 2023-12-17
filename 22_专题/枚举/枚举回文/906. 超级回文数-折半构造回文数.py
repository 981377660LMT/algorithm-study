# 如果一个正整数自身是回文数，而且它也是一个回文数的平方，那么我们称这个数为超级回文数。
# 返回包含在范围 [L, R] 中的超级回文数的数目。
# L 和 R 是表示 [1, 10^18) 范围的整数的字符串。

from enumeratePalindrome import emumeratePalindromeByLength


palindromes = [int(v) for v in emumeratePalindromeByLength(1, 9)]


class Solution:
    def superpalindromesInRange(self, left: str, right: str) -> int:
        res = []
        lower, upper = int(left), int(right)
        for num in palindromes:
            square = num**2
            if square >= lower and square <= upper and square == int(str(square)[::-1]):
                res.append(square)
        return len(res)


if __name__ == "__main__":
    print(Solution().superpalindromesInRange(left="4", right="1000"))
    # 输出：4
    # 解释：
    # 4，9，121，以及 484 是超级回文数。
    # 注意 676 不是一个超级回文数： 26 * 26 = 676，但是 26 不是回文数。
