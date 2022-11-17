# O(n*m)子序列匹配


def isSubsequnce(longer: str, shorter: str) -> bool:
    if len(shorter) > len(longer):
        return False
    it = iter(longer)
    return all(need in it for need in shorter)


def isSubsequnce2(longer: str, shorter: str) -> bool:
    i, j = 0, 0
    while i < len(longer) and j < len(shorter):
        if longer[i] == shorter[j]:
            j += 1
        if j == len(shorter):
            return True
        i += 1
    return False


if __name__ == "__main__":
    assert isSubsequnce("aabbccdd", "abc") == True
    assert isSubsequnce2("aabbccdd", "abc") == True
