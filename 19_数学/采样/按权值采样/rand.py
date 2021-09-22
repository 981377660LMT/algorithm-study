import random
from itertools import accumulate

print(random.randint(1, 4))  # [1,4]
print(*accumulate([1, 2, 3]))  # 1,3,6
