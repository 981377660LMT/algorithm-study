from typing import List, Optional, Tuple

MOD = int(1e9 + 7)
INF = int(1e20)

# palindrome = [1, 2, 3, 4, 5, 6, 7, 8, 9]
# # 构造九位数的回文只需要左右部分最多四位 一共有11*10000+9个回文数
# for side in range(1, 10000):
#     s1 = str(side) + str(side)[::-1]
#     palindrome.append(int(s1))
#     for mid in range(10):
#         s2 = str(side) + str(mid) + str(side)[::-1]
#         palindrome.append(int(s2))
# # palindrome.sort()

# print(palindrome)


class Solution:
    def kthPalindrome(self, queries: List[int], intLength: int) -> List[int]:
        res = [-1] * len(queries)
        if intLength == 1:
            palindrome = [1, 2, 3, 4, 5, 6, 7, 8, 9]
            for i, q in enumerate(queries):
                if q - 1 >= len(palindrome):
                    res[i] = -1
                else:
                    res[i] = palindrome[q - 1]
            return res

        # palindrome = []
        # start = pow(10, ((intLength) >> 1) - 1)

        # if intLength & 1:
        #     for side in range(start, start * 10):
        #         for mid in range(10):
        #             s2 = str(side) + str(mid) + str(side)[::-1]
        #             palindrome.append(int(s2))
        # else:
        #     for side in range(start, start * 10):
        #         s1 = str(side) + str(side)[::-1]
        #         palindrome.append(int(s1))

        start = pow(10, ((intLength - 1) >> 1))
        end = start * 10
        count = end - start
        for i, q in enumerate(queries):
            if q - 1 >= count:
                continue
            if intLength & 1:
                half = start + q - 1
                res[i] = int(str(half)[:-1] + str(half)[::-1])
            else:
                half = start + q - 1
                res[i] = int(str(half) + str(half)[::-1])

        return res


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
