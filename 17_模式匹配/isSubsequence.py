# O(n+m)子序列匹配


from typing import Any, Sequence


def isSubsequence(longer: str, shorter: str) -> bool:
    if len(shorter) > len(longer):
        return False
    it = iter(longer)
    return all(need in it for need in shorter)


from typing import Sequence, Any


def isSubsequence2(longer: Sequence[Any], shorter: Sequence[Any]) -> bool:
    """判断shorter是否是longer的子序列"""
    if len(shorter) > len(longer):
        return False
    if len(shorter) == 0:
        return True
    i, j = 0, 0
    while i < len(longer) and j < len(shorter):
        if longer[i] == shorter[j]:
            j += 1
            if j == len(shorter):
                return True
        i += 1
    return False


if __name__ == "__main__":
    assert isSubsequence("aabbccdd", "abc")
    assert isSubsequence2("aabbccdd", "abc")
