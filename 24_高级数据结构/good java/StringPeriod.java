package template.string;

import template.algo.DivisionDescending;
import template.string.KMPAutomaton;


public class KMPAutomaton {
  int[] data;
  int[] fail;
  int buildLast;
  public int matchLast = 0;
  int length;

  public KMPAutomaton(int cap) {
    data = new int[cap + 2];
    fail = new int[cap + 2];
    fail[0] = -1;
    buildLast = 0;
    length = cap;
  }

  public void init() {
    buildLast = 0;
  }

  /**
   * Get the border of s[0...i]
   */
  public int maxBorder(int i) {
    return fail[i + 1];
  }

  public KMPAutomaton(KMPAutomaton automaton) {
    data = automaton.data;
    fail = automaton.fail;
    buildLast = automaton.buildLast;
    length = automaton.length;
  }

  public boolean isMatch() {
    return matchLast == length;
  }

  public int length() {
    return length;
  }

  public void beginMatch() {
    matchLast = 0;
  }

  public void match(int c) {
    matchLast = visit(c, matchLast) + 1;
  }

  public int visit(int c, int trace) {
    while (trace >= 0 && data[trace + 1] != c) {
      trace = fail[trace];
    }
    return trace;
  }

  public void build(int c) {
    buildLast++;
    fail[buildLast] = visit(c, fail[buildLast - 1]) + 1;
    data[buildLast] = c;
  }

}


public class DivisionDescending {
  /**
   * <pre>
   * 给定一个check函数，存在一个数p，check(x)=true当且仅当x=kp，且check(n)=true。
   * 时间复杂度：O(n^{0.5})，以及调用O(\log_2n)次check函数。
   * </pre>
   */
  public static long find(long n, LongPredicate predicate) {
    LongArrayList list = Factorization.factorizeNumberPrime(n);
    for (long x : list.toArray()) {
      while (n % x == 0 && predicate.test(n / x)) {
        n /= x;
      }
    }
    return n;
  }

  /**
   * @see #find(long, LongPredicate)
   * @param factorization
   * @param predicate
   * @return
   */
  public static long find(LongFactorization factorization, LongPredicate predicate) {
    long n = factorization.g;
    for (long x : factorization.primes) {
      while (n % x == 0 && predicate.test(n / x)) {
        n /= x;
      }
    }
    return n;
  }
}



public class StringPeriod {
  /**
   * <pre>
   * 查询字符串的最小周期p，满足s[i+p]=s[i]，如果i+p<n
   * 时间复杂度：O(n)
   * 空间复杂度：O(n)
   * </pre>
   */
  public static int minPeriod(char[] s, int n) {
    KMPAutomaton kmp = new KMPAutomaton(n);
    for (int i = 0; i < n; i++) {
      kmp.build(s[i]);
    }
    return n - kmp.maxBorder(n);
  }

  /**
   * <pre>
   * 查询字符串的最小旋转周期p，满足s[(i+p)%n]=s[i]
   * 时间复杂度：O(n\log_2n)
   * 空间复杂度：O(\log_2n)
   * </pre>
   */
  public static int minRotatePeriod(char[] s, int n) {
    return (int) DivisionDescending.find(n, m -> {
      for (int i = 0, j = (int) m; i < n; i++, j++) {
        if (j >= n) {
          j -= n;
        }
        if (s[i] != s[j]) {
          return false;
        }
      }
      return true;
    });
  }

  /**
   * <pre>
   * 查询字符串的最小回文旋转周期p，满足s[p..n)+s[0..p)是回文
   * 前置条件：s[0..n)是回文
   * 时间复杂度：O(n\log_2n)
   * 空间复杂度：O(\log_2n)
   * </pre>
   */
  public static int minPalindromeRotatePeriod(char[] s, int n) {
    assert isPalindrome(s, n);
    int rp = minRotatePeriod(s, n);
    if (rp % 2 == 0) {
      rp /= 2;
    }
    return rp;
  }

  /**
   * <pre>
   * 判断字符串s是否是回文
   * 时间复杂度：O(n)
   * 空间复杂度：O(1)
   * </pre>
   */
  public static boolean isPalindrome(char[] s, int n) {
    int l = 0;
    int r = n - 1;
    while (l < r) {
      if (s[l] != s[r]) {
        return false;
      }
      l++;
      r--;
    }
    return true;
  }
}
