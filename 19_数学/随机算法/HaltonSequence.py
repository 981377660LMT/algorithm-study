# HaltonSequence(高维序列)
# https://blog.csdn.net/Amber_amber/article/details/47421053
# https://github.com/hos-lyric/libra/blob/60b8b56ecae5860f81d75de28510d94336f5dad9/other/halton.cpp#L3
# Halton序列是用于为蒙特卡罗模拟等数值方法生成空间点的序列.
# !选取一个质数base作为基数，然后不断地切分，从而形成一些不重复并且均匀的点，每个点的坐标都在0~1之间。
# !对区间有一个很好的覆盖度

# 1/2 -> 0.1
# 1/4 -> 0.01
# 3/4 -> 0.11
# 1/8 -> 0.001
# 5/8 -> 0.101
# 3/8 -> 0.011
# 7/8 -> 0.111
# 1/16 -> 0.0001
# 9/16 -> 0.1001
# !可以看出来，小数点后的部分正是二进制正整数序列1, 10, 11, 100, 101, 110, 111, 1000, 1001……每个数的各位逆序。


def genHaltonSequnce(baseP: int, exp: int):
    """
    生成分母在[1, baseP^exp)间的Halton序列.
    https://www.zhihu.com/question/28157920/answer/1115592808
    """
    ds = [0] * exp
    xs = [0] * exp
    b = 1
    for i in range(exp):
        b *= baseP
    ds[exp - 1] = 1 + baseP
    for i in range(exp - 2, -1, -1):
        ds[i] = ds[i + 1] * baseP
    for i in range(exp):
        ds[i] -= b

    y = 0
    while True:
        for i in range(exp):
            xs[i] += 1
            if xs[i] == baseP:
                xs[i] = 0
            else:
                y += ds[i]
                break
        yield y / b


if __name__ == "__main__":
    h2 = genHaltonSequnce(2, 6)
    # [0.5, 0.25, 0.75, 0.125, 0.625, 0.375, 0.875, 0.0625, 0.5625, 0.3125]
    print([next(h2) for _ in range(100)])
