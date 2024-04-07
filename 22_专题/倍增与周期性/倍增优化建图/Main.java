import java.io.OutputStream;
import java.io.IOException;
import java.io.InputStream;
import java.util.Arrays;
import java.util.HashMap;
import java.util.Deque;
import java.util.ArrayList;
import java.util.Map;
import java.io.OutputStreamWriter;
import java.io.OutputStream;
import java.util.Collection;
import java.io.IOException;
import java.io.UncheckedIOException;
import java.util.List;
import java.io.Closeable;
import java.io.Writer;
import java.util.ArrayDeque;
import java.io.InputStream;

/**
 * Built using CHelper plug-in Actual solution is at the top
 */
public class Main {
  public static void main(String[] args) throws Exception {
    Thread thread = new Thread(null, new TaskAdapter(), "", 1 << 27);
    thread.start();
    thread.join();
  }

  static class TaskAdapter implements Runnable {
    @Override
    public void run() {
      InputStream inputStream = System.in;
      OutputStream outputStream = System.out;
      FastInput in = new FastInput(inputStream);
      FastOutput out = new FastOutput(outputStream);
      EALT solver = new EALT();
      solver.solve(1, in, out);
      out.close();
    }
  }

  static class EALT {
    int n;
    int m;
    LongMinimumCloseSubGraph subGraph;
    int idAlloc;

    public void solve(int testNumber, FastInput in, FastOutput out) {
      n = in.readInt();
      m = in.readInt();
      idAlloc = m + n - 1;
      Node[] nodes = new Node[n];
      for (int i = 0; i < n; i++) {
        nodes[i] = new Node();
        nodes[i].id = i + 1;
      }

      MultiWayIntegerStack edges = new MultiWayIntegerStack(n, n * 2);
      for (int i = 1; i < n; i++) {
        int aId = in.readInt() - 1;
        int bId = in.readInt() - 1;
        Node a = nodes[aId];
        Node b = nodes[bId];
        Edge e = new Edge();
        e.a = a;
        e.b = b;
        e.id = i - 1;
        a.next.add(e);
        b.next.add(e);

        edges.addLast(aId, bId);
        edges.addLast(bId, aId);
      }

      LcaOnTree lcaOnTree = new LcaOnTree(edges, 0);

      int vertexNum = m + n - 1 + (n - 1) * 15;
      long[] weights = new long[vertexNum];
      for (int i = 0; i < m; i++) {
        weights[idOfCitizen(i)] = 1;
      }
      for (int i = 0; i < n - 1; i++) {
        weights[idOfEdge(i)] = -1;
      }
      subGraph = new LongMinimumCloseSubGraph(weights);
      dfs(nodes[0], null, null, 0);

      for (int i = 0; i < m; i++) {
        int a = in.readInt() - 1;
        int b = in.readInt() - 1;
        int lca = lcaOnTree.lca(a, b);
        findLCA(nodes[a], nodes[lca].depth, i);
        findLCA(nodes[b], nodes[lca].depth, i);
      }

      subGraph.solve();
      IntegerList citizen = new IntegerList(m);
      IntegerList guard = new IntegerList(n);
      boolean[] status = subGraph.getStatus();
      for (int i = 0; i < m; i++) {
        if (!status[i]) {
          citizen.add(i);
        }
      }
      for (int i = 0; i < n - 1; i++) {
        if (status[m + i]) {
          guard.add(i);
        }
      }

      out.println(citizen.size() + guard.size());
      out.append(citizen.size()).append(' ');
      for (int i = 0; i < citizen.size(); i++) {
        out.append(citizen.get(i) + 1).append(' ');
      }
      out.println();
      out.append(guard.size()).append(' ');
      for (int i = 0; i < guard.size(); i++) {
        out.append(guard.get(i) + 1).append(' ');
      }
    }

    public int idOfCitizen(int i) {
      return i;
    }

    public int idOfEdge(int i) {
      return m + i;
    }

    public int nextId() {
      return ++idAlloc;
    }

    public void dfs(Node root, Node p, Edge from, int depth) {
      root.depth = depth;
      if (from != null) {
        root.jump[0] = p;
        root.dependencyIds[0] = idOfEdge(from.id);
        for (int i = 0; root.jump[i] != null; i++) {
          root.jump[i + 1] = root.jump[i].jump[i];
          if (root.jump[i + 1] != null) {
            root.dependencyIds[i + 1] = nextId();
            subGraph.addDependency(root.dependencyIds[i + 1], root.dependencyIds[i]);
            subGraph.addDependency(root.dependencyIds[i + 1], root.jump[i].dependencyIds[i]);
          }
        }
      }
      for (Edge e : root.next) {
        if (e == from) {
          continue;
        }
        Node node = e.other(root);
        dfs(node, root, e, depth + 1);
      }
    }

    public void findLCA(Node a, int targetDepth, int citizen) {
      if (a.depth == targetDepth) {
        return;
      }
      int differ = a.depth - targetDepth;
      int log = CachedLog2.floorLog(differ);
      subGraph.addDependency(idOfCitizen(citizen), a.dependencyIds[log]);
      findLCA(a.jump[log], targetDepth, citizen);
    }

  }

  static class IntegerList implements Cloneable {
    private int size;
    private int cap;
    private int[] data;
    private static final int[] EMPTY = new int[0];

    public IntegerList(int cap) {
      this.cap = cap;
      if (cap == 0) {
        data = EMPTY;
      } else {
        data = new int[cap];
      }
    }

    public IntegerList(IntegerList list) {
      this.size = list.size;
      this.cap = list.cap;
      this.data = Arrays.copyOf(list.data, size);
    }

    public IntegerList() {
      this(0);
    }

    public void ensureSpace(int req) {
      if (req > cap) {
        while (cap < req) {
          cap = Math.max(cap + 10, 2 * cap);
        }
        data = Arrays.copyOf(data, cap);
      }
    }

    private void checkRange(int i) {
      if (i < 0 || i >= size) {
        throw new ArrayIndexOutOfBoundsException("index " + i + " out of range");
      }
    }

    public int get(int i) {
      checkRange(i);
      return data[i];
    }

    public void add(int x) {
      ensureSpace(size + 1);
      data[size++] = x;
    }

    public void addAll(int[] x, int offset, int len) {
      ensureSpace(size + len);
      System.arraycopy(x, offset, data, size, len);
      size += len;
    }

    public void addAll(IntegerList list) {
      addAll(list.data, 0, list.size);
    }

    public int size() {
      return size;
    }

    public int[] toArray() {
      return Arrays.copyOf(data, size);
    }

    public String toString() {
      return Arrays.toString(toArray());
    }

    public IntegerIterator iterator() {
      return new IntegerIterator() {
        int i = 0;


        public boolean hasNext() {
          return i < size;
        }


        public int next() {
          return data[i++];
        }
      };
    }

    public boolean equals(Object obj) {
      if (!(obj instanceof IntegerList)) {
        return false;
      }
      IntegerList other = (IntegerList) obj;
      return SequenceUtils.equal(data, 0, size - 1, other.data, 0, other.size - 1);
    }

    public int hashCode() {
      int h = 1;
      for (int i = 0; i < size; i++) {
        h = h * 31 + Integer.hashCode(data[i]);
      }
      return h;
    }

    public IntegerList clone() {
      IntegerList ans = new IntegerList();
      ans.addAll(this);
      return ans;
    }

  }

  static class SequenceUtils {
    public static boolean equal(int[] a, int al, int ar, int[] b, int bl, int br) {
      if ((ar - al) != (br - bl)) {
        return false;
      }
      for (int i = al, j = bl; i <= ar; i++, j++) {
        if (a[i] != b[j]) {
          return false;
        }
      }
      return true;
    }

  }

  static class Node {
    List<Edge> next = new ArrayList<>();
    int depth;
    Node[] jump = new Node[16];
    int[] dependencyIds = new int[16];
    int id;

    public String toString() {
      return "" + id;
    }

  }

  static class LongMinimumCloseSubGraph {
    private LongISAP isap;
    private int n;
    private long sumOfPositive;
    private boolean solved;
    private long minCut;
    private boolean[] status;

    public LongMinimumCloseSubGraph(long[] weight) {
      this.n = weight.length;
      isap = new LongISAP(n + 2);
      int idOfSrc = n;
      int idOfDst = n + 1;
      isap.setSource(idOfSrc);
      isap.setTarget(idOfDst);

      for (int i = 0; i < n; i++) {
        if (weight[i] > 0) {
          isap.getChannel(idOfSrc, i).reset(weight[i], 0);
          sumOfPositive += weight[i];
        }
        if (weight[i] < 0) {
          isap.getChannel(i, idOfDst).reset(-weight[i], 0);
        }
      }
    }

    public void addDependency(int from, int to) {
      if (!(from >= 0 && from < n && to >= 0 && to < n)) {
        throw new IllegalArgumentException();
      }
      isap.getChannel(from, to).reset((long) 1e18, 0);
    }

    public long solve() {
      if (!solved) {
        minCut = isap.send((long) 1e18);
        solved = true;
      }
      return sumOfPositive - minCut;
    }

    public boolean[] getStatus() {
      if (status == null) {
        solve();
        status = new boolean[n];
        for (IntegerIterator iterator = isap.getComponentS().iterator(); iterator.hasNext();) {
          int node = iterator.next();
          if (node < n) {
            status[node] = true;
          }
        }
      }
      return status;
    }

  }

  static interface IntegerIterator {
    boolean hasNext();

    int next();

  }

  static class MultiWayIntegerStack {
    private int[] values;
    private int[] next;
    private int[] heads;
    private int alloc;
    private int stackNum;

    public IntegerIterator iterator(final int queue) {
      return new IntegerIterator() {
        int ele = heads[queue];


        public boolean hasNext() {
          return ele != 0;
        }


        public int next() {
          int ans = values[ele];
          ele = next[ele];
          return ans;
        }
      };
    }

    private void doubleCapacity() {
      int newSize = Math.max(next.length + 10, next.length * 2);
      next = Arrays.copyOf(next, newSize);
      values = Arrays.copyOf(values, newSize);
    }

    public void alloc() {
      alloc++;
      if (alloc >= next.length) {
        doubleCapacity();
      }
      next[alloc] = 0;
    }

    public int stackNumber() {
      return stackNum;
    }

    public MultiWayIntegerStack(int qNum, int totalCapacity) {
      values = new int[totalCapacity + 1];
      next = new int[totalCapacity + 1];
      heads = new int[qNum];
      stackNum = qNum;
    }

    public void addLast(int qId, int x) {
      alloc();
      values[alloc] = x;
      next[alloc] = heads[qId];
      heads[qId] = alloc;
    }

    public String toString() {
      StringBuilder builder = new StringBuilder();
      for (int i = 0; i < stackNum; i++) {
        builder.append(i).append(": ");
        for (IntegerIterator iterator = iterator(i); iterator.hasNext();) {
          builder.append(iterator.next()).append(",");
        }
        if (builder.charAt(builder.length() - 1) == ',') {
          builder.setLength(builder.length() - 1);
        }
        builder.append('\n');
      }
      return builder.toString();
    }

  }

  static class Edge {
    Node a;
    Node b;
    int id;

    Node other(Node x) {
      return x == a ? b : a;
    }

  }

  static class CachedLog2 {
    private static final int BITS = 16;
    private static final int LIMIT = 1 << BITS;
    private static final byte[] CACHE = new byte[LIMIT];

    static {
      int b = 0;
      for (int i = 0; i < LIMIT; i++) {
        while ((1 << (b + 1)) <= i) {
          b++;
        }
        CACHE[i] = (byte) b;
      }
    }

    public static int floorLog(int x) {
      return x < LIMIT ? CACHE[x] : (BITS + CACHE[x >>> BITS]);
    }

  }

  static class LongISAP {
    LongISAP.Node[] nodes;
    int[] distanceCnt;
    LongISAP.Node source;
    LongISAP.Node target;
    int nodeNum;
    Map<Long, LongISAP.DirectLongChannel> channelMap;
    Deque<LongISAP.Node> deque;
    long totalFlow;

    public IntegerList getComponentS() {
      IntegerList result = new IntegerList(nodeNum);
      for (int i = 0; i < nodeNum; i++) {
        nodes[i].visited = false;
      }
      deque.addLast(source);
      source.visited = true;
      while (!deque.isEmpty()) {
        LongISAP.Node head = deque.removeFirst();
        result.add(head.id);
        for (LongISAP.LongChannel channel : head.channelList) {
          if (channel.getFlow() == channel.getCapacity()) {
            continue;
          }
          LongISAP.Node node = channel.getDst();
          if (node.visited) {
            continue;
          }
          node.visited = true;
          deque.addLast(node);
        }
      }
      return result;
    }

    private Collection<LongISAP.DirectLongChannel> getChannels() {
      return channelMap.values();
    }

    private LongISAP.DirectLongChannel addChannel(int src, int dst) {
      LongISAP.DirectLongChannel channel =
          new LongISAP.DirectLongChannel(nodes[src], nodes[dst], 0, 0);
      nodes[src].channelList.add(channel);
      nodes[dst].channelList.add(channel.getInverse());
      return channel;
    }

    public LongISAP.DirectLongChannel getChannel(int src, int dst) {
      Long id = (((long) src) << 32) | dst;
      LongISAP.DirectLongChannel channel = channelMap.get(id);
      if (channel == null) {
        channel = addChannel(src, dst);
        channelMap.put(id, channel);
      }
      return channel;
    }

    public LongISAP(int nodeNum) {
      channelMap = new HashMap<>(nodeNum);
      this.nodeNum = nodeNum;
      deque = new ArrayDeque(nodeNum);
      nodes = new LongISAP.Node[nodeNum];
      distanceCnt = new int[nodeNum + 2];
      for (int i = 0; i < nodeNum; i++) {
        LongISAP.Node node = new LongISAP.Node();
        node.id = i;
        nodes[i] = node;
      }
    }

    public long send(long flow) {
      long sent = 0;
      bfs();
      while (flow > sent && source.distance < nodeNum) {
        sent += send(source, flow - sent);
      }
      totalFlow += sent;
      return sent;
    }

    private long send(LongISAP.Node node, long flowRemain) {
      if (node == target) {
        return flowRemain;
      }

      long sent = 0;
      int nextDistance = node.distance - 1;
      for (LongISAP.LongChannel channel : node.channelList) {
        long channelRemain = channel.getCapacity() - channel.getFlow();
        LongISAP.Node dst = channel.getDst();
        if (channelRemain == 0 || dst.distance != nextDistance) {
          continue;
        }
        long actuallySend = send(channel.getDst(), Math.min(flowRemain - sent, channelRemain));
        channel.sendFlow(actuallySend);
        sent += actuallySend;
        if (flowRemain == sent) {
          break;
        }
      }

      if (sent == 0) {
        if (--distanceCnt[node.distance] == 0) {
          distanceCnt[source.distance]--;
          source.distance = nodeNum;
          distanceCnt[source.distance]++;
          if (node != source) {
            distanceCnt[++node.distance]++;
          }
        } else {
          distanceCnt[++node.distance]++;
        }
      }

      return sent;
    }

    public void setSource(int id) {
      source = nodes[id];
    }

    public void setTarget(int id) {
      target = nodes[id];
    }

    private void bfs() {
      Arrays.fill(distanceCnt, 0);
      deque.clear();

      for (int i = 0; i < nodeNum; i++) {
        nodes[i].distance = nodeNum;
      }

      target.distance = 0;
      deque.addLast(target);

      while (!deque.isEmpty()) {
        LongISAP.Node head = deque.removeFirst();
        distanceCnt[head.distance]++;
        for (LongISAP.LongChannel channel : head.channelList) {
          LongISAP.LongChannel inverse = channel.getInverse();
          if (inverse.getCapacity() == inverse.getFlow()) {
            continue;
          }
          LongISAP.Node dst = channel.getDst();
          if (dst.distance != nodeNum) {
            continue;
          }
          dst.distance = head.distance + 1;
          deque.addLast(dst);
        }
      }
    }

    public String toString() {
      StringBuilder builder = new StringBuilder();
      for (LongISAP.DirectLongChannel channel : getChannels()) {
        if (channel.getFlow() == 0) {
          continue;
        }
        builder.append(channel).append('\n');
      }

      for (LongISAP.DirectLongChannel channel : getChannels()) {
        if (channel.getFlow() != 0) {
          continue;
        }
        builder.append(channel).append('\n');
      }
      return builder.toString();
    }

    public static interface LongChannel {
      public LongISAP.Node getSrc();

      public LongISAP.Node getDst();

      public long getCapacity();

      public long getFlow();

      public void sendFlow(long volume);

      public LongISAP.LongChannel getInverse();

    }

    public static class DirectLongChannel implements LongISAP.LongChannel {
      final LongISAP.Node src;
      final LongISAP.Node dst;
      final int id;
      long capacity;
      long flow;
      LongISAP.LongChannel inverse;

      public DirectLongChannel(LongISAP.Node src, LongISAP.Node dst, int capacity, int id) {
        this.src = src;
        this.dst = dst;
        this.capacity = capacity;
        this.id = id;
        inverse = new LongISAP.InverseLongChannelWrapper(this);
      }

      public void reset(long cap, long flow) {
        this.flow = flow;
        this.capacity = cap;
      }

      public String toString() {
        return String.format("%s--%s/%s-->%s", getSrc(), getFlow(), getCapacity(), getDst());
      }

      public LongISAP.Node getSrc() {
        return src;
      }

      public LongISAP.LongChannel getInverse() {
        return inverse;
      }

      public LongISAP.Node getDst() {
        return dst;
      }

      public long getCapacity() {
        return capacity;
      }

      public long getFlow() {
        return flow;
      }

      public void sendFlow(long volume) {
        flow += volume;
      }

    }

    public static class InverseLongChannelWrapper implements LongISAP.LongChannel {
      final LongISAP.LongChannel channel;

      public InverseLongChannelWrapper(LongISAP.LongChannel channel) {
        this.channel = channel;
      }

      public LongISAP.LongChannel getInverse() {
        return channel;
      }

      public LongISAP.Node getSrc() {
        return channel.getDst();
      }

      public LongISAP.Node getDst() {
        return channel.getSrc();
      }

      public long getCapacity() {
        return channel.getFlow();
      }

      public long getFlow() {
        return 0;
      }

      public void sendFlow(long volume) {
        channel.sendFlow(-volume);
      }

      public String toString() {
        return String.format("%s--%s/%s-->%s", getSrc(), getFlow(), getCapacity(), getDst());
      }

    }

    private static class Node {
      int id;
      int distance;
      boolean visited;
      List<LongISAP.LongChannel> channelList = new ArrayList(1);

      public String toString() {
        return "" + id;
      }

    }

  }

  static class LcaOnTree {
    int[] parent;
    int[] preOrder;
    int[] i;
    int[] head;
    int[] a;
    int time;

    void dfs1(MultiWayIntegerStack tree, int u, int p) {
      parent[u] = p;
      i[u] = preOrder[u] = time++;
      for (IntegerIterator iterator = tree.iterator(u); iterator.hasNext();) {
        int v = iterator.next();
        if (v == p)
          continue;
        dfs1(tree, v, u);
        if (Integer.lowestOneBit(i[u]) < Integer.lowestOneBit(i[v])) {
          i[u] = i[v];
        }
      }
      head[i[u]] = u;
    }

    void dfs2(MultiWayIntegerStack tree, int u, int p, int up) {
      a[u] = up | Integer.lowestOneBit(i[u]);
      for (IntegerIterator iterator = tree.iterator(u); iterator.hasNext();) {
        int v = iterator.next();
        if (v == p)
          continue;
        dfs2(tree, v, u, a[u]);
      }
    }

    public LcaOnTree(MultiWayIntegerStack tree, int root) {
      int n = tree.stackNumber();
      preOrder = new int[n];
      i = new int[n];
      head = new int[n];
      a = new int[n];
      parent = new int[n];

      dfs1(tree, root, -1);
      dfs2(tree, root, -1, 0);
    }

    private int enterIntoStrip(int x, int hz) {
      if (Integer.lowestOneBit(i[x]) == hz)
        return x;
      int hw = 1 << CachedLog2.floorLog(a[x] & (hz - 1));
      return parent[head[i[x] & -hw | hw]];
    }

    public int lca(int x, int y) {
      int hb = i[x] == i[y] ? Integer.lowestOneBit(i[x]) : (1 << CachedLog2.floorLog(i[x] ^ i[y]));
      int hz = Integer.lowestOneBit(a[x] & a[y] & -hb);
      int ex = enterIntoStrip(x, hz);
      int ey = enterIntoStrip(y, hz);
      return preOrder[ex] < preOrder[ey] ? ex : ey;
    }

  }

  static class FastInput {
    private final InputStream is;
    private byte[] buf = new byte[1 << 13];
    private int bufLen;
    private int bufOffset;
    private int next;

    public FastInput(InputStream is) {
      this.is = is;
    }

    private int read() {
      while (bufLen == bufOffset) {
        bufOffset = 0;
        try {
          bufLen = is.read(buf);
        } catch (IOException e) {
          bufLen = -1;
        }
        if (bufLen == -1) {
          return -1;
        }
      }
      return buf[bufOffset++];
    }

    public void skipBlank() {
      while (next >= 0 && next <= 32) {
        next = read();
      }
    }

    public int readInt() {
      int sign = 1;

      skipBlank();
      if (next == '+' || next == '-') {
        sign = next == '+' ? 1 : -1;
        next = read();
      }

      int val = 0;
      if (sign == 1) {
        while (next >= '0' && next <= '9') {
          val = val * 10 + next - '0';
          next = read();
        }
      } else {
        while (next >= '0' && next <= '9') {
          val = val * 10 - next + '0';
          next = read();
        }
      }

      return val;
    }

  }

  static class FastOutput implements AutoCloseable, Closeable {
    private StringBuilder cache = new StringBuilder(10 << 20);
    private final Writer os;

    public FastOutput(Writer os) {
      this.os = os;
    }

    public FastOutput(OutputStream os) {
      this(new OutputStreamWriter(os));
    }

    public FastOutput append(char c) {
      cache.append(c);
      return this;
    }

    public FastOutput append(int c) {
      cache.append(c);
      return this;
    }

    public FastOutput println(int c) {
      cache.append(c).append('\n');
      return this;
    }

    public FastOutput println() {
      cache.append('\n');
      return this;
    }

    public FastOutput flush() {
      try {
        os.append(cache);
        os.flush();
        cache.setLength(0);
      } catch (IOException e) {
        throw new UncheckedIOException(e);
      }
      return this;
    }

    public void close() {
      flush();
      try {
        os.close();
      } catch (IOException e) {
        throw new UncheckedIOException(e);
      }
    }

    public String toString() {
      return cache.toString();
    }

  }
}
