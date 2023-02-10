// java矩阵快速幂

import java.util.ArrayList;

/**
 * 矩阵快速幂
 */
class MatPow {
  public static long[][] mul(long[][] mat1, long[][] mat2, int mod) {
    long[][] res = new long[mat1.length][mat2[0].length];
    for (int i = 0; i < mat1.length; i++) {
      for (int k = 0; k < mat2.length; k++) {
        for (int j = 0; j < mat2[0].length; j++) {
          res[i][j] = (res[i][j] + mat1[i][k] * mat2[k][j]) % mod;
          if (res[i][j] < 0) {
            res[i][j] += mod;
          }
        }
      }
    }
    return res;
  }

  private int n;
  private int mod;
  private long[] base;
  private int cacheLevel;
  private boolean useCache;
  private ArrayList<ArrayList<long[]>> cache;

  /**
   * 
   * @param base       转移矩阵,必须是方阵
   * @param mod        矩阵快速幂的模数
   * @param cacheLevel
   *                   矩阵快速幂的log底数.启用缓存时一般设置为 {@code 4}.
   *                   当调用 {@link MatPow#pow} 次数很多时,可以设置为{@code 16}.
   *                   小于 {@code 2} 时不启用缓存.
   */
  public MatPow(long[][] base, int mod, int cacheLevel) {
    this.n = base.length;
    this.mod = mod;
    this.base = new long[n * n];
    for (int i = 0; i < n; i++) {
      for (int j = 0; j < n; j++) {
        this.base[i * n + j] = base[i][j];
      }
    }
    this.cacheLevel = cacheLevel;
    this.useCache = cacheLevel >= 2;
    if (useCache) {
      this.cache = new ArrayList<ArrayList<long[]>>(cacheLevel - 1);
      for (int i = 0; i < cacheLevel - 1; i++) {
        this.cache.add(new ArrayList<long[]>(64));
      }
    }
  }

  /**
   * 时间复杂度 {@code O(n^3*logk)}.
   * 
   * @param exp 幂次/矩阵转移次数
   */
  public long[][] pow(long exp) {
    if (!useCache) {
      return powWithoutCache(exp);
    }

    if (cache.get(0).isEmpty()) {
      cache.get(0).add(base);
      for (int i = 1; i < cacheLevel - 1; i++) {
        cache.get(i).add(mul(cache.get(i - 1).get(0), base));
      }
    }

    long[] e = eye(n);
    int div = 0;
    while (exp > 0) {
      if (div == cache.get(0).size()) {
        cache.get(0).add(mul(cache.get(cacheLevel - 2).get(div - 1), cache.get(0).get(div - 1)));
        for (int i = 1; i < cacheLevel - 1; i++) {
          cache.get(i).add(mul(cache.get(i - 1).get(div), cache.get(0).get(div)));
        }
      }

      int mod = (int) (exp % cacheLevel);
      if (mod > 0) {
        e = mul(e, cache.get(mod - 1).get(div));
      }
      exp /= cacheLevel;
      div++;
    }

    return to2D(e);
  }

  private long[] mul(long[] mat1, long[] mat2) {
    long res[] = new long[n * n];
    for (int i = 0; i < n; i++) {
      for (int k = 0; k < n; k++) {
        for (int j = 0; j < n; j++) {
          res[i * n + j] = (res[i * n + j] + mat1[i * n + k] * mat2[k * n + j]) % mod;
          if (res[i * n + j] < 0) {
            res[i * n + j] += mod;
          }
        }
      }
    }
    return res;
  }

  private long[][] powWithoutCache(long exp) {
    long[] e = eye(n);
    long[] b = base.clone();
    while (exp > 0) {
      if ((exp & 1) == 1) {
        e = mul(e, b);
      }
      b = mul(b, b);
      exp >>= 1;
    }

    return to2D(e);
  }

  private long[] eye(int n) {
    long[] res = new long[n * n];
    for (int i = 0; i < n; i++) {
      res[i * n + i] = 1;
    }
    return res;
  }

  private long[][] to2D(long[] mat) {
    long[][] res = new long[n][n];
    for (int i = 0; i < n; i++) {
      for (int j = 0; j < n; j++) {
        res[i][j] = mat[i * n + j];
      }
    }
    return res;
  }
}