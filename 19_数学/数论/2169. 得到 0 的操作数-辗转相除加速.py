# 辗转相除加速
# 每一步 操作 中，如果 num1 >= num2 ，你必须用 num1 减 num2 ；否则，你必须用 num2 减 num1


class Solution:
    def countOperations(self, num1: int, num2: int) -> int:
        res = 0
        # O(log) 辗转相除法加速，直到出现0
        while num1 and num2:
            res += num2 // num1
            num1, num2 = num2 % num1, num1
        return res


print(Solution().countOperations(1, 2))
print(Solution().countOperations(10, 10))
print(Solution().countOperations(2, 3))


# D - Count Subtractions
# https://atcoder.jp/contests/abc297/tasks/abc297_d
# 如果 num1 >= num2 ，你必须用 num1 减 num2 ；否则，你必须用 num2 减 num1
# 求 nums1 == nums2 的最小操作数
def countSubtractions(a: int, b: int) -> int:
    res = 0
    while a and b:
        if a < b:
            a, b = b, a
        div, mod = a // b, a % b
        res += div
        a, b = b, mod
    return res - 1


if __name__ == "__main__":
    a, b = map(int, input().split())
    print(countSubtractions(a, b))
