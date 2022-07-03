# print to file
import os


with open(os.path.join(os.path.dirname(__file__), 'output.txt'), 'w') as f:
    f.write('[' + ','.join(['"MajorityChecker"'] + ['"query"'] * 10 ** 4) + ']')
    f.write('\r')
    f.write('[' + ','.join(map(str, [[list(range(1, 10001))]] + [[1, 1, 1]] * 10 ** 4)) + ']')
