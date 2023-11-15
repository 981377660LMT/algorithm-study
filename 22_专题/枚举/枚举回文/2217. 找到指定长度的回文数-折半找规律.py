# 反思：这道题打表去了，实际上打表是过不了的
# 2217. 找到指定长度的回文数
#
# https://leetcode.cn/problems/find-palindrome-with-fixed-length/
# intLength<=15 有1e8种情况
# 1 <= queries.length <= 5 * 104
# 1 <= queries[i] <= 109
# 1 <= intLength <= 15


from typing import List


class Solution:
    def kthPalindrome(self, queries: List[int], intLength: int) -> List[int]:
        res = [-1] * len(queries)

        # 长为3，4的回文都是从10开始的，所以只需要构造10-99的回文即可
        start = pow(10, ((intLength - 1) >> 1))
        count = (start * 10 - 1) - start + 1

        for i, q in enumerate(queries):
            if q - 1 >= count:
                continue
            half = start + q - 1
            if intLength & 1:
                res[i] = int(str(half)[:-1] + str(half)[::-1])
            else:
                res[i] = int(str(half) + str(half)[::-1])
        return res


print(
    Solution().kthPalindrome(
        queries=[1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 9, 8], intLength=1
    )
)
print(Solution().kthPalindrome(queries=[1, 2, 3, 4, 5, 90], intLength=3))
print(Solution().kthPalindrome(queries=[2, 4, 6], intLength=4))
print(Solution().kthPalindrome(queries=[2, 4, 6], intLength=2))
print(
    Solution().kthPalindrome(
        queries=[2, 201429812, 8, 520498110, 492711727, 339882032, 462074369, 9, 7, 6], intLength=1
    )
)
print(
    Solution().kthPalindrome(
        queries=[
            475098318,
            62,
            457771600,
            85,
            476799241,
            23,
            73,
            600686743,
            58,
            264628531,
            26,
            25,
            9,
        ],
        intLength=13,
    )
)
