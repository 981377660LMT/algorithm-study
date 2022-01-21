import java.io.*;
import java.util.*;

public class meituan009小团的装饰物 {
  static int MOD = 998244353;
  static Map<String, Integer> cache = new HashMap<>();

  public static void main(String[] args) throws Exception {
    BufferedReader reader = new BufferedReader(new InputStreamReader(System.in));
    args = reader.readLine().split(" ");
    int limit = Integer.parseInt(args[0]), need = Integer.parseInt(args[1]);
    int res = dfs(1, need - 1, limit) % MOD;
    System.out.println(res);
  }

  static int dfs(int cur, int need, int limit) {
    if (need == 0)
      return 1;

    String key = cur + "," + need;
    Integer cached = cache.get(key);
    if (cached != null)
      return cached;

    int res = 0;
    for (int nextMoney = cur; nextMoney <= limit; nextMoney += cur) {
      res += dfs(nextMoney, need - 1, limit);
      res %= MOD;
    }

    cache.put(key, res);
    return res;
  }
}
