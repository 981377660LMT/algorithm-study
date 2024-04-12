package template.algo;

import template.math.Factorization;
import template.math.LongFactorization;
import template.primitve.generated.datastructure.LongArrayList;

import java.util.function.LongPredicate;


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
