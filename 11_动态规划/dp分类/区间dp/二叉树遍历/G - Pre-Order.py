"""
给定一个dfs序,问能够构建出来的子树有多少种
每一次dfs都是优先选择节点的编号最小的节点。
n<=500

区间dp
类似于479. 加分二叉树-中序遍历分割左右子树
"""
import sys
import os

sys.setrecursionlimit(int(1e6))
input = lambda: sys.stdin.readline().rstrip("\r\n")
MOD = 998244353


def main() -> None:
    def dfs(left: int, right: int) -> int:
        """区间[left,right]形成的子树有多少种"""
        if left >= right:
            return 1
        if memo[left][right] != -1:
            return memo[left][right]
        res = dfs(left + 1, right)  # left肯定是合法的根节点(按照dfs序)
        for i in range(left + 1, right + 1):  # 选出来一些点，使得其和 left 并列成为子节点
            if preOrder[left] < preOrder[i]:
                res += dfs(left + 1, i - 1) * dfs(i, right)
                res %= MOD
        memo[left][right] = res
        return res % MOD

    n = int(input())
    preOrder = [int(num) - 1 for num in input().split()]
    memo = [[-1] * (n + 10) for _ in range((n + 10))]
    res = dfs(1, n - 1)  # 0为唯一的根节点,不考虑
    print(res)


if __name__ == "__main__":
    if os.environ.get("USERNAME", " ") == "caomeinaixi":
        while True:
            main()
    else:
        main()
