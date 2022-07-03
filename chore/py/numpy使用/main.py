import numpy as np


nums = [1, 2, 3]
nums = np.array(nums)

print(nums.tolist())


# 多项式乘法(处理生成函数)
nums1 = np.array([1, 2, 3])
nums2 = np.array([4, 5, 6])
print(nums1 * nums2)
print(np.polymul(nums1, nums2))

print(np.poly1d(nums1))
print(np.poly1d(nums2))
print(np.polymul(np.poly1d(nums1), np.poly1d(nums2)))
print(np.polymul(np.poly1d(nums1), np.poly1d(nums2)).coeffs.tolist())
