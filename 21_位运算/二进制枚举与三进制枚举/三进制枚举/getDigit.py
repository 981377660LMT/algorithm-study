def getDigit(num: int, /, *, index: int, radix: int):
    """返回 `radix` 进制下 `num` 的 `index` 位的数字，`index` 最低位(最右)为 0 """
    assert radix >= 2 and index >= 0
    prefix = num // pow(radix, index)
    return prefix % radix


print(getDigit(190, index=1, radix=10))
print(getDigit(110, index=1, radix=2))

