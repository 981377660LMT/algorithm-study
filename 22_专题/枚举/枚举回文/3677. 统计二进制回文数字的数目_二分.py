# 3677. 统计二进制回文数字的数目
# 给你一个 非负 整数 n。
# 如果一个 非负 整数的二进制表示（不含前导零）正着读和倒着读都一样，则称该数为 二进制回文数。
# 返回满足 0 <= k <= n 且 k 的二进制表示是回文数的整数 k 的数量。
# 注意： 数字 0 被认为是二进制回文数，其表示为 "0"。
#
# !确定回文中心后，二分回文前缀


class Solution:
    def countBinaryPalindromes(self, n: int) -> int:
        def solve(center: str) -> int:
            left, right = 0, int(1e8)
            while left <= right:
                mid = (left + right) // 2
                half = bin(mid)[2:]
                palindrome = half + center + half[::-1]
                if int(palindrome, 2) <= n:
                    left = mid + 1
                else:
                    right = mid - 1
            return left

        if n <= 1:
            return n + 1
        res = solve("") + solve("0") + solve("1")
        return res - 1
