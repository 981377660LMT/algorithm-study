# C-时针分针的夹角与距离(不超过180度的夹角)


from math import cos, radians


def calDegree(hour: int, minute: int) -> float:
    minuteAngle = minute * 6
    hourAngle = hour * 30 + minute * 0.5
    return min(abs(minuteAngle - hourAngle), 360 - abs(minuteAngle - hourAngle))


def calDist(hourLength: int, minuteLength: int, degree: float) -> float:
    return (
        hourLength**2 + minuteLength**2 - 2 * hourLength * minuteLength * cos(radians(degree))
    ) ** 0.5


if __name__ == "__main__":
    A, B, H, M = map(int, input().split())
    angle = calDegree(H, M)
    print(calDist(A, B, angle))
