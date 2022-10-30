# 给你两个正整数 n 和 target 。
# 如果某个整数每一位上的数字相加小于或等于 target ，
# 则认为这个整数是一个 美丽整数 。
# !找出并返回满足 n + x 是 美丽整数 的最小非负整数 x 。
# 生成的输入保证总可以使 n 变成一个美丽整数。


class Solution:
    def makeIntegerBeautiful(self, n: int, target: int) -> int:
        def carry(num: int) -> int:
            """最低的非0位进位(carry)

            16 -> 20
            19 -> 20
            20 -> 100
            101 -> 110
            110 -> 200
            """
            base = 10
            while True:
                div, mod = divmod(num, base)
                if mod == 0:
                    base *= 10
                else:
                    return (div + 1) * base

        cur = n
        while True:
            if sum(map(int, str(cur))) <= target:
                return cur - n
            cur = carry(cur)


print(Solution().makeIntegerBeautiful(n=16, target=6))
print(Solution().makeIntegerBeautiful(n=467, target=6))
print(Solution().makeIntegerBeautiful(n=1, target=1))
print(Solution().makeIntegerBeautiful(n=3, target=2))
