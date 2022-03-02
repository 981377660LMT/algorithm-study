# 查找用find/rfind而不是index/rindex，不会引发 ValueError

print('asa'.find('a'))
print('asa'.find('c'))
print('asa'.rfind('a'))


# 删除前后缀两种方法
# 删除所有前缀(strip:去除)：
# 实际上 chars 参数并非指定单个前缀；而是会移除参数值的所有组合:
print('aaaab'.lstrip('ba'))
print('aaaab'.lstrip('a'))
print('aaaab'.rstrip('b'))
print('aaaab'.strip('b'))

# str.removeprefix() ，该方法将删除单个前缀字符串，而不是全部给定集合中的字符
# print('aaaab'.removeprefix('a'))


# 转换大小写
print('Aa'.swapcase())

# 返回原字符串的标题版本，其中每个单词第一个字母为大写，其余字母为小写。
print('hello world'.title())

# 左填充0
print('as'.zfill(5))
