# 1.检查一个数是否存在

visitedState = 187
print((visitedState >> 7) & 1)


# 2.检查一组数是否存在

# 检查是否都看过(都存在于集合，交小并大)
visitedState = 187
print(visitedState | 3 == visitedState)

# 检查是否都没看过(都不存在于集合，交小并大): availabelState反向模式
availabelState = 187
print(availabelState | 3 == availabelState)
