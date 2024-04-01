package template.problem;

import template.binary.Log2;

import java.util.Arrays;

/**
 * <pre>
 * m场电影，第i场电影的开始时间为ai（包含）,结束时间为bi（不包含）。
 * 处理q个请求，第i个请求表示从li（包含）时间来，ri（不包含）时间离开，最多能看多少场电影。
 * </pre>
 * 
 * <pre>
 * 等价问题，给定一副包含n个顶点的有向拓扑图，之后有m条边(l_i,r_i)，每条边的长度为1，且顶点x到顶点x+1有一条长度为0的边。
 * 之后处理q个请求，第i个请求表示从li到ri的最长距离。
 * </pre>
 * 
 * <pre>
 * 预处理时间空间复杂度为O(m\log_2m)。
 * 回答每个请求的时间复杂度为O(\log_2m)
 * </pre>
 */
public class MaximumNotIntersectIntervals {
  Interval[] movies;
  int m;
  int[][] jump;

  public MaximumNotIntersectIntervals(Interval[] movies) {
    this.movies = movies;
    Arrays.sort(movies, (a, b) -> Long.compare(a.l, b.l));
    int wpos = 1;
    for (int i = 1; i < movies.length; i++) {
      while (wpos > 0 && movies[wpos - 1].r >= movies[i].r) {
        wpos--;
      }
      if (wpos == 0 || movies[wpos - 1].l < movies[i].l) {
        movies[wpos++] = movies[i];
      }
    }
    m = wpos;

    int log = Log2.floorLog(m);
    jump = new int[log + 1][m];
    for (int i = m - 1, r = m; i >= 0; i--) {
      while (r - 1 >= i && movies[r - 1].l >= movies[i].r) {
        r--;
      }
      jump[0][i] = r;
    }

    for (int i = 0; i + 1 <= log; i++) {
      for (int j = 0; j < m; j++) {
        jump[i + 1][j] = jump[i][j] == m ? m : jump[i][jump[i][j]];
      }
    }
  }

  /**
   * 查询时间区间为[l,r)，最多能看多少场电影
   */
  public int query(long l, long r) {
    if (m == 0) {
      return 0;
    }
    int lo = 0;
    int hi = m - 1;
    while (lo < hi) {
      int mid = (lo + hi) >> 1;
      if (movies[mid].l < l) {
        lo = mid + 1;
      } else {
        hi = mid;
      }
    }
    if (movies[lo].l < l || movies[lo].r > r) {
      return 0;
    }
    int start = lo;
    int ans = 1;
    for (int i = jump.length - 1; i >= 0; i--) {
      if (jump[i][start] == m || movies[jump[i][start]].r > r) {
        continue;
      }
      ans += 1 << i;
      start = jump[i][start];
    }
    return ans;
  }

  /**
   * [l, r)
   */
  public static class Interval {
    public long l;
    public long r;

    public Interval(long l, long r) {
      this.l = l;
      this.r = r;
    }

    @Override
    public String toString() {
      return String.format("[%d, %d]", l, r);
    }
  }

  @Override
  public String toString() {
    return Arrays.toString(Arrays.copyOf(movies, m));
  }
}
