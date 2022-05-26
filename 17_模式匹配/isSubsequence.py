def isSubsequnce(longer: str, shorter: str) -> bool:
    it = iter(longer)
    return all(need in it for need in shorter)


def isSubsequnce2(longer: str, shorter: str) -> bool:
    j = 0
    for i in range(len(longer)):
        if longer[i] == shorter[j]:
            j += 1
        if j == len(shorter):
            return True
    return False


if __name__ == '__main__':
    assert isSubsequnce('aabbccdd', 'abc') == True
    assert isSubsequnce2('aabbccdd', 'abc') == True

