package template.datastructure;

import template.binary.Bits;
import template.utils.Buffer;
import template.utils.Pair;

/**
 * @deprecated
 */
public class BinaryTree implements Cloneable {
    public BinaryTree left;
    public BinaryTree right;
    public int size;
    static Buffer<BinaryTree> buf = new Buffer<>(BinaryTree::new, t -> {
        t.left = t.right = null;
        t.size = 0;
    });

    private BinaryTree get(int i) {
        if (i == 0) {
            if (left == null) {
                left = buf.alloc();
            }
            return left;
        } else {
            if (right == null) {
                right = buf.alloc();
            }
            return right;
        }
    }


    public void destroy() {
        if (left != null) {
            left.destroy();
        }
        if (right != null) {
            right.destroy();
        }
        buf.release(this);
    }

    public int size(int i) {
        if (i == 0) {
            return left == null ? 0 : left.size;
        }
        return right == null ? 0 : right.size;
    }

    public void add(int x, int height, int mod) {
        if (height < 0) {
            size += mod;
            return;
        }
        pushDown(height);
        get(Bits.get(x, height)).add(x, height - 1, mod);
        pushUp();
    }

    public int find(int x, int height) {
        if (height < 0) {
            return size;
        }
        pushDown(height);
        return get(Bits.get(x, height)).find(x, height - 1);
    }

    public int kthElement(int k, int height) {
        if (height < 0) {
            return 0;
        }
        pushDown(height);
        if (size(0) >= k) {
            return get(0).kthElement(k, height - 1);
        }
        return get(1).kthElement(k - size(0), height - 1);
    }

    public int prefixSum(int x, int height) {
        if (height < 0) {
            return size;
        }
        pushDown(height);
        int ans = get(Bits.get(x, height)).prefixSum(x, height - 1);
        if (Bits.get(x, height) == 1) {
            ans += size(0);
        }
        return ans;
    }

    public int interval(int l, int r, int h) {
        int ans = prefixSum(r, h);
        if (l > 0) {
            ans -= prefixSum(l - 1, h);
        }
        return ans;
    }


    public void pushUp() {
        size = 0;
        if (left != null) {
            size += left.size;
        }
        if (right != null) {
            size += right.size;
        }
    }

    public void pushDown(int height) {

    }

    public int maxXor(int x, int height) {
        if (height < 0) {
            return 0;
        }
        pushDown(height);
        int prefer = Bits.get(x, height) ^ 1;
        int ans;
        if (size(prefer) > 0) {
            ans = get(prefer).maxXor(x, height - 1);
            ans |= prefer << height;
        } else {
            ans = get(1 ^ prefer).maxXor(x, height - 1);
            ans |= (1 ^ prefer) << height;
        }
        return ans;
    }

    public int minXor(int x, int height) {
        if (height < 0) {
            return 0;
        }
        pushDown(height);
        int prefer = Bits.get(x, height);
        int ans;
        if (size(prefer) > 0) {
            ans = get(prefer).minXor(x, height - 1);
            ans |= prefer << height;
        } else {
            ans = get(1 ^ prefer).minXor(x, height - 1);
            ans |= (1 ^ prefer) << height;
        }
        return ans;
    }

    /**
     * res.a <= key and res.b > key
     *
     * @param key
     * @return
     */
    public static Pair<BinaryTree, BinaryTree> split(BinaryTree bt, int key, int height, boolean toLeft) {
        if (bt == null || height == -1) {
            if (toLeft) {
                return new Pair<>(bt, null);
            } else {
                return new Pair<>(null, bt);
            }
        }
        Pair<BinaryTree, BinaryTree> ans;
        bt.pushDown(height);
        if (Bits.get(key, height) == 0) {
            ans = split(bt.left, key, height - 1, toLeft);
            bt.left = ans.b;
            BinaryTree a = new BinaryTree();
            a.left = ans.a;
            ans.a = a;
            ans.b = bt;
        } else {
            ans = split(bt.right, key, height - 1, toLeft);
            bt.right = ans.a;
            BinaryTree b = new BinaryTree();
            b.right = ans.b;
            ans.a = bt;
            ans.b = b;
        }
        ans.a.pushUp();
        ans.b.pushUp();
        return ans;
    }

    public static BinaryTree merge(BinaryTree a, BinaryTree b, int height) {
        if (a == null) {
            return b;
        }
        if (b == null) {
            return a;
        }
        if (height == -1) {
            a.size += b.size;
            return a;
        }
        a.pushDown(height);
        b.pushDown(height);
        a.left = merge(a.left, b.left, height - 1);
        a.right = merge(a.right, b.right, height - 1);
        a.pushUp();
        return a;
    }

    @Override
    public BinaryTree clone() {
        try {
            return (BinaryTree) super.clone();
        } catch (CloneNotSupportedException e) {
            throw new RuntimeException(e);
        }
    }
}
