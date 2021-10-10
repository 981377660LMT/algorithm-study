from  sortedcontainers import SortedDict



# 左查右 右查左(参考考场就座)
class StreamSummary:
    def __init__(self):
      self.start=SortedDict() # 表示以 x 为区间左端点的区间的右端点
      self.end=SortedDict()  # 表示以 x 为区间右端点的区间的左端点

        



    # 每次 add 都会增加一个 [val, val] 的左右闭合的区间
    # 如果 add 的区间与左边或者右边能够合并，我们需要将其合并
    # 合并左右/左/右/不合并
    def add(self, val):
      if val - 1 in self.end and val + 1 in self.start:
            # [a, val-1] + [val,val] + [val+1, b] -> [a, b]
            self.end[self.start[val + 1]] = self.end[val - 1]
            self.start[self.end[val - 1]] = self.start[val + 1]
            del self.start[val + 1]
            del self.end[val - 1]
        elif val - 1 in self.end:
            # [a, val -1] + [val, val] -> [a, val]
            self.end[val] = self.end[val - 1]
            self.start[self.end[val]] = val
            del self.end[val - 1]
        elif val + 1 in self.start:
            # [val,val] + [val+1, b] -> [val, b]
            self.start[val] = self.start[val + 1]
            self.end[self.start[val]] = val
            del self.start[val + 1]
        else:
            self.start[val] = val
            self.end[val] = val
        
    # get 需要返回合并之后的区间。
    def get(self):
      res=[]
      for l,r in self.start.items():
        res.append([l,r])
      return res
        
ss=StreamSummary()
