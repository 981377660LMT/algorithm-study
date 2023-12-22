# isSubstring


from typing import Any, List, Sequence, TypeVar


T = TypeVar("T")


def isSubarray(longer: Sequence[T], shorter: Sequence[T]) -> int:
    """O(n)判断shorter是否为longer的子串."""
    if len(shorter) > len(longer):
        return False
    if len(shorter) == 0:
        return True
    n, m = len(longer), len(shorter)
    st = list(shorter) + list(longer)
    z = zAlgo(st)
    for i in range(m, n + m):
        if z[i] >= m:
            return True
    return False


def zAlgo(seq: Sequence[Any]) -> List[int]:
    n = len(seq)
    if n == 0:
        return []
    z = [0] * n
    j = 0
    for i in range(1, n):
        if j + z[j] <= i:
            k = 0
        else:
            k = min(j + z[j] - i, z[i - j])
        while i + k < n and seq[k] == seq[i + k]:
            k += 1
        if j + z[j] < i + z[i]:
            j = i
        z[i] = k
    z[0] = n
    return z


if __name__ == "__main__":

    def randomCheck():
        from random import randint

        for _ in range(10000):
            s = [randint(0, 10) for _ in range(100)]
            t = [randint(0, 10) for _ in range(100)]
            if isSubarray(s, t) != ("".join(map(str, t)) in "".join(map(str, s))):
                print(s)
                print(t)
                print(isSubarray(s, t))
                print("".join(map(str, t)) in "".join(map(str, s)))
                raise ValueError

        print("pass")

    randomCheck()
