def isSubsequence(longer: str, shorter: str) -> bool:
    if len(shorter) > len(longer):
        return False
    it = iter(longer)
    return all(need in it for need in shorter)


s = input()
res = ""
while isSubsequence(s, res):
    res += "a"
print(res)
