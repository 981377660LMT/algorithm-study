package template.problem;

public class BinaryTreeTraversalProblem {
    public static Node travelByPreAndIn(int[] preOrder, int[] inOrder) {
        int n = preOrder.length;
        int[] inRev = new int[n];
        for (int i = 0; i < n; i++) {
            inRev[inOrder[i]] = i;
        }
        return buildByPreAndIn(preOrder, inRev, 0, n - 1, 0);
    }

    private static Node buildByPreAndIn(int[] preOrder, int[] invInOrder, int l, int r, int ll) {
        if (l > r) {
            return null;
        }
        Node root = new Node(preOrder[l]);
        int left = invInOrder[root.id] - ll;
        root.left = buildByPreAndIn(preOrder, invInOrder, l + 1, l + left, ll);
        root.right = buildByPreAndIn(preOrder, invInOrder, l + 1 + left, r, ll + left + 1);
        return root;
    }


    public static class Node {
        public Node left;
        public Node right;
        public int id;

        public Node(int id) {
            this.id = id;
        }
    }
}
