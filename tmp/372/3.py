from typing import List, Tuple, Optional
from collections import defaultdict, Counter
from sortedcontainers import SortedList


MOD = int(1e9 + 7)
INF = int(1e20)

# 给你三个整数 a ，b 和 n ，请你返回 (a XOR x) * (b XOR x) 的 最大值 且 x 需要满足 0 <= x < 2n。

# 由于答案可能会很大，返回它对 109 + 7 取余 后的结果。


# 注意，XOR 是按位异或操作。

# class Solution(object):
#     def maximumXorProduct(self, a, b, n):
#         max_val = a * b # 初始化最大值为a * b（不异或）
#         for i in range(n): # 遍历每一位
#             cur = (a ^ (1 << i)) * (b ^ (1 << i)) # 当前位异或后的乘积
#             if max_val < cur: # 比之前的最大值大
#                 a ^= 1 << i # 异或
#                 b ^= 1 << i # 异或
#                 max_val = cur # 更新最大值

#         return max_val % (10 ** 9 + 7) # 返回取余后的最大值


class Solution:
    def maximumXorProduct(self, a: int, b: int, n: int) -> int:
        def solve(counter1: List[int], counter2: List[int], f: int) -> int:
            f = 0
            pos = [[], []]
            for i in range(n - 1, -1, -1):
                if counter1[i] == 1 and counter2[i] == 0:
                    pos[f].append(i)
                    counter1[i] = 0
                    f ^= 1
                elif counter1[i] == 0 and counter2[i] == 1:
                    pos[f].append(i)
                    counter2[i] = 0
                    f ^= 1
            num1 = sum((1 << i) * counter1[i] for i in range(64))
            num2 = sum((1 << i) * counter2[i] for i in range(64))
            for i in pos[0]:
                num1 += 1 << i
            for i in pos[1]:
                num2 += 1 << i
            return num1 * num2

        def solve2(counter1: List[int], counter2: List[int]) -> int:
            pos = []
            for i in range(n - 1, -1, -1):
                if counter1[i] == 1 and counter2[i] == 0:
                    pos.append(i)
                    counter1[i] = 0
                elif counter1[i] == 0 and counter2[i] == 1:
                    pos.append(i)
                    counter2[i] = 0
            num1 = sum((1 << i) * counter1[i] for i in range(64))
            num2 = sum((1 << i) * counter2[i] for i in range(64))
            for i in pos:
                num1 += 1 << i
            return num1 * num2

        def solve3(counter1: List[int], counter2: List[int]) -> int:
            pos = []
            for i in range(n - 1, -1, -1):
                if counter1[i] == 1 and counter2[i] == 0:
                    pos.append(i)
                    counter1[i] = 0
                elif counter1[i] == 0 and counter2[i] == 1:
                    pos.append(i)
                    counter2[i] = 0
            num1 = sum((1 << i) * counter1[i] for i in range(64))
            num2 = sum((1 << i) * counter2[i] for i in range(64))
            for i in pos[:1]:
                num1 += 1 << i
            for i in pos[1:]:
                num2 += 1 << i
            return num1 * num2

        bitCounter1 = [0] * 64
        bitCounter2 = [0] * 64
        for i in range(64):
            bitCounter1[i] = (a >> i) & 1
            bitCounter2[i] = (b >> i) & 1

        for i in range(n):
            if bitCounter1[i] == 0 and bitCounter2[i] == 0:
                bitCounter1[i] = 1
                bitCounter2[i] = 1

        res1 = solve(bitCounter1.copy(), bitCounter2.copy(), 0)
        res2 = solve(bitCounter1.copy(), bitCounter2.copy(), 1)
        res3 = solve2(bitCounter1.copy(), bitCounter2.copy())
        res4 = solve2(bitCounter2.copy(), bitCounter1.copy())
        res5 = solve3(bitCounter1.copy(), bitCounter2.copy())
        res6 = solve3(bitCounter2.copy(), bitCounter1.copy())
        return max(res1, res2, res3, res4, res5, res6, a * b) % MOD


def bruteForce(a: int, b: int, n: int) -> int:
    # 每位选还是不选
    res = 0
    for state in range(1 << n):
        curNum = sum((1 << i) * (state >> i & 1) for i in range(n))
        res = max(res, (a ^ curNum) * (b ^ curNum))
    return res % MOD


# a = 12, b = 5, n = 4
# print(Solution().maximumXorProduct(5, 2, 1))
# print(bruteForce(5, 2, 1))
# print(Solution().maximumXorProduct(7, 8, 2))
# print(bruteForce(7, 8, 2))
# print(Solution().maximumXorProduct(1, 6, 1))
print(Solution().maximumXorProduct(2, 5, 3))
print(bruteForce(2, 5, 3))
print(bruteForce(2, 5, 3))
print(bruteForce(2, 5, 3))
if __name__ == "__main__":
    import random

    for _ in range(100):
        a = random.randint(1, 100)
        b = random.randint(1, 100)
        n = random.randint(1, 10)
        if Solution().maximumXorProduct(a, b, n) != bruteForce(a, b, n):
            print(a, b, n, "bad")
            break
