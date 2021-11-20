nums = [[1], [23], [4]]
print(sum(nums, []))
#########################################

from functools import reduce
from operator import iconcat

print(reduce(iconcat, nums, []))

