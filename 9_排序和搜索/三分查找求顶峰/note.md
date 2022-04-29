- 三分法可以用来查找凸函数的最大（小）值
  如果 lmid 和 rmid 在最大（小）值的`同一侧`：由于单调性，一定是二者中较大（小）的那个离最值近一些，较远的那个点对应的区间不可能包含最值，所以可以舍弃。
  如果在`两侧`：由于最值在二者中间，我们舍弃两侧的一个区间后，也不会影响最值，所以可以舍弃。
  ```Python
  lmid = left + (right - left >> 1)
  rmid = lmid + (right - lmid >> 1)  # 对右侧区间取半
  if cal(lmid) > cal(rmid):
    right = rmid
  else:
    left = lmid
  ```
