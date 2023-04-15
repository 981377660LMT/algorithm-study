#  https://yukicoder.me/problems/no/1142
#  !两个数组的区间异或等于k的四元组个数
#  n,m<=2e5
#  nums[i],k<=1024


from typing import List


MOD = int(1e9) + 7
INV = pow(2, MOD - 2, MOD)


def xorConvolution(a: List[int], b: List[int]) -> List[int]:
    a, b = a[:], b[:]
    a = _walshHadamardTransform(a, 1)
    b = _walshHadamardTransform(b, 1)
    for i in range(len(a)):
        a[i] = a[i] * b[i]
    res = _walshHadamardTransform(a, INV)
    return res


def _walshHadamardTransform(f, op):
    n = len(f)
    l_, k = 2, 1
    while l_ <= n:
        for i in range(0, n, l_):
            for j in range(k):
                f[i + j], f[i + j + k] = (f[i + j] + f[i + j + k]) * op, (
                    f[i + j] - f[i + j + k]
                ) * op
        l_, k = l_ << 1, k << 1
    return f


if __name__ == "__main__":
    n, m, k = map(int, input().split())
    nums1 = list(map(int, input().split()))
    nums2 = list(map(int, input().split()))

    MAX = 1024
    f, g = [0] * MAX, [0] * MAX
    f[0], g[0] = 1, 1  # 统计前缀xor的频率
    xor_ = 0
    for num in nums1:
        xor_ ^= num
        f[xor_] += 1
    xor_ = 0
    for num in nums2:
        xor_ ^= num
        g[xor_] += 1

    f, g = xorConvolution(f, f), xorConvolution(g, g)
    f[0], g[0] = (f[0] - n - 1), (g[0] - m - 1)

    res = xorConvolution(f, g)
    print((res[k] // 4) % MOD)
