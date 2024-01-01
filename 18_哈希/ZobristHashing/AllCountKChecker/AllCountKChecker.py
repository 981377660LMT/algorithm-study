# 判断数据结构中每个数出现的次数是否均为k.
# 等价于:
#  1. 数据结构中每个数出现的次数均为k的倍数：异或哈希.
#  2. 数据结构中每个数出现的次数均不超过k：双指针.
#     在右指针扫到 i 的时候，不停将左指针向右移动并减去这个桶的出现次数，
#     直到 nums[i] 的出现次数小于等于 k 为止。此时再统计答案，两个限制都可以满足。


from typing import List
from collections import defaultdict
from random import randint


def countSubarrayWithFrequencyEqualToK(arr: List[int], k: int) -> int:
    """统计满足`每个元素出现的次数均为k`条件的子数组的个数."""
    n = len(arr)
    if n == 0 or k <= 0 or k > n:
        return 0

    pool = defaultdict(lambda: randint(1, (1 << 61) - 1))
    id_ = defaultdict(lambda: len(id_))
    arr = [id_[v] for v in arr]
    counter = [0] * len(id_)
    random = [pool[v] for v in arr]
    hashPreSum = [0] * (n + 1)  # 哈希之和的前缀和
    for i, v in enumerate(arr):
        hashPreSum[i + 1] = hashPreSum[i]
        hashPreSum[i + 1] -= counter[v] * random[i]
        counter[v] = (counter[v] + 1) % k
        hashPreSum[i + 1] += counter[v] * random[i]

    countPreSum = defaultdict(int, {0: 1})
    counter = [0] * len(id_)
    res, left = 0, 0
    for right, num in enumerate(arr):
        counter[num] += 1
        while counter[num] > k:
            counter[arr[left]] -= 1
            countPreSum[hashPreSum[left]] -= 1
            left += 1
        res += countPreSum[hashPreSum[right + 1]]
        countPreSum[hashPreSum[right + 1]] += 1
    return res


if __name__ == "__main__":
    # https://leetcode.cn/problems/count-complete-substrings

    class Solution:
        def countCompleteSubstrings(self, word: str, k: int) -> int:
            n = len(word)
            ords = [ord(x) - 97 for x in word]
            groups = []
            ptr = 0
            while ptr < n:
                leader = ords[ptr]
                group = [leader]
                ptr += 1
                while ptr < n and abs(ords[ptr] - ords[ptr - 1]) <= 2:
                    group.append(ords[ptr])
                    ptr += 1
                groups.append(group)
            return sum(countSubarrayWithFrequencyEqualToK(group, k) for group in groups)

    # https://www.luogu.com.cn/problem/CF1418G
    import sys

    input = lambda: sys.stdin.readline().rstrip("\r\n")

    n = int(input())
    arr = list(map(int, input().split()))
    print(countSubarrayWithFrequencyEqualToK(arr, 3))
