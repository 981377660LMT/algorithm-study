# f(x) 是 x! 末尾是 0 的数量。
# !如果f(x)=k有解，则解的个数只可能是5

# 给定 k，找出返回能满足 f(x) = k 的非负整数 x 的数量。
# !原问题转化为：在非负整数中，有多少个数进行阶乘分解后，所含质因数 5 的个数恰好为 k 个。
# !容斥原理
class Solution:
    def preimageSizeFZF(self, k: int) -> int:
        """给定 k,找出返回能满足 f(x) = k 的非负整数 x 的数量。(答案只可能是0或5)"""

        def countTrailingZeros(n) -> int:
            """求n!的尾随0的个数,时间复杂度O(log(n))"""
            res = 0
            div = 5
            while n // div:
                res += n // div
                div *= 5
            return res

        def cal(count5: int) -> int:
            """a!所含质因数 5 的个数恰好为 count5 个时a的最大值"""
            left, right = 0, int(1e10)
            while left <= right:
                mid = (left + right) // 2
                if countTrailingZeros(mid) > count5:
                    right = mid - 1
                else:
                    left = mid + 1
            return right

        return cal(k) - cal(k - 1)


print(Solution().preimageSizeFZF(5))  # 0
print(Solution().preimageSizeFZF(3))  # 5
