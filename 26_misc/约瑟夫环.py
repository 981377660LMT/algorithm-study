# josephus
def josephus(n: int, jump: int, k: int) -> int:
    """
    约瑟夫环问题(从1开始)
    n: 总人数
    jump: 每次跳过的人数
    k: 第K个被点到的人(从1开始)
    O(jump*logn)
    """
    k = n - k
    if jump <= 1:
        return n - k
    i = k
    while i < n:
        r = (i - k + jump - 2) // (jump - 1)
        if (i + r) > n:
            r = n - i
        elif r == 0:
            r = 1
        i += r
        k = (k + (r * jump)) % i
    return k + 1


if __name__ == "__main__":
    print(josephus(10, 2, 4))
