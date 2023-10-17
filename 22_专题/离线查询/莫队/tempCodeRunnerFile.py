from typing import List


def mex(nums: List[int], mexStart=0) -> int:
    n = len(nums)
    counter = [0] * (mexStart + n + 1)
    for num in nums:
        if num < mexStart or num > mexStart + n:
            continue
        counter[num - mexStart] += 1
    mex = mexStart
    while counter[mex - mexStart]:
        mex += 1
    return mex


if __name__ == "__main__":
    nums = [0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10]
    print(mex(nums, mexStart=12))
