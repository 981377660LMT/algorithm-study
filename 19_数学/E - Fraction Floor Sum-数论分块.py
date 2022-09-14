# E - Fraction Floor Sum
# !求floor(n/1)+floor(n/2)+floor(n/3)+...+floor(n/n)的值 (n<=1e12)

# !数论分块 https://oi-wiki.org/math/number-theory/sqrt-decomposition/
# 数论分块可以O(sqrt(n))快速计算一些含有除法向下取整的和式
# 即将floor(n/i)相同的数打包同时计算


def floorSum(lower: int, upper: int, num: int) -> int:
    """
    快速计算
    sum(floor(`num`/i) for i in range(`lower`,`upper`+1))
    """
    res = 0
    upper = min(upper, num)

    start = lower
    while True:
        if start > upper:
            break
        end = min(upper, num // (num // start))
        res += (end - start + 1) * (num // start)
        start = end + 1

    return res


if __name__ == "__main__":
    n = int(input())
    print(floorSum(1, n, n))
