def checkRotation(str1, str2):
    if len(str1) != len(str2):
        return False
    str1 = str1 + str1
    return str2 in str1
