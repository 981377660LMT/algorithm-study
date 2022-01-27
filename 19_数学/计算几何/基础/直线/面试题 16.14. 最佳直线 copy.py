from typing import List
from collections import defaultdict


def gcd(a, b):
    return a if b == 0 else gcd(b, a % b)


class Solution:
    def bestLine(self, points: List[List[int]]) -> List[int]:
        n = len(points)
        res = []
        maxCount = 0

        for i in range(n):  # 第一层for循环
            x1, y1 = points[i]
            counter = defaultdict(int)
            firstpair = defaultdict(list)

            for j in range(i + 1, n):  # 第二层for循环
                x2, y2 = points[j]
                A, B = (y2 - y1), (x2 - x1)

                if B == 0:
                    key = (0, 0)
                else:
                    gcd_ = gcd(A, B)
                    key = (A / gcd_, B / gcd_)

                counter[key] += 1
                firstpair.setdefault(key, [i, j])
                if counter[key] > maxCount:  # 只有更多，才更新
                    print(key)
                    maxCount = counter[key]
                    res = firstpair[key]
        return res


print(
    Solution().bestLine(
        [
            [-24272, -29606],
            [-37644, -4251],
            [2691, -22513],
            [-14592, -33765],
            [-21858, 28550],
            [-22264, 41303],
            [-6960, 12785],
            [-39133, -41833],
            [25151, -26643],
            [-19416, 28550],
            [-17420, 22270],
            [-8793, 16457],
            [-4303, -25680],
            [-14405, 26607],
            [-49083, -26336],
            [22629, 20544],
            [-23939, -25038],
            [-40441, -26962],
            [-29484, -30503],
            [-32927, -18287],
            [-13312, -22513],
            [15026, 12965],
            [-16361, -23282],
            [7296, -15750],
            [-11690, -21723],
            [-34850, -25928],
            [-14933, -16169],
            [23459, -9358],
            [-45719, -13202],
            [-26868, 28550],
            [4627, 16457],
            [-7296, -27760],
            [-32230, 8174],
            [-28233, -8627],
            [-26520, 28550],
            [5515, -26001],
            [-16766, 28550],
            [21888, -3740],
            [1251, 28550],
            [15333, -26322],
            [-27677, -19790],
            [20311, 7075],
            [-10751, 16457],
            [-47762, -44638],
            [20991, 24942],
            [-19056, -11105],
            [-26639, 28550],
            [-19862, 16457],
            [-27506, -4251],
            [-20172, -5440],
            [-33757, -24717],
            [-9411, -17379],
            [12493, 29906],
            [0, -21755],
            [-36885, -16192],
            [-38195, -40088],
            [-40079, 7667],
            [-29294, -34032],
            [-55968, 23947],
            [-22724, -22513],
            [20362, -11530],
            [-11817, -23957],
            [-33742, 5259],
            [-10350, -4251],
            [-11690, -22513],
            [-20241, -22513],
        ]
    )
)
