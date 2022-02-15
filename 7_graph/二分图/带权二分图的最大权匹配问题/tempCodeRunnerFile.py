ans = 0
            # print(KM(g).solve())
            for i, j in enumerate(KM(adjMatrix).getResult()):
                if j != -1:
                    ans += adjMatrix[j][i]
            return ans