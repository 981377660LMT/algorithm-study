import java.util.HashMap;
import java.util.Map;

class Solution {
  public int maxPoints(int[][] ps) {
    int n = ps.length;
    int ans = 1;
    for (int i = 0; i < n; i++) {
      Map<String, Integer> map = new HashMap<>();
      // 由当前点 i 发出的直线所经过的最多点数量
      int max = 0;
      for (int j = i + 1; j < n; j++) {
        int x1 = ps[i][0], y1 = ps[i][1], x2 = ps[j][0], y2 = ps[j][1];
        int a = x1 - x2, b = y1 - y2;
        int k = gcd(a, b);
        String key = (a / k) + "_" + (b / k);
        map.put(key, map.getOrDefault(key, 0) + 1);
        max = Math.max(max, map.get(key));
      }
      ans = Math.max(ans, max + 1);
    }
    return ans;
  }

  static int gcd(int a, int b) {
    return b == 0 ? a : gcd(b, a % b);
  }

  public static void main(String[] args) {
    int[][] ps = { { 1, 1 }, { 3, 2 }, { 5, 3 }, { 4, 1 }, { 2, 3 }, { 1, 4 } };
    System.out.println(new Solution().maxPoints(ps));
    System.out.println(gcd(-2, -1));
  }
}
