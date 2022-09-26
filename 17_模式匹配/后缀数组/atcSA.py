from typing import List, Sequence, Tuple


def sa_is(s: Sequence[int], upper: int) -> List[int]:
    """SA-IS, linear-time suffix array construction

    Args:
        s (Sequence[int]): Sequence of integers in [0, upper]
        upper (int): Upper bound of the integers in s

    Returns:
        List[int]: Suffix array
    """

    n = len(s)
    if n == 0:
        return []
    if n == 1:
        return [0]
    if n == 2:
        return [0, 1] if s[0] < s[1] else [1, 0]

    sa = [0] * n
    ls = [False] * n
    for i in range(n - 2, -1, -1):
        ls[i] = ls[i + 1] if s[i] == s[i + 1] else (s[i] < s[i + 1])
    sum_l = [0] * (upper + 1)
    sum_s = [0] * (upper + 1)
    for i in range(n):
        if not ls[i]:
            sum_s[s[i]] += 1
        else:
            sum_l[s[i] + 1] += 1
    for i in range(upper + 1):
        sum_s[i] += sum_l[i]
        if i < upper:
            sum_l[i + 1] += sum_s[i]

    def induce(lms):
        for i in range(n):
            sa[i] = -1
        buf = sum_s[:]
        for d in lms:
            if d == n:
                continue
            sa[buf[s[d]]] = d
            buf[s[d]] += 1
        buf = sum_l[:]
        sa[buf[s[n - 1]]] = n - 1
        buf[s[n - 1]] += 1
        for i in range(n):
            v = sa[i]
            if v >= 1 and not ls[v - 1]:
                sa[buf[s[v - 1]]] = v - 1
                buf[s[v - 1]] += 1
        buf = sum_l[:]
        for i in range(n - 1, -1, -1):
            v = sa[i]
            if v >= 1 and ls[v - 1]:
                buf[s[v - 1] + 1] -= 1
                sa[buf[s[v - 1] + 1]] = v - 1

    lms_map = [-1] * (n + 1)
    m = 0
    for i in range(1, n):
        if not ls[i - 1] and ls[i]:
            lms_map[i] = m
            m += 1
    lms = []
    for i in range(1, n):
        if not ls[i - 1] and ls[i]:
            lms.append(i)
    induce(lms)

    if m:
        sorted_lms = []
        for v in sa:
            if lms_map[v] != -1:
                sorted_lms.append(v)
        rec_s = [0] * m
        rec_upper = 0
        rec_s[lms_map[sorted_lms[0]]] = 0
        for i in range(1, m):
            l, r = sorted_lms[i - 1], sorted_lms[i]
            end_l = lms[lms_map[l] + 1] if lms_map[l] + 1 < m else n
            end_r = lms[lms_map[r] + 1] if lms_map[r] + 1 < m else n
            same = True
            if end_l - l != end_r - r:
                same = False
            else:
                while l < end_l:
                    if s[l] != s[r]:
                        break
                    l += 1
                    r += 1
                if l == n or s[l] != s[r]:
                    same = False
            if not same:
                rec_upper += 1
            rec_s[lms_map[sorted_lms[i]]] = rec_upper
        rec_sa = sa_is(rec_s, rec_upper)
        for i in range(m):
            sorted_lms[i] = lms[rec_sa[i]]
        induce(sorted_lms)
    return sa


def rank_lcp(s: Sequence[int], sa: List[int]) -> Tuple[List[int], List[int]]:
    """Rank and LCP array construction

    Args:
        s (Sequence[int]): Sequence of integers in [0, upper]
        sa (List[int]): Suffix array

    Returns:
        Tuple[List[int], List[int]]: Rank array and LCP array

    example:
    ```
    ords = [1, 2, 3, 1, 2, 3]
    sa = sa_is(ords, max(ords))
    rank, lcp = rank_lcp(ords, sa)
    print(rank, lcp)  # [1, 3, 5, 0, 2, 4] [3, 0, 2, 0, 1]
    ```
    """
    n = len(s)
    assert n >= 1
    rank = [0] * n
    for i in range(n):
        rank[sa[i]] = i
    lcp = [0] * (n - 1)
    h = 0
    for i in range(n):
        if h > 0:
            h -= 1
        if rank[i] == 0:
            continue
        j = sa[rank[i] - 1]
        while j + h < n and i + h < n:
            if s[j + h] != s[i + h]:
                break
            h += 1
        lcp[rank[i] - 1] = h
    return rank, lcp


if __name__ == "__main__":
    # !求每个后缀与所有后缀的公共前缀长度和
    n = int(input())
    S = str(input())
    ords = [ord(c) for c in S]
    sa = sa_is(ords, max(ords))
    _, lcp = rank_lcp(ords, sa)
    # print(sa)
    # print(lcp)

    ans = [0] * n
    stack = []
    cur = 0
    for i in range(n - 1):
        length = 1
        while stack and stack[-1][0] >= lcp[i]:
            a, l = stack.pop()
            length += l
            cur -= a * l
        cur += lcp[i] * length
        stack.append((lcp[i], length))
        ans[sa[i + 1]] += cur

    sa.reverse()
    lcp.reverse()
    stack = []
    cur = 0
    for i in range(n - 1):
        length = 1
        while stack and stack[-1][0] >= lcp[i]:
            a, l = stack.pop()
            length += l
            cur -= a * l
        cur += lcp[i] * length
        stack.append((lcp[i], length))
        ans[sa[i + 1]] += cur

    for i in range(n):
        ans[i] += n - i
    print(*ans, sep="\n")
