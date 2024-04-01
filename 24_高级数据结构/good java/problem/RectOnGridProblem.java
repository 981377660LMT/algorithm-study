package template.problem;

import template.datastructure.LinkedListBeta;
import template.primitve.generated.datastructure.*;
import template.primitve.generated.utils.IntegerBinaryConsumer;
import template.utils.Pair;

import java.util.ArrayList;
import java.util.Arrays;
import java.util.List;


public class RectOnGridProblem {
    /**
     * /**
     * <pre>
     * 给定二维矩阵g，找到最大面积的子矩形，要求矩形中不能包含任意的非0整数
     * </pre>
     * <pre>
     * return area, l, r, b, t
     * </pre>
     * <pre>
     * 时间复杂度为O(nm)
     * </pre>
     */
    public static int[] maximumAllZeroRectArea(Int2ToIntegerFunction mat, int n, int m) {
        int[] low = new int[m];
        int[] lb = new int[m];
        int[] rb = new int[m];
        IntegerDequeImpl dq = new IntegerDequeImpl(m);

        int best = 0;
        int left = 0;
        int right = -1;
        int up = -1;
        int bottom = 0;
        Arrays.fill(low, n);
        for (int i = n - 1; i >= 0; i--) {
            for (int j = 0; j < m; j++) {
                if (mat.apply(i, j) != 0) {
                    low[j] = i;
                }
            }
            dq.clear();
            for (int j = 0; j < m; j++) {
                while (!dq.isEmpty() && low[dq.peekLast()] >= low[j]) {
                    dq.removeLast();
                }
                lb[j] = dq.isEmpty() ? 0 : dq.peekLast() + 1;
                dq.addLast(j);
            }
            dq.clear();
            for (int j = m - 1; j >= 0; j--) {
                while (!dq.isEmpty() && low[dq.peekFirst()] >= low[j]) {
                    dq.removeFirst();
                }
                rb[j] = dq.isEmpty() ? m - 1 : dq.peekFirst() - 1;
                dq.addFirst(j);
            }
            for (int j = 0; j < m; j++) {
                int cur = (rb[j] - lb[j] + 1) * (low[j] - i);
                if (cur > best) {
                    best = cur;
                    left = lb[j];
                    right = rb[j];
                    up = i;
                    bottom = low[j] - 1;
                }
            }
        }

        return new int[]{best, left, right, bottom, up};
    }

    public static long countExactlyKOneRect(int[][] mat, int k) {
        return new ExactlyKOneRectSovler().count(mat, k);
    }

    private static class ExactlyKOneRectSovler {
        int[] left;
        int[] right;
        int[] matchSize;
        long matchPair;
        int[] size;
        void add(int x, int sign) {
            assert size[x] > 0;
            matchPair += sign * (right[x] - left[x] + 1) * matchSize[x];
            assert matchPair >= 0;
        }

        /**
         * O(max(1,k^2)nm)
         *
         * @param mat
         * @param k
         * @return
         */
        public long count(int[][] mat, int k) {
            int n = mat.length;
            int m = mat[0].length;

            if (k == 0) {
                long[][] res = countAllZeroRect((i, j) -> mat[i][j], n, m);
                long ans = 0;
                for (int i = 1; i <= n; i++) {
                    for (int j = 1; j <= m; j++) {
                        ans += res[i][j];
                    }
                }

                return ans;
            }

            int[][] next = new int[n + 1][m];
            Arrays.fill(next[n], n);
            for (int i = n - 1; i >= 0; i--) {
                for (int j = 0; j < m; j++) {
                    next[i][j] = next[i + 1][j];
                    if (mat[i][j] == 1) {
                        next[i][j] = i;
                    }
                }
            }

            IntegerMultiWayStack stack = new IntegerMultiWayStack(n, (m + 1) * k);
            size = new int[m + 1];
            left = new int[m + 1];
            right = new int[m + 1];
            matchSize = new int[m + 1];

            long ans = 0;
            for (int i = 0; i < n; i++) {
                stack.clear();
                Arrays.fill(size, 0);
                size[m] = k + 1;
                for (int j = 0; j < m; j++) {
                    int h = i;
                    while (size[j] <= k) {
                        h = next[h][j];
                        if (h >= n) {
                            break;
                        }
                        size[j]++;
                        stack.addLast(h, j);
                        h++;
                    }
                }
                int last = 0;
                for (int j = 0; j <= m; j++) {
                    if (size[j] != 0) {
                        left[last] = left[j] = last;
                        right[last] = right[j] = j;
                        last = j + 1;
                    }
                }
                Arrays.fill(matchSize, 0);
                int sum = 0;
                int consider = right[0];
                matchPair = 0;
                for (int j = right[0]; j < m; j = right[j + 1]) {
                    if (consider < j) {
                        consider = j;
                        sum = 0;
                    }
                    while (sum + size[consider] <= k) {
                        sum += size[consider];
                        consider = right[consider + 1];
                    }
                    if (sum == k) {
                        matchSize[j] = right[consider] - left[consider] + 1;
                    } else {
                        matchSize[j] = 0;
                    }
                    add(j, 1);
                    sum -= size[j];
                }

                long contrib = 0;
                for (int j = n - 1; j >= i; j--) {
                    contrib += matchPair;
                    while (!stack.isEmpty(j)) {
                        int head = stack.removeLast(j);
                        add(head, -1);
                        size[head]--;

                        int end = head;
                        if (size[head] == 0) {
                            int leftPart = left[head];
                            int rightPart = right[head + 1];
                            add(rightPart, -1);
                            end = rightPart;
                            left[leftPart] = left[rightPart] = leftPart;
                            right[leftPart] = right[rightPart] = rightPart;
                        }
                        //from last k
                        int begin = end;
                        for (int step = 0; left[begin] > 0 && step < k; step++) {
                            begin = left[begin] - 1;
                            assert size[begin] > 0;
                            add(begin, -1);
                        }

                        sum = 0;
                        consider = -1;
                        while (begin < m && begin <= end) {
                            if (consider < begin) {
                                consider = begin;
                                sum = 0;
                            }
                            while (sum + size[consider] <= k) {
                                sum += size[consider];
                                consider = right[consider + 1];
                            }
                            if (sum == k) {
                                matchSize[begin] = right[consider] - left[consider] + 1;
                            } else {
                                matchSize[begin] = 0;
                            }
                            add(begin, 1);
                            sum -= size[begin];
                            begin = right[begin + 1];
                        }
                    }
                }
                ans += contrib;
            }

            return ans;
        }
    }

    /**
     * <pre>
     * 统计不同大小的所有可行矩阵的数量，其中一个可行矩阵只能覆盖值为0的单元格
     * </pre>
     * <pre>
     * 时间复杂度O(nm)
     * </pre>
     */
    public static long[][] countAllZeroRect(Int2ToIntegerFunction mat, int n, int m) {
        long[][] tag = new long[n + 1][m + 1];
        int[] low = new int[m];
        boolean[] active = new boolean[m];
        int[] left = new int[m];
        int[] right = new int[m];
        Arrays.fill(low, n);
        LinkedListBeta<Integer> list = new LinkedListBeta<>();
        LinkedListBeta.Node<Integer>[] nodes = new LinkedListBeta.Node[m];
        for (int i = 0; i < m; i++) {
            nodes[i] = list.addLast(i);
        }
        for (int i = n - 1; i >= 0; i--) {
            for (int j = 0; j < m; j++) {
                if (mat.apply(i, j) != 0) {
                    list.remove(nodes[j]);
                    list.addLast(nodes[j]);
                    low[j] = i;
                }
            }
            Arrays.fill(active, false);
            for (int j : list) {
                int row = low[j] - 1;
                int high = row - i + 1;
                active[j] = true;
                int l = j;
                int r = j;
                tag[high][1]++;
                while (l > 0 && active[l - 1]) {
                    //merge
                    int lr = l - 1;
                    int ll = left[lr];
                    tag[high][lr - ll + 1]--;
                    tag[high][r - l + 1]--;
                    l = ll;
                    tag[high][r - l + 1]++;
                }
                while (r + 1 < m && active[r + 1]) {
                    //merge
                    int rl = r + 1;
                    int rr = right[rl];
                    tag[high][rr - rl + 1]--;
                    tag[high][r - l + 1]--;
                    r = rr;
                    tag[high][r - l + 1]++;
                }
                left[l] = left[r] = l;
                right[r] = right[l] = r;
            }
        }

        for (int i = n - 1; i >= 0; i--) {
            for (int j = 0; j <= m; j++) {
                tag[i][j] += tag[i + 1][j];
            }
        }
        for (int i = 1; i <= n; i++) {
            long cc = 0;
            long last = 0;
            for (int j = m; j >= 0; j--) {
                cc += tag[i][j];
                last += cc;
                tag[i][j] = last;
            }
        }

        for (int i = 0; i <= m; i++) {
            tag[0][i] = 0;
        }

        for (int i = 0; i <= n; i++) {
            tag[i][0] = 0;
        }

        return tag;
    }

    public static int[][] maxSquareContainsAtMostKDistinctNumbers(int[][] mat, int k) {
        int n = mat.length;
        int m = mat[0].length;
        if (k == 0) {
            return new int[n][m];
        }

        int[][] maxSquareSize = new int[n][m];
        IntegerArrayList list = new IntegerArrayList(n * m);
        for (int[] r : mat) {
            list.addAll(r);
        }
        list.unique();
        Point2[][] pts = new Point2[n][m];
        for (int i = 0; i < n; i++) {
            for (int j = 0; j < m; j++) {
                pts[i][j] = new Point2(i, j, list.binarySearch(mat[i][j]));
            }
        }
        Point2[] registries = new Point2[list.size()];
        for (int i = 0; i < n; i++) {
            for (int j = 0; j < m; j++) {
                pts[i][j].left = registries[pts[i][j].v];
                registries[pts[i][j].v] = pts[i][j];
            }
            for (int j = 0; j < m; j++) {
                registries[pts[i][j].v] = null;
            }
        }
        for (int j = 0; j < m; j++) {
            for (int i = n - 1; i >= 0; i--) {
                pts[i][j].bot = registries[pts[i][j].v];
                registries[pts[i][j].v] = pts[i][j];
            }
            for (int i = n - 1; i >= 0; i--) {
                registries[pts[i][j].v] = null;
            }
        }
        LinkedListBeta.Node<Point2>[][] forCols = new LinkedListBeta.Node[n][m];
        LinkedListBeta.Node<Point2>[][] forRows = new LinkedListBeta.Node[n][m];
        for (int i = 0; i < n; i++) {
            for (int j = 0; j < m; j++) {
                forCols[i][j] = new LinkedListBeta.Node<>(pts[i][j]);
                forRows[i][j] = new LinkedListBeta.Node<>(pts[i][j]);
            }
        }
        LinkedListBeta<Point2>[] rows = new LinkedListBeta[n];
        LinkedListBeta<Point2>[] cols = new LinkedListBeta[m];
        for (int i = 0; i < n; i++) {
            rows[i] = new LinkedListBeta<>();
        }
        for (int i = 0; i < m; i++) {
            cols[i] = new LinkedListBeta<>();
        }
        DistinctCounter dc = new DistinctCounter(new int[list.size()], 0);
        for (int i = n - 1; i >= 0; i--) {
            //update cols
            for (int j = 0; j < m; j++) {
                Point2 cur = pts[i][j];
                if (cur.bot != null) {
                    LinkedListBeta.Node<Point2> botNode = forCols[cur.bot.x][cur.bot.y];
                    if (!botNode.singleton()) {
                        cols[j].remove(botNode);
                    }
                }
                cols[j].addLast(forCols[cur.x][cur.y]);
                if (cols[j].size() > k + 1) {
                    cols[j].remove(cols[j].begin());
                }
            }

            assert dc.distinct == 0;
            assert Arrays.stream(cols).mapToInt(LinkedListBeta::size).max().orElse(-1) <= k + 1;
            for (int j = 0, r = -1; j < m; j++) {
                int lo = (r - j) + i;
                while (lo + 1 < n && r + 1 < m) {
                    lo++;
                    r++;
                    //add col
                    for (LinkedListBeta.Node<Point2> node = cols[r].begin(); node != cols[r].end(); node = node.next) {
                        Point2 pt = node.val;
                        if (pt.x <= lo) {
                            dc.modify(pt.v, 1);
                        }
                    }
                    //add row
                    for (LinkedListBeta.Node<Point2> node = rows[lo].begin(); node != rows[lo].end(); node = node.next) {
                        Point2 pt = node.val;
                        dc.modify(pt.v, 1);
                    }

                    if (dc.distinct > k) {
                        //revoke

                        //remove col
                        for (LinkedListBeta.Node<Point2> node = cols[r].begin(); node != cols[r].end(); node = node.next) {
                            Point2 pt = node.val;
                            if (pt.x <= lo) {
                                dc.modify(pt.v, -1);
                            }
                        }
                        //remove row
                        for (LinkedListBeta.Node<Point2> node = rows[lo].begin(); node != rows[lo].end(); node = node.next) {
                            Point2 pt = node.val;
                            dc.modify(pt.v, -1);
                        }

                        r--;
                        lo--;
                        break;
                    }

                    //apply
                    for (LinkedListBeta.Node<Point2> node = cols[r].begin(); node != cols[r].end(); node = node.next) {
                        Point2 pt = node.val;
                        //remove prev
                        if (pt.left != null) {
                            LinkedListBeta.Node<Point2> leftNode = forRows[pt.left.x][pt.left.y];
                            if (!leftNode.singleton()) {
                                if (pt.x <= lo) {
                                    dc.modify(pt.v, -1);
                                }
                                rows[pt.x].remove(leftNode);
                            }
                        }
                        rows[pt.x].addLast(forRows[pt.x][pt.y]);
                        if (rows[pt.x].size() > k + 1) {
                            LinkedListBeta.Node<Point2> first = rows[pt.x].begin();
                            assert pt.x > lo;
//                            dc.modify(first.val.v, -1);
                            rows[pt.x].remove(first);
                        }
                    }

                    assert Arrays.stream(rows).mapToInt(LinkedListBeta::size).max().orElse(-1) <= k + 1;
                }
                maxSquareSize[i][j] = r - j + 1;

                //remove col
                for (LinkedListBeta.Node<Point2> node = cols[j].begin(); node != cols[j].end(); node = node.next) {
                    Point2 pt = node.val;
                    if (!forRows[pt.x][pt.y].singleton()) {
                        if (pt.x <= lo) {
                            dc.modify(pt.v, -1);
                        }
                        rows[pt.x].remove(forRows[pt.x][pt.y]);
                    }
                }
                //remove row
                for (LinkedListBeta.Node<Point2> node = rows[lo].begin(); node != rows[lo].end(); node = node.next) {
                    Point2 pt = node.val;
                    dc.modify(pt.v, -1);
                }
            }
        }

        assert dc.distinct == 0;

        return maxSquareSize;
    }

    private static class Point2 {
        int x;
        int y;
        int v;
        Point2 left;
        Point2 bot;

        public Point2(int x, int y, int v) {
            this.x = x;
            this.y = y;
            this.v = v;
        }

        @Override
        public String toString() {
            return String.format("(%d, %d, %d)", x, y, v);
        }

    }

    private static class DistinctCounter {
        int[] occur;
        int distinct;

        public DistinctCounter(int[] occur, int distinct) {
            this.occur = occur;
            this.distinct = distinct;
        }

        public void modify(int i, int x) {
            if (occur[i] > 0) {
                distinct--;
            }
            occur[i] += x;
            assert occur[i] >= 0;
            if (occur[i] > 0) {
                distinct++;
            }
        }

        @Override
        public String toString() {
            StringBuilder ans = new StringBuilder();
            for (int i = 0; i < occur.length; i++) {
                if (occur[i] != 0) {
                    ans.append(i).append(',');
                }
            }
            if (ans.length() > 0) {
                ans.setLength(ans.length() - 1);
            }
            return ans.toString();
        }
    }

    /**
     * 给定一个矩阵，找到其中最大的一个子矩阵，矩阵中没有元素都互不相同，时间复杂度为$O(n^2m+nm\log_2n)$。
     *
     * @return
     */
    public static Pair<Integer, int[]> maxAreaDistinctRect(int[][] mat) {
        int best = 1;
        int rectL = 0;
        int rectR = 0;
        int rectB = 0;
        int rectU = 0;

        int n = mat.length;
        int m = mat[0].length;
        int[][] left = new int[n][m];
        int[][] right = new int[n][m];
        IntegerArrayList all = new IntegerArrayList(n * m);
        for (int i = 0; i < n; i++) {
            all.addAll(mat[i]);
        }
        all.unique();
        for (int i = 0; i < n; i++) {
            for (int j = 0; j < m; j++) {
                mat[i][j] = all.binarySearch(mat[i][j]);
            }
        }

        int k = all.size();
        int[] leftReg = new int[k];
        Arrays.fill(leftReg, -1);
        int[] rightReg = new int[k];
        Arrays.fill(rightReg, m);
        int[] size = new int[m];
        IntegerMinQueue dq = new IntegerMinQueue(m, IntegerComparator.REVERSE_ORDER);
        for (int u = n - 1; u >= 0; u--) {
            //update left and right
            for (int j = 0; j < m; j++) {
                left[u][j] = leftReg[mat[u][j]];
                leftReg[mat[u][j]] = j;
                for (int r = u + 1; r < n; r++) {
                    left[r][j] = Math.max(left[r][j], leftReg[mat[r][j]]);
                }
            }
            for (int j = 0; j < m; j++) {
                leftReg[mat[u][j]] = -1;
            }
            for (int j = m - 1; j >= 0; j--) {
                right[u][j] = rightReg[mat[u][j]];
                rightReg[mat[u][j]] = j;
                for (int r = u + 1; r < n; r++) {
                    right[r][j] = Math.min(right[r][j], rightReg[mat[r][j]]);
                }
            }
            for (int j = 0; j < m; j++) {
                rightReg[mat[u][j]] = m;
            }
            for (int j = 0; j < m; j++) {
                size[j] = m - j;
            }
            //consider top to bot
            for (int b = u; b < n; b++) {
                dq.clear();
                for (int j = m - 1; j >= 0; j--) {
                    dq.addLast(left[b][j]);
                    if (j + 1 < m) {
                        size[j] = Math.min(size[j], size[j + 1] + 1);
                    }
                    size[j] = Math.min(size[j], right[b][j] - j);
                    while (dq.size() > size[j]) {
                        dq.removeFirst();
                    }
                    while (!dq.isEmpty() && dq.min() >= j) {
                        dq.removeFirst();
                    }
                    size[j] = dq.size();
                    int area = size[j] * (b - u + 1);
                    if (area > best) {
                        best = area;
                        rectL = j;
                        rectR = rectL + size[j] - 1;
                        rectU = u;
                        rectB = b;
                    }
                }
            }
        }

        return new Pair<>(best, new int[]{rectL, rectR, rectU, rectB});
    }
}
