# https://atcoder.jp/contests/abc300/tasks/abc300_g


# k-smooth number的定义为：所有质因子都小于等于k的正整数。
# !给定正整数n和不超过100的质数p，求n以内的p-smooth number的个数。
# 1<=n<=1e16 2<=p<=100

# 1.100以内的质数不多(25个)，可以预处理出来。
# 2.折半搜索


from typing import List


def isPrime(n: int) -> bool:
    """判断n是否为质数"""
    if n < 2:
        return False
    for i in range(2, int(n**0.5) + 1):
        if n % i == 0:
            return False
    return True


def pSmoothNumber(n: int, p: int) -> int:
    def expand(nums: List[int], num: int) -> None:
        len_ = len(nums)
        for i in range(len_):
            cur = nums[i]
            while True:
                cur *= num
                if cur > n:
                    break
                nums.append(cur)

    ps = [i for i in range(2, p + 1) if isPrime(i)]
    left, right = [1], [1]
    for num in ps:
        if len(left) > len(right):  # !每次都把较短的那个数组扩展
            left, right = right, left
        expand(left, num)

    left, right = sorted(left), sorted(right)
    res = 0
    j = len(right) - 1
    for i in range(len(left)):
        ceiling = n // left[i]
        while j >= 0 and right[j] > ceiling:
            j -= 1
        if j < 0:
            break
        res += j + 1
    return res


if __name__ == "__main__":
    n, p = map(int, input().split())
    print(pSmoothNumber(n, p))
