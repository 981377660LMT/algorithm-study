# 650. 两个键的键盘
# https://leetcode.cn/problems/2-keys-keyboard/description/


class Solution:
    def minSteps(self, n: int) -> int:
        """
        质因数分解法：
        观察：要从 1 个 A 变成 n 个 A，最优操作数等于 n 的质因数之和。
        因为每个质因子 p 意味着一次 “复制全部” + (p-1) 次 “粘贴” 共 p 步，
        将当前长度扩大 p 倍。

        算法：对 n 从小到大试除所有可能因子 d，
        每次能整除就把 d 加进结果，并把 n //= d，直到 n 归 1。
        """
        res = 0
        d = 2
        while d * d <= n:
            while n % d == 0:
                res += d
                n //= d
            d += 1
        if n > 1:
            res += n
        return res


if __name__ == "__main__":
    sol = Solution()
    for n in [1, 2, 3, 4, 5, 6, 9, 18]:
        print(n, sol.minSteps(n))

    # 预期输出：
    # 1 0   （已是 1，不需任何操作）
    # 2 2   （Copy + Paste）
    # 3 3   （Copy + Paste + Paste）
    # 4 4   （Copy + Paste -> 2A, Copy + Paste -> 4A）
    # 5 5
    # 6 5   （质因数分解 6=2×3，2+3=5）
    # 9 6   （9=3×3，3+3=6）
    # 18 8  （18=2×3×3，2+3+3=8）
