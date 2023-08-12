                spread1()  # 如果你到达安全屋后，火马上到了安全屋，这视为你能够安全到达安全屋。
                if visited[targetR][targetC] == 1:
                    return True