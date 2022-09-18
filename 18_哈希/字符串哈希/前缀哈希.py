from typing import List


def genHash1(s: str, base=131, mod=13331, offset=96) -> List[int]:
    """求字符串所有前缀的哈希值"""
    n = len(s)
    res = [0] * n
    for i in range(n):
        res[i] = (res[i - 1] * base + ord(s[i]) - offset) % mod
    return res


def genHash2(nums: List[int], base=131, mod=13331, offset=96) -> List[int]:
    """求数组所有前缀的哈希值"""
    n = len(nums)
    res = [0] * n
    for i in range(n):
        res[i] = (res[i - 1] * base + nums[i] - offset) % mod
    return res


if __name__ == "__main__":
    s = "abcdefg"
    print(genHash1(s))
    nums = [1, 2, 3, 4, 5, 6, 7]
    print(genHash2(nums))
