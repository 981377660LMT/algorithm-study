import java.util.*;

public class Main {

  static int[][][] ROLE = {
      { { 0, 1 }, { 0, -1 }, { -1, 0 }, { 1, 0 } },
      { { 1, 2 }, { 1, -2 }, { -1, 2 }, { -1, -2 }, { 2, 1 }, { 2, -1 }, { -2, 1 }, { -2, -1 } },
  };

  public static void main(String[] args) {
    Scanner sc = new Scanner(System.in);
    int ROW = sc.nextInt(), COL = sc.nextInt();
    sc.nextLine();
    char[][] grid = new char[ROW][COL];
    for (int i = 0; i < ROW; i++) {
      grid[i] = sc.nextLine().toCharArray();
    }

    int sr = 0, sc = 0, er = ROW - 1, ec = COL - 1;
    Queue<int[]> queue = new LinkedList<>();
    queue.add(new int[] { sr, sc, 0, 0 }); // (row, col, step, role)
    Set<String> visited = new HashSet<>();
    visited.add(sr + "," + sc + "," + 0);
    while (!queue.isEmpty()) {
      int[] cur = queue.poll();
      int curRow = cur[0], curCol = cur[1], step = cur[2], role = cur[3];
      if (curRow == er && curCol == ec) {
        System.out.println(step);
        System.exit(0);
      }

      if (grid[curRow][curCol] == 'S') {
        int nextRole = role ^ 1;
        String nextState = curRow + "," + curCol + "," + nextRole;
        if (!visited.contains(nextState)) {
          queue.add(new int[] { curRow, curCol, step + 1, nextRole });
          visited.add(nextState);
        }
      }

      for (int[] drdc : ROLE[role]) {
        int nr = curRow + drdc[0], nc = curCol + drdc[1];
        if (0 <= nr && nr < ROW && 0 <= nc && nc < COL && grid[nr][nc] != 'X') {
          String nextState = nr + "," + nc + "," + role;
          if (!visited.contains(nextState)) {
            queue.add(new int[] { nr, nc, step + 1, role });
            visited.add(nextState);
          }
        }
      }
    }

    System.out.println(-1);
  }

}