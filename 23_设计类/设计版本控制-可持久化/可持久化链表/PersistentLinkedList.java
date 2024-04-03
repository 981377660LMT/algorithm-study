package template.datastructure;

import java.util.Random;

public class PersistentLinkedList<T> {
    private PersistentTreap data;

    private PersistentLinkedList(PersistentTreap data) {
        this.data = data;
    }

    public PersistentLinkedList() {
        this(PersistentTreap.NIL);
    }

    public int size() {
        return data.size;
    }

    /**
     * Get the i-th element
     */
    public T get(int i) {
        return (T) PersistentTreap.getValueByRank(data, i);
    }

    /**
     * Insert before the i-the element
     */
    public PersistentLinkedList<T> insert(int i, T val) {
        PersistentTreap[] split = PersistentTreap.splitByRank(data, i - 1);
        PersistentTreap newNoe = new PersistentTreap();
        newNoe.value = val;
        newNoe.pushUp();
        split[0] = PersistentTreap.merge(split[0], newNoe);

        PersistentLinkedList<T> ans = new PersistentLinkedList<>(PersistentTreap.merge(split[0], split[1]));
        return ans;
    }

    /**
     * Get the interval include the l-th and r-th element
     */
    public PersistentLinkedList<T> interval(int l, int r) {
        int size = size();
        if (l == 1 && r == size) {
            return this;
        }
        if (l > r) {
            return new PersistentLinkedList<>();
        }
        if (l == 1) {
            PersistentTreap[] split = PersistentTreap.splitByRank(data, r);
            return new PersistentLinkedList<>(split[0]);
        }
        if (r == size) {
            PersistentTreap[] split = PersistentTreap.splitByRank(data, l - 1);
            return new PersistentLinkedList<>(split[1]);
        }
        PersistentTreap[] split = PersistentTreap.splitByRank(data, r);
        return new PersistentLinkedList<>(PersistentTreap.splitByRank(split[0], l - 1)[1]);
    }

    /**
     * Split this array into 1,2,3,...,mid and mid+1,...,size
     */
    public PersistentLinkedList<T>[] split(int mid) {
        PersistentTreap[] split = PersistentTreap.splitByRank(data, mid);
        return new PersistentLinkedList[]{new PersistentLinkedList(split[0]), new PersistentLinkedList(split[1])};
    }

    public PersistentLinkedList<T> merge(PersistentLinkedList<T> other) {
        return new PersistentLinkedList<>(PersistentTreap.merge(data, other.data));
    }

    public PersistentLinkedList<T> modify(int i, T val) {
        PersistentTreap[] split = PersistentTreap.splitByRank(data, i - 1);
        PersistentTreap[] split2 = PersistentTreap.splitByRank(split[1], 1);
        split2[0] = split2[0].clone();
        split2[0].value = val;
        split[1] = PersistentTreap.merge(split2[0], split2[1]);
        return new PersistentLinkedList<>(PersistentTreap.merge(split[0], split[1]));
    }

    @Override
    public String toString() {
        StringBuilder builder = new StringBuilder("[");
        for (int i = 1; i <= size(); i++) {
            builder.append(get(i)).append(',');
        }
        if (builder.charAt(builder.length() - 1) == ',') {
            builder.setLength(builder.length() - 1);
        }
        builder.append("]");
        return builder.toString();
    }

    private static class PersistentTreap implements Cloneable {
        private static Random random = new Random();
        static PersistentTreap NIL = new PersistentTreap();

        static {
            NIL.left = NIL.right = NIL;
            NIL.size = 0;
        }

        PersistentTreap left = NIL;
        PersistentTreap right = NIL;
        int size = 1;
        Object value;
        //int height;

        @Override
        public PersistentTreap clone() {
            if (this == NIL) {
                return this;
            }
            try {
                return (PersistentTreap) super.clone();
            } catch (CloneNotSupportedException e) {
                throw new RuntimeException(e);
            }
        }

        public void pushUp() {
            if (this == NIL) {
                return;
            }
            size = left.size + right.size + 1;
            //height = Math.max(left.height, right.height) + 1;
        }

        /**
         * split by rank and the node whose rank is argument will stored at result[0]
         */
        public static PersistentTreap[] splitByRank(PersistentTreap root, int rank) {
            if (root == NIL) {
                return new PersistentTreap[]{NIL, NIL};
            }
            root = root.clone();
            PersistentTreap[] result;
            if (root.left.size >= rank) {
                result = splitByRank(root.left, rank);
                root.left = result[1];
                result[1] = root;
            } else {
                result = splitByRank(root.right, rank - (root.size - root.right.size));
                root.right = result[0];
                result[0] = root;
            }
            root.pushUp();
            return result;
        }

        public static PersistentTreap merge(PersistentTreap a, PersistentTreap b) {
            if (a == NIL) {
                return b;
            }
            if (b == NIL) {
                return a;
            }
            if (random.nextInt(a.size + b.size) < a.size) {
                a = a.clone();
                a.right = merge(a.right, b);
                a.pushUp();
                return a;
            } else {
                b = b.clone();
                b.left = merge(a, b.left);
                b.pushUp();
                return b;
            }
        }

        public static PersistentTreap clone(PersistentTreap root) {
            if (root == NIL) {
                return NIL;
            }
            PersistentTreap clone = root.clone();
            clone.left = clone(root.left);
            clone.right = clone(root.right);
            return clone;
        }

        public static Object getValueByRank(PersistentTreap root, int k) {
            while (root.size > 1) {
                if (root.left.size >= k) {
                    root = root.left;
                } else {
                    k -= root.left.size;
                    if (k == 1) {
                        break;
                    }
                    k--;
                    root = root.right;
                }
            }
            return root.value;
        }
    }
}
