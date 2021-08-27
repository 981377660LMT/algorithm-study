from itertools import accumulate, product

presum = [0] + list(accumulate([1, 2, 3]))
print(presum)
print(len(list(product(['a', 'b', 'c'], ['d', 'e'], ['f', 'g', 'h']))))

