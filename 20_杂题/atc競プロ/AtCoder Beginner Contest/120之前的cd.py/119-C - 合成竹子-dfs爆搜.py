# 3<=n<=8
# 变化竹子 O(4^n)
# !有n根竹子,要获得长度为A,B,C的三根竹子(不必使用所有竹子)

# !有三种操作:
# 1. 花费1,将一根竹子长度加1
# 2. 花费1,将一根竹子长度减1
# 3. 花费10,将两根竹子拼接

# 问合成三根长度为A,B,C的竹子最少需要多少费用
# !每根竹子四种转移,所以是4^n种状态

import sys

sys.setrecursionlimit(int(1e9))
input = lambda: sys.stdin.readline().rstrip("\r\n")
MOD = 998244353
INF = int(4e18)

if __name__ == "__main__":

    def dfs(index: int, curA: int, curB: int, curC: int) -> int:
        if index == n:
            if curA > 0 and curB > 0 and curC > 0:
                return abs(curA - a) + abs(curB - b) + abs(curC - c)
            return INF

        # !当前的竹子最后组成ABC中的哪一个
        res = dfs(index + 1, curA, curB, curC)  # 不用当前竹子
        cand1 = dfs(index + 1, curA + nums[index], curB, curC) + (10 if curA > 0 else 0)
        cand2 = dfs(index + 1, curA, curB + nums[index], curC) + (10 if curB > 0 else 0)
        cand3 = dfs(index + 1, curA, curB, curC + nums[index]) + (10 if curC > 0 else 0)
        return min(res, cand1, cand2, cand3)

    n, a, b, c = map(int, input().split())
    nums = [int(input()) for _ in range(n)]
    print(dfs(0, 0, 0, 0))
