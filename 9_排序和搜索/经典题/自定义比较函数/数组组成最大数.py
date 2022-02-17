from functools import cmp_to_key

# 给定一组非负整数，重新排列它们的顺序使之组成一个最大的整数。
# 数组组成最大数
nums = input().replace('[', '').replace(']', '').split(',')

nums = sorted(nums, key=cmp_to_key(lambda x, y: int(y + x) - int(x + y)))

print(''.join(nums))


# cmp_to_key 将compare函数转换成key


# 示例 1：


# 输入：[10,1,2]
# 输出：2110
# 示例 2：


# 输入：[3,30,34,5,9]
# 输出：9534330
