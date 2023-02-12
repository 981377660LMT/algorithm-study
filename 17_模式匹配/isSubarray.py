# O(n^2)子串匹配,一般可以用语言自带的indexOf/find 函数
def isSubarray(longer: str, shorter: str) -> bool:
    """判断shorter是否是longer的子串"""
    if len(shorter) > len(longer):
        return False
    if len(shorter) == 0:
        return True
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
