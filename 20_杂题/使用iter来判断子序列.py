def isSubsequence(source: str, target: str) -> bool:
    it = iter(target)
    return all(s in it for s in source)


print(isSubsequence('abc', 'aabbccdd'))

