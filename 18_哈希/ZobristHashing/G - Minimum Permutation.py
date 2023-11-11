from collections import defaultdict
from random import randint
import sys

sys.setrecursionlimit(int(1e6))
input = lambda: sys.stdin.readline().rstrip("\r\n")
MOD = 998244353
INF = int(4e18)


if __name__ == "__main__":
    n, m = map(int, input().split())
    nums = list(map(int, input().split()))
    pool = defaultdict(lambda: randint(1, (1 << 61) - 1))
    targets = [0] * (m + 1)  # !targets[i] 1~i 的集合的哈希值
    for i in range(1, m + 1):
        targets[i] = targets[i - 1] ^ pool[i]
    all_ = targets[-1]

    # 二分最短长度,使得这之间出现了1-m所有数
    def check(mid: int) -> bool:
        """mid长度的滑窗中是否包含1-m所有数"""
        counter = defaultdict(int)
        curXor = 0
        for right in range(n):
            cur = nums[right]
            counter[cur] += 1
            if counter[cur] == 1:
                curXor ^= pool[cur]
            if right >= mid:
                num = nums[right - mid]
                counter[num] -= 1
                if counter[num] == 0:
                    curXor ^= pool[num]
            if right >= mid - 1 and curXor == all_:
                return True
        return False

    left, right = 0, n
    while left <= right:
        mid = (left + right) // 2
        if check(mid):
            right = mid - 1
        else:
            left = mid + 1

    print(left)
