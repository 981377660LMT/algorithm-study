# O(n^2)子串匹配
def isSubarray(longer: str, shorter: str) -> bool:
    i, j = 0, 0
    while i < len(longer) and j < len(shorter):
        if longer[i] == shorter[j]:
            j += 1
            i += 1
        else:
            break
        if j == len(shorter):
            return True
    return False
