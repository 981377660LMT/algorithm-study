template<int N> class LinearSequenceDisjointSet {
private:
    static const int L = log2(N); 
    static const int MASK = ( 1 << L ) - 1; 
    
    int pre[N / L + 1];
    int dat[N / L + 1];
    
    int findPre( int x ) 
      { return x == pre[x] ? x : findPre( pre[x] ); }
    
public:
    void init() {
        for( int i = 0; i <= N / L; i ++ ) {
            pre[i] = i;
            dat[i] = MASK;
        }
    }
 
    int find( int x ) {
        int b = x / L;
        int p = x % L;
        
        int m = dat[b] & ( (1 << p) - 1 );
        
        if( !m ) {
            b = findPre(b);
            m = dat[b];
        }
 
        return b * L + log2(m);
    }
 
    void join( int x ) {
        int b = x / L;
        int p = x % L;
        
        dat[b] &= ( MASK - (1 << p) );
        
        if( p == 0 and b != 0 ) 
          { pre[ findPre(b) ] = findPre(b - 1); }
    }
};