# 直接从1到num开始筛法 复杂度nlogn
# 这类题的特点是nums[i]<=10^5 筛法枚举因子是nlogn


from typing import List


def countMulti(nums: List[int]) -> List[int]:
    """对于每个数，原数组中有多少个他的倍数."""
    upper = max(nums) + 1
    c1, c2 = [0] * upper, [0] * upper
    for v in nums:
        c1[v] += 1
    for f in range(1, upper):
        for m in range(f, upper, f):
            c2[f] += c1[m]
    return c2


if __name__ == "__main__":
    print(countMulti([1, 2, 3, 4, 5]))  # [5, 2, 1, 1, 1]
