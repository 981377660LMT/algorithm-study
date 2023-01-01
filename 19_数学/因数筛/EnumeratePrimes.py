# EnumeratePrimes


# 所有的质数
# p0=2,p1=3,p2=5,...
# 给定n,a,b 求出所有不超过n的质数p(ai+b)


from typing import List, Tuple


def faster_eratosthenes(n: int, a: int, b: int) -> Tuple[List[int], int]:
    """
    返回所有(不超过n的质数p(ai+b),不超过n的质数个数)
    n<=5e8
    """
    if n < 30:
        res = [x for x in [2, 3, 5, 7, 11, 13, 17, 19, 23, 29] if x <= n]
        return res[b::a], len(res)
    remains = [1, 7, 11, 13, 17, 19, 23, 29]
    inv_remains = {x: i for i, x in enumerate(remains)}
    msk = 255  # (1 << 8) - 1
    div30 = [i * j // 30 for j in remains for i in remains]
    mod30 = [inv_remains[i * j % 30] for j in remains for i in remains]
    msk8 = [msk - (1 << i) for i in range(8)]
    inv_msk = {1 << i: i for i in range(8)}
    res = [2, 3, 5][b::a]
    count = 3
    max_k = n // 30
    import bisect

    max_m = bisect.bisect_right(remains, n % 30) - 1
    sqrtn = int(n**0.5) + 1
    max_sqrt_k = sqrtn // 30
    max_sqrt_m = bisect.bisect_right(remains, sqrtn % 30) - 1
    table = bytearray([msk] * (max_k + 1))
    table[max_k] = (1 << (max_m + 1)) - 1
    table[0] -= 1  # remove 1
    for k in range(max_sqrt_k + 1):
        x = table[k]
        while x:
            u = x & (-x)
            if table[k] & u:
                m = inv_msk[u]
                if k == max_sqrt_k and m > max_sqrt_m:
                    break
                # k_before = k
                m_before = m
                i = k * (30 * k + 2 * remains[m]) + div30[(m << 3) + m]
                j = mod30[(m << 3) + m]
                while i < max_k or (i == max_k and j <= max_m):
                    table[i] &= msk8[j]
                    if m_before == 7:
                        i += 2 * k + remains[m] + div30[m << 3] - div30[(m << 3) + 7]
                        j = mod30[m << 3]
                        # k_before += 1
                        m_before = 0
                    else:
                        i += (
                            k * (remains[m_before + 1] - remains[m_before])
                            + div30[(m << 3) + m_before + 1]
                            - div30[(m << 3) + m_before]
                        )
                        j = mod30[(m << 3) + m_before + 1]
                        m_before += 1
            x &= x - 1
    i30 = 0
    for t in table:
        while t:
            j = inv_msk[t & (-t)]
            if count % a == b:
                res.append(i30 + remains[j])
            count += 1
            t &= t - 1
        i30 += 30
    return res, count


n, a, b = map(int, input().split())
res, count = faster_eratosthenes(n, a, b)
print(count, len(res))
print(*res)
