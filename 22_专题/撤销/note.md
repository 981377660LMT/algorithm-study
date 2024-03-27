https://taodaling.github.io/blog/2020/10/11/%E4%B8%80%E7%B1%BB%E6%92%A4%E9%94%80%E9%97%AE%E9%A2%98/

## 撤销的实现方式

如果一个操作对整体的改动很小，那么我们可以记录这个操作执行后的修改，并记录在案，之后就可以很方便的撤销这些修改。
比如所我们某个操作将 x 改变为 x+1，那么我们只需要**记录这个操作**信息，之后在撤销的时候将 x 改变为 x−1 即可。
但是不是所有操作都可以通过仅记录操作信息就能撤销的，比如 chmin 操作，将某个变量 x 修改为 min(x,5)，这样我们就不能简单的撤销这样的操作了。
更加好用的方式是**记录变量被修改前的状态**，比如 x=x0，这样就能保证对任何操作都能执行撤销。

## 栈式撤销

## 队列撤销

[[Tutorial] Supporting Queue-like Undoing on DS](https://codeforces.com/blog/entry/83467)

## 优先队列撤销

[[Tutorial] Supporting Priority-Queue-like Undoing on DS](https://codeforces.com/blog/entry/111117)
