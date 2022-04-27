def isSubsequnce(pattern: str, needle: str) -> bool:
    it = iter(pattern)
    return all(p in it for p in needle)


def isSubsequnce2(pattern: str, needle: str) -> bool:
    j = 0
    for i in range(len(pattern)):
        if pattern[i] == needle[j]:
            j += 1
        if j == len(needle):
            return True
    return False


if __name__ == '__main__':
    assert isSubsequnce('aabbccdd', 'abc') == True
    assert isSubsequnce2('aabbccdd', 'abc') == True

