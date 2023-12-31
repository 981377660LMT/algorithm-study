class PersistentStack {
    PersistentStack next = null;
    int val;

    public PersistentStack(int val, PersistentStack next) {
        this.val = val;
        this.next = next;
    }
    
    public PersistentStack() {
    }
    
    int top() {
        return val;
    }
    
    PersistentStack pop() {
        return next != null ? (next = next.eval()) : next;
    }
    
    static PersistentStack push(PersistentStack x, int v) {
        return new PersistentStack(v, x);
    }
    
    static PersistentStack concat(PersistentStack x, PersistentStack y) {
        if (x == null) return y.eval();
        else return new PersistentStack(x.val, new PersistentStack() {
                @Override
                PersistentStack eval() {
                    return PersistentStack.concat(x.pop(), y);
                }
            });
    }
    
    static PersistentStack reverse(PersistentStack head) {
        return new PersistentStack () {
            @Override
            PersistentStack eval() {
                PersistentStack ret = null;
                for (PersistentStack x=head;x!=null;x=x.pop()) {
                    ret=PersistentStack.push(ret, x.top());
                }
                return ret;
            }
        };
    }
    
    synchronized PersistentStack eval() {
        return this;
    }
}

class PersistentQueue {
    int fsize = 0;
    int rsize = 0;
    PersistentStack f = null;
    PersistentStack r = null;
    
    public PersistentQueue(PersistentStack f, int fsize, PersistentStack r, int rsize) {
        this.fsize = fsize;
        this.rsize = rsize;
        this.f = f;
        this.r = r;
    }
    
    public PersistentQueue() {
    }
    
    boolean isEmpty() {
        return fsize == 0;
    }
    
    int top() {
        return f.top();
    }
    
    PersistentQueue push(int x) {
        return new PersistentQueue(f, fsize, PersistentStack.push(r, x), rsize+1).normalize();
    }
    
    PersistentQueue pop() {
        return new PersistentQueue(f.pop(), fsize-1, r, rsize).normalize();
    }
    
    PersistentQueue normalize() {
        if (fsize >= rsize) return this;
        else {
            return new PersistentQueue(PersistentStack.concat(f, PersistentStack.reverse(r)), fsize+rsize, null, 0);
        }
    }
}