
// 提交的 JAVA 代码的类名必须为 Main 且类前不加修饰符。代码中不能包含 package 语句。
import java.io.BufferedReader;
import java.io.InputStreamReader;

// sys.setrecursionlimit(int(1e9))
// input = lambda: sys.stdin.readline().rstrip("\r\n")
// MOD = 998244353
// INF = int(4e18)

// if __name__ == "__main__":
//     MOD = 998244353
//     n = int(input()) + 1
//     fac = [1] * (n + 1)
//     ifac = [1] * (n + 1)
//     for i in range(1, n + 1):
//         fac[i] = fac[i - 1] * i % MOD
//     ifac[n] = pow(fac[n], MOD - 2, MOD)
//     for i in range(n, 0, -1):
//         ifac[i - 1] = ifac[i] * i % MOD

//     q = [1]
//     acc = [0, 1]
//     for i in range(1, n):
//         q.append(((i * q[-1] * ifac[2] + acc[-2] * ifac[i]) * ifac[i]) % MOD)
//         acc.append((acc[-1] + q[-1]) % MOD)

//     tmp = acc[-1]
//     v = pow(tmp, MOD - 2, MOD)
//     q = [num * v % MOD for num in q]
//     print(*q)

//     q = [1]
//     acc = [0, 1]
//     for i in range(1, n):
//         q.append(((i * q[-1] + acc[-2]) * ifac[i] * 2) % MOD)
//         acc.append((acc[-1] + q[-1]) % MOD)

//     tmp = acc[-1]
//     v = pow(tmp, MOD - 2, MOD)
//     q = [num * v % MOD for num in q]
//     print(*q)

public class Main {
  static int a, b;

  public static void main(String[] args) throws Exception {
    BufferedReader input = new BufferedReader(new InputStreamReader(System.in));
    int a = Integer.parseInt(input.readLine());

    String[] s = input.readLine().split(" ");
    int b = Integer.parseInt(s[0]);
    int c = Integer.parseInt(s[1]);

    String s2 = input.readLine();

    // a+b+c s2
    System.out.println(a + b + c + " " + s2);
  }

}
