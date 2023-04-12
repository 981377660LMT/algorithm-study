
template<typename T,unsigned lg=20>class PersistentQueue{
    struct Node{T val; Node* cs[lg];};
    Node *fnode=nullptr,*bnode=nullptr; int sz=0;
public:
    PersistentQueue(){}
    PersistentQueue(Node* front,Node* back,int sz):fnode(front),bnode(back),sz(sz){}
    PersistentQueue push(const T& x){
        Node *t=new Node();
        t->val=x;
        t->cs[0]=bnode;
        rep(i,1,lg){
            Node* c=t->cs[i-1];
            if(c)t->cs[i]=c->cs[i-1];
            else break;
        }
        return PersistentQueue(fnode?fnode:t,t,sz+1);
    }
    PersistentQueue pop(){
        if(!fnode||!bnode||sz==1)return PersistentQueue();
        int d=sz-2;
        Node *t=bnode;
        while(d){
            int k=31-__builtin_clz(d);
            d-=1<<k;
            t=t->cs[k];
        }
        return PersistentQueue(t,bnode,sz-1);
    }
    T front()const{return fnode->val;}
    T back()const{return bnode->val;}
};
