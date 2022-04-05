# Return the number of pairs (i, j) such that lower ≤ a[i] * a[i] + b[j] * b[j] ≤ upper
# atMoskK

# 平方和在范围里的对数 排序首尾指针
class Solution:
    def solve(self, a, b, lower, upper):
        def atMoskK(threshold):
            left = 0
            right = len(b) - 1
            res = 0
            while left < len(a) and right >= 0:
                if a[left] ** 2 + b[right] ** 2 <= threshold:
                    res += right + 1  # 0 到 right 都可以
                    left += 1
                else:
                    right -= 1
            return res

        a = sorted([abs(x) for x in a])
        b = sorted([abs(x) for x in b])
        return atMoskK(upper) - atMoskK(lower - 1)


print(Solution().solve(a=[3, -1, 9], b=[100, 5, -2], lower=7, upper=99))
