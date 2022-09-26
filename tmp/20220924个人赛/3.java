// DIR4 = {
//   0: (0, 1),
//   1: (1, 0),
//   2: (0, -1),
//   3: (-1, 0),
// }  # 顺时针

// class Solution:
//   def ballGame(self, num: int, plate: List[str]) -> List[List[int]]:
//       ROW, COL = len(plate), len(plate[0])
//       queue = []
//       visited = [[False] * 4 for _ in range(ROW * COL)]
//       for i in range(ROW):
//           for j in range(COL):
//               if plate[i][j] == "O":
//                   for dir in range(4):
//                       queue.append((i, j, dir, 0))

//       res = []
//       BAD = set([(0, 0), (0, COL - 1), (ROW - 1, 0), (ROW - 1, COL - 1)])
//       while queue:
//           curRow, curCol, curDir, curStep = queue.pop()
//           hash_ = curRow * COL + curCol
//           if visited[hash_][curDir]:
//               continue
//           visited[hash_][curDir] = True

//           if curStep > num:
//               continue
//           if plate[curRow][curCol] == "W":
//               nextDir = (curDir + 1) % 4  # 顺时针
//           elif plate[curRow][curCol] == "E":
//               nextDir = (curDir - 1) % 4  # 逆时针
//           else:
//               nextDir = curDir

//           nextRow, nextCol = curRow + DIR4[nextDir][0], curCol + DIR4[nextDir][1]

//           if nextRow < 0 or nextRow >= ROW or nextCol < 0 or nextCol >= COL:
//               # 四个角除外
//               if plate[curRow][curCol] == "." and (curRow, curCol) not in BAD:
//                   res.append((curRow, curCol))
//           else:
//               hash_ = nextRow * COL + nextCol
//               if curStep + 1 <= num:
//                   queue.append((nextRow, nextCol, nextDir, curStep + 1))

//       return res

import java.util.*;

class Solution {
  public int[][] ballGame(int num, String[] plate) {

    Map<Integer, int[]> DIR4 = new HashMap<>();
    DIR4.put(0, new int[] { 0, 1 });
    DIR4.put(1, new int[] { 1, 0 });
    DIR4.put(2, new int[] { 0, -1 });
    DIR4.put(3, new int[] { -1, 0 });

    int ROW = plate.length;
    int COL = plate[0].length();
    Deque<int[]> queue = new ArrayDeque<>();
    int[][] visited = new int[ROW * COL][4];

    for (int i = 0; i < ROW; i++) {
      for (int j = 0; j < COL; j++) {
        if (plate[i].charAt(j) == 'O') {
          for (int dir = 0; dir < 4; dir++) {
            queue.offerLast(new int[] { i, j, dir, 0 });
            visited[i * COL + j][dir] = 1;
          }
        }
      }
    }

    List<List<Integer>> res = new ArrayList<>();
    Set<Integer> BAD = new HashSet<>(Arrays.asList(0, COL - 1, (ROW - 1) * COL, (ROW - 1) * COL + COL - 1));

    while (!queue.isEmpty()) {

      int[] cur = queue.pollFirst();
      int curRow = cur[0];
      int curCol = cur[1];
      int curDir = cur[2];
      int curStep = cur[3];

      if (curStep > num) {
        continue;
      }

      int nextDir;
      if (plate[curRow].charAt(curCol) == 'W') {
        nextDir = (curDir + 1) % 4; // 顺时针
      } else if (plate[curRow].charAt(curCol) == 'E') {
        nextDir = (curDir - 1 + 4) % 4; // 逆时针
      } else {
        nextDir = curDir;
      }

      int[] dir = DIR4.get(nextDir);
      int nextRow = curRow + dir[0];
      int nextCol = curCol + dir[1];
      if (nextRow < 0 || nextRow >= ROW || nextCol < 0 || nextCol >= COL) {
        // !四个角除外
        int hash = curRow * COL + curCol;
        if (!BAD.contains(hash) && plate[curRow].charAt(curCol) == '.') {
          res.add(Arrays.asList(curRow, curCol));
        }
      } else {
        int hash = nextRow * COL + nextCol;
        if (visited[hash][nextDir] == 0 && curStep + 1 <= num) {
          queue.offerLast(new int[] { nextRow, nextCol, nextDir, curStep + 1 });
          visited[hash][nextDir] = 1;
        }
      }
    }

    int[][] ans = new int[res.size()][2];
    for (int i = 0; i < res.size(); i++) {
      ans[i][0] = res.get(i).get(0);
      ans[i][1] = res.get(i).get(1);
    }
    return ans;
  }
}
