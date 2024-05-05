import java.util.Arrays;
import java.util.Scanner;

class Solution {
  public void solve() {

    Scanner sc = new Scanner(System.in);
    while (sc.hasNext()) {
      int n = sc.nextInt();

      int[] pos = new int[n];
      for (int i = 0; i < n; i++) {
        pos[i] = sc.nextInt();
      }

      int[] radius = new int[n];
      for (int i = 0; i < n; i++) {
        radius[i] = sc.nextInt();
      }

      int[] value = new int[n];
      for (int i = 0; i < n; i++) {
        value[i] = sc.nextInt();
      }
    }
  }
}