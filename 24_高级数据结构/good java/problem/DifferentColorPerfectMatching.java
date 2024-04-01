package template.problem;

import template.datastructure.LinkedListBeta;

import java.util.Arrays;

/**
 * 给定一个二分图$(A,B)$，其中$A$与$B$中均包含$n$个顶点。每个顶点都有各自的颜色，如果$a\in A$与$b\in B$拥有不同的颜色，则二者之间建立一条边。现在要求判断是否存在完美匹配，如果存在则找到一组。
 */
public class DifferentColorPerfectMatching {
    Node[] na;
    Node[] nb;
    LinkedListBeta<Void>[] pqA;
    LinkedListBeta<Void>[] pqB;
    int thresholdA;
    int thresholdB;
    int top;
    int[] regA;
    int[] regB;
    int[] nextA;
    int[] nextB;
    int n;
    int m;

    Node firstA;
    Node firstB;

    boolean hasNextA() {
        if (firstA != null) {
            return true;
        }
        while (thresholdA < top && pqA[thresholdA].isEmpty()) {
            thresholdA++;
        }
        if (thresholdA < top) {
            firstA = (Node) pqA[thresholdA].begin();
            pqA[thresholdA].remove(firstA);
        }
        return firstA != null;
    }

    Node popA() {
        Node ans = firstA;
        firstA = null;
        return ans;
    }

    void addBackA(Node x) {
        if (x.w >= thresholdA) {
            pqA[x.w].addLast(x);
        } else {
            if (firstA != null) {
                pqA[firstA.w].addLast(firstA);
            }
            firstA = x;
        }
    }

    void removeA(Node x){
        if(firstA == x){
            firstA = null;
        }else{
            pqA[x.w].remove(x);
        }
    }

    void removeB(Node x){
        if(firstB == x){
            firstB = null;
        }else{
            pqB[x.w].remove(x);
        }
    }

    void addBackB(Node x) {
        if (x.w >= thresholdB) {
            pqB[x.w].addLast(x);
        } else {
            if (firstB != null) {
                pqB[firstB.w].addLast(firstB);
            }
            firstB = x;
        }
    }

    boolean hasNextB() {
        if (firstB != null) {
            return true;
        }
        while (thresholdB < top && pqB[thresholdB].isEmpty()) {
            thresholdB++;
        }
        if (thresholdB < top) {
            firstB = (Node) pqB[thresholdB].begin();
            pqB[thresholdB].remove(firstB);
        }
        return firstB != null;
    }

    Node popB() {
        Node ans = firstB;
        firstB = null;
        return ans;
    }

    /**
     * m means the max color, and n is the max nodes
     *
     * @param n
     * @param m
     */
    public DifferentColorPerfectMatching(int n, int m) {
        this.n = n;
        this.m = m;
        regA = new int[m];
        regB = new int[m];
        nextA = new int[n];
        nextB = new int[n];
        Arrays.fill(regA, n);
        Arrays.fill(regB, n);
        na = new Node[m];
        nb = new Node[m];
        for (int i = 0; i < m; i++) {
            na[i] = new Node(null);
            nb[i] = new Node(null);
            na[i].c = i;
            nb[i].c = i;
            na[i].init();
            nb[i].init();
        }
        pqA = new LinkedListBeta[n * 2];
        pqB = new LinkedListBeta[n * 2];
        for (int i = 0; i < n * 2; i++) {
            pqA[i] = new LinkedListBeta<>();
            pqB[i] = new LinkedListBeta<>();
        }
    }

    /**
     * sol[i] means (i, sol[i]) is matched.
     * <p>
     * O(a.length + b.length)
     */
    public boolean solve(int[] a, int[] b, int[] sol) {
        if (a.length != b.length) {
            return false;
        }
        top = a.length;
        thresholdA = 0;
        thresholdB = 0;
        Arrays.fill(nextA, n);
        Arrays.fill(nextB, n);
        for (int i = a.length - 1; i >= 0; i--) {
            int x = a[i];
            na[x].cnt++;
            nextA[i] = regA[x];
            regA[x] = i;
        }
        for (int i = b.length - 1; i >= 0; i--) {
            int x = b[i];
            nb[x].cnt++;
            nextB[i] = regB[x];
            regB[x] = i;
        }

        boolean ans = true;
        for (int x : a) {
            if (na[x].cnt + nb[x].cnt > a.length) {
                ans = false;
            }
        }

        if (ans) {
            for (int x : a) {
                if (na[x].visited) {
                    continue;
                }
                na[x].visited = true;
                na[x].w = a.length - na[x].cnt - nb[x].cnt;
                pqA[na[x].w].addLast(na[x]);
            }
            for (int x : b) {
                if (nb[x].visited) {
                    continue;
                }
                nb[x].visited = true;
                nb[x].w = a.length - na[x].cnt - nb[x].cnt;
                pqB[nb[x].w].addLast(nb[x]);
            }

            while (hasNextA()) {
                Node h0 = popA();
                Node addBack0 = null;
                if (nb[h0.c].w >= 0) {
                    addBack0 = nb[h0.c];
                    removeB(addBack0);
                }
                hasNextB();
                Node h1 = popB();
                Node addBack1 = null;
                if (na[h1.c].w >= 0) {
                    addBack1 = na[h1.c];
                    removeA(addBack1);
                }

                int index0 = regA[h0.c];
                int index1 = regB[h1.c];
                regA[h0.c] = nextA[index0];
                regB[h1.c] = nextB[index1];

                h0.w++;
                h1.w++;
                if (addBack0 != null) {
                    addBack0.w++;
                }
                if (addBack1 != null) {
                    addBack1.w++;
                }

                top++;
                sol[index0] = index1;
                if (regA[h0.c] < a.length) {
                    addBackA(h0);
                } else {
                    h0.w = -1;
                }
                if (addBack0 != null) {
                    addBackB(addBack0);
                }
                if (addBack1 != null) {
                    addBackA(addBack1);
                }
                if (regB[h1.c] < b.length) {
                    addBackB(h1);
                } else {
                    h1.w = -1;
                }
            }
        }

        for (int x : a) {
            na[x].init();
            regA[x] = n;
        }
        for (int x : b) {
            nb[x].init();
            regB[x] = n;
        }

        return ans;
    }

    static class Node extends LinkedListBeta.Node<Void> {

        int cnt;
        int w;
        int c;
        boolean visited;

        public Node(Void val) {
            super(val);
        }

        void init() {
            cnt = 0;
            w = -1;
            visited = false;
        }
    }
}
