# 给定一个长度为n的正整数序列A，你需要找到两个不同的A的子序列B、C，
# 使得两个子序列中所有元素的和除以200的余数相等。你需要输出两个子序列在原数组A中的下标。
# 数据范围
# !2<=n<=200,1≤Ai≤1e9

# 鸽巢原理 2**8个子集中必定有两个子集的和模200相等

from typing import List, Tuple


def findTwoSubsequnce(nums: List[int]) -> Tuple[List[int], List[int], bool]:
    n = min(8, len(nums))
    group = [[] for _ in range(200)]
    for state in range(1, 1 << n):
        cur = []
        mod = 0
        for i in range(n):
            if state >> i & 1:
                cur.append(i)
                mod = (mod + nums[i]) % 200
        if group[mod]:
            return group[mod], cur, True
        group[mod] = cur
    return [], [], False


if __name__ == "__main__":
    n = int(input())
    nums = list(map(int, input().split()))
    res1, res2, ok = findTwoSubsequnce(nums)
    if not ok:
        print("No")
        exit(0)
    print("Yes")
    print(len(res1), *(i + 1 for i in res1))
    print(len(res2), *(i + 1 for i in res2))
