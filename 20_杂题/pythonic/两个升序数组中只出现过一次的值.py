nums1 = [1, 2, 3]
nums2 = [2, 3, 4]
S1, S2 = set(nums1), set(nums2)
print(*sorted(S1 ^ S2))  # 1 4

# !symmetric_difference 等于异或操作
# !把set的运算想象成大数位运算即可
