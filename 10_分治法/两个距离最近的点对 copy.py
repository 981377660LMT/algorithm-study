# nlogn解法
# 2 ≤ n ≤ 100,000


class Solution:
    def minDist(self, coordinates):
        def dist(p1, p2):
            return abs(p1[0] - p2[0]) ** 2 + abs(p1[1] - p2[1]) ** 2

        def bruteForce(points):
            res = int(1e20)
            for i in range(len(points)):
                for j in range(i + 1, len(points)):
                    res = min(res, dist(points[i], points[j]))
            return res

        def closestInStrip(points, d):
            for i in range(len(points)):
                for j in range(i + 1, len(points)):
                    if points[j][1] - points[i][1] >= d:
                        break
                    d = min(d, dist(points[i], points[j]))
            return d

        def closest(xPoints, yPoints):
            if len(xPoints) <= 3:
                return bruteForce(xPoints)
            mid = len(xPoints) // 2
            midPoint = xPoints[mid]
            yLeft = [point for point in yPoints if point[0] <= midPoint[0]]
            yRight = [point for point in yPoints if point[0] > midPoint[0]]
            leftMin = closest(xPoints[:mid], yLeft)
            rightMin = closest(xPoints[mid:], yRight)
            d = min(leftMin, rightMin)
            strip = [point for point in yPoints if abs(point[0] - midPoint[0]) < d]
            return closestInStrip(strip, d)

        xPoints = sorted(coordinates, key=lambda x: x[0])
        yPoints = sorted(coordinates, key=lambda x: x[1])
        return closest(xPoints, yPoints)


print(Solution().minDist(coordinates=[[0, 0], [1, 1], [2, 4]]))
