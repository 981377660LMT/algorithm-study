# F. The Number of Subpermutations 题解(异或hash维护全排列)
# https://www.luogu.com.cn/problem/CF1175F
# https://blog.51cto.com/u_15077535/3585174
# !求一个序列中是排列的子串数量。
# !枚举 1的位置。每次向左/右搜，往一个方向搜的时候记录扫到的最大值，
# 当作排列的长度（最大值），在判断这个子串是否合法就行了。


from collections import defaultdict
from random import randint
from typing import List


def numberOfSubpermutations(nums: List[int]) -> int:
    """求一个序列中是排列的子数组数."""
    n = len(nums)
    pool = defaultdict(lambda: randint(1, (1 << 61) - 1))
    targets = [0] * (n + 1)  # !targets[i] 1~i 的集合的哈希值
    for i in range(1, n + 1):
        targets[i] = targets[i - 1] ^ pool[i]

    def cal(curNums: List[int]) -> int:
        """最大值在1的右侧时满足题意的子数组数."""
        preXor = [0] * (n + 1)
        for i in range(1, n + 1):
            preXor[i] = preXor[i - 1] ^ pool[curNums[i - 1]]

        res = 0
        for i, v in enumerate(curNums):
            if v != 1:
                continue
            max_ = 1
            for j in range(i + 1, n):
                if curNums[j] == 1:
                    break
                max_ = curNums[j] if curNums[j] > max_ else max_
                if j + 1 - max_ >= 0 and preXor[j + 1] ^ preXor[j + 1 - max_] == targets[max_]:
                    res += 1
        return res

    return cal(nums) + cal(nums[::-1]) + nums.count(1)


if __name__ == "__main__":
    _ = int(input())
    nums = list(map(int, input().split()))
    print(numberOfSubpermutations(nums))
