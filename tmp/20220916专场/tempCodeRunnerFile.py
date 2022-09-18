              number -= count
                lists.discard((price, expire, count))
                res += price * count

        # print(res, 999)
        disCountRes = res * self.userCount[customer] / 100
        self.userCount[customer] = max(70, self.userCount[customer] - 1)
        return disCountRes
