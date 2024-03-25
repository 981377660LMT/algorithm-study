const int		inf	= 1 << 30;
static const int	maxn	= 210000;       /* maxn 应当大于所有集合大小的总和 */
struct node {                                   /* 该算法实现的整数集合为可重集合（multiset） */
	int l, r, v;                            /* 左节点编号 右节点编号 值 */
}	T[maxn * 25];
int	len, sz, info[maxn];                    /* len 是值域线段树的长度 */
int ranking( int v )
{
/* 如果不需要数据离散化直接 return v 即可，同时还要将所有的 info[v] 改为 v，此时线段树的值域为 , → [1, len] 。*/
	return(lower_bound( info + 1, info + len + 1, v ) - info);
}


void ins( int & i, int l, int r, int p )   /* 在区间 [l, r] 中插入 p */
{
	if ( !i )
	{
		i	= ++sz;
		T[i].v	= T[i].l = T[i].r = 0;
	}
	int m = (l + r) / 2;
	T[i].v++; /* 递增当前区间内值的个数 */
	if ( l == r )
		return;
	if ( p <= m )
		ins( T[i].l, l, m, p );
	else ins( T[i].r, m + 1, r, p );
}


void insert( int & x, int v )                                   /* 向集合 x 中添加数值 v */
{
	ins( x, 1, len, ranking( v ) );
}


void init( int* A, int length )                                 /* 输入数组 A 的下标从 1 开始，length 表示 A 的长度 */
{
	sz = 0;                                                 /* A 记录集合中所有可能出现的整数 */
	copy( A + 1, A + length + 1, info + 1 );
	sort( info + 1, info + length + 1 );                    /* 如果没有离散化需要对 len 进行赋值来指明线段树的值域。 */
	len = unique( info + 1, info + length + 1 ) - info - 1; /* 对 info 数组进行排序去重 */
}


int kth( int x, int k )                                         /* 查找集合 x 中的第 k 大元素 */
{
	int l = 1, r = len;
	while ( l < r )
	{
		int m = (l + r) / 2, t = T[T[x].l].v;

		if ( k <= t )
			x = T[x].l, r = m;
		else
			x = T[x].r, l = m + 1, k -= t;
	}
	if ( k > T[x].v )       /* 如果 k 大于集合 x 的大小，则返回 inf */
		return(inf);
	return(info[r]);
}


int ask( int x, int v )         /* 返回集合 x 中比 v 小的数的个数 */
{
	int	l	= 1, r = len, k = 0;
	int	p	= ranking( v ) - 1;
	if ( p <= 0 )
		return(0);
	while ( l < r )
	{
		int m = (l + r) / 2, t = T[T[x].l].v;
		if ( p <= m )
			x = T[x].l, r = m;
		else
			x = T[x].r, l = m + 1, k += t;
	}
	k += T[x].v;
	return(k);
}


int pre( int x, int v )   /* 在集合 x 中查询 v 的前驱，找不到返回-inf */
{
	int k = ask( x, v );
	if ( k == 0 )
		return(-inf);
	return(kth( x, k ) );
}


int next( int x, int v )        /* 在集合 x 中查询 v 的后继，找不到返回 inf */
{
	int k = ask( x, v + 1 ) + 1;
	return(kth( x, k ) );
}


int count( int x, int v )       /* 返回集合 x 中数值 v 的个数 */
{
	int	l	= 1, r = len;
	int	p	= ranking( v );
	if ( p > len || info[p] != v )
		return(0);
	while ( l < r && T[x].v > 0 )
	{
		int m = (l + r) / 2;
		if ( p <= m )
			x = T[x].l, r = m;
		else
			x = T[x].r, l = m + 1;
	}
	return(T[x].v);
}


bool find( int x, int v )       /* 返回集合 x 中是否存在 v */
{
	return(count( x, v ) >= 1);
}


int maximum( int x )            /* 返回集合 x 中的最大值 */
{
	int l = 1, r = len;
	if ( T[x].v == 0 )      /* 如果集合为空则返回-inf */
		return(-inf);
	while ( l < r )
	{
		int m = (l + r) / 2, t = T[T[x].r].v;
		if ( t == 0 )
			x = T[x].l, r = m;
		else
			x = T[x].r, l = m + 1;
	}
	return(info[r]);
}


int minimum( int x )            /* 返回集合 x 中的最小值 */
{
	int l = 1, r = len;
	if ( T[x].v == 0 )      /* 如果集合为空则返回 inf */
		return(inf);
	while ( l < r )
	{
		int m = (l + r) / 2, t = T[T[x].l].v;
		if ( t > 0 )
			x = T[x].l, r = m;
		else
			x = T[x].r, l = m + 1;
	}
	return(info[r]);
}


int merge( int x, int y, int l = 1, int r = len )       /* 将集合 x 与集合 y 合并成一个新的集合后返回 */
{
	if ( !x )
		return(y);                              /* 合并之后 x、y 均会失效 */
	if ( !y )
		return(x);
	if ( l == r )
	{
		T[x].v += T[y].v;
	}else  {
		int m = (l + r) / 2;
		T[x].l	= merge( T[x].l, T[y].l, l, m );
		T[x].r	= merge( T[x].r, T[y].r, m + 1, r );
		T[x].v	= T[T[x].l].v + T[T[x].r].v;
	}
	return(x);
}


int diff( int x, int y, int l = 1, int r = len )        /* 返回集合 x 与集合 y 的差集（x - y）, 执行之,→ 后 x、y 均会失效 */
{
	if ( l == r )                                   /* 时间复杂度 size(y)log(len) */
	{
		T[x].v = max( 0, T[x].v - T[y].v );
	}else  {
		int m = (l + r) / 2;
		if ( T[y].l )
			T[x].l = diff( T[x].l, T[y].l, l, m );
		if ( T[y].r )
			T[x].r = diff( T[x].r, T[y].r, m + 1, r );
		T[x].v = T[T[x].l].v + T[T[x].r].v;
	}
	return(x);
}


int intersect( int x, int y, int l = 1, int r = len )   /* 返回集合 x 与集合 y 的交集，执行之后 x、,→ y 均会失效 */
{
	if ( l == r )                                   /* 时间复杂度 min{size(x),size(y)} * log(len) */
	{
		T[x].v = min( T[x].v, T[y].v );
	}else  {
		int m = (l + r) / 2;
		T[x].l	= (!T[x].l || !T[y].l ? 0 : intersect( T[x].l, T[y].l, l, m ) );
		T[x].r	= (!T[x].r || !T[y].r ? 0 : intersect( T[x].r, T[y].r, m + 1, r ) );
		T[x].v	= T[T[x].l].v + T[T[x].r].v;
	}
	return(T[x].v == 0 ? 0 : x);
}


int symmetric( int x, int y, int l = 1, int r = len )   /* 返回集合 x 与集合 y 的对称差分，执行之,→ 后 x、y 均会失效 */
{
	if ( !x )
		return(y);                              /* 时间复杂度 min{size(x),size(y)} *,→ log(len) */
	if ( !y )
		return(x);
	if ( l == r )
	{
		T[x].v = abs( T[x].v - T[y].v );
	}else  {
		int m = (l + r) / 2;
		T[x].l	= symmetric( T[x].l, T[y].l, l, m );
		T[x].r	= symmetric( T[x].r, T[y].r, m + 1, r );
		T[x].v	= T[T[x].l].v + T[T[x].r].v;
	}
	return(x);
}


void left( vector<int> & X, vector<int> & Y )
{
	for ( auto & i : X )
		i = T[i].l;
	for ( auto & i : Y )
		i = T[i].l;
}


void right( vector<int> & X, vector<int> & Y )
{
	for ( auto & i : X )
		i = T[i].r;
	for ( auto & i : Y )
		i = T[i].r;
}


int left_value( vector<int> & X, vector<int> & Y )
{
	int tot = 0;
	for ( auto i : X )
		tot += T[T[i].l].v;
	for ( auto i : Y )
		tot -= T[T[i].l].v;
	return(tot);
}


int right_value( vector<int> & X, vector<int> & Y )
{
	int tot = 0;
	for ( auto i : X )
		tot += T[T[i].r].v;
	for ( auto i : Y )
		tot -= T[T[i].r].v;
	return(tot);
}


int value( vector<int> & X, vector<int> & Y )
{
	int tot = 0;
	for ( auto i : X )
		tot += T[i].v;
	for ( auto i : Y )
		tot -= T[i].v;
	return(tot);
}


int kth( vector<int> X, vector<int> Y, int k )  /* 计算 X 中集合的并集 减去 Y 中集合的并集 所得集,→ 合的第 k 大 */
{
	int l = 1, r = len;                     /* 必须保证 Y 的并集 是 X 的并集 的子集 */
	while ( l < r )
	{
		int m = (l + r) / 2, t = left_value( X, Y );
		if ( k <= t )
			left( X, Y ), r = m;
		else
			right( X, Y ), l = m + 1, k -= t;
	}
	if ( k > value( X, Y ) )
		return(inf);
	return(info[r]);
}


int ask( vector<int> X, vector<int> Y, int v )  /* 计算 X 中集合的并集 减去 Y 中集合的并集 所得集,→ 合中比 v 小的数的个数 */
{
	int	l	= 1, r = len, k = 0;    /* 必须保证 Y 的并集 是 X 的并集 的子集 */
	int	p	= ranking( v ) - 1;
	if ( p <= 0 )
		return(0);
	while ( l < r )
	{
		int m = (l + r) / 2;
		if ( p <= m )
			left( X, Y ), r = m;
		else {
			k += left_value( X, Y );
			right( X, Y );
			l = m + 1;
		}
	}
	k += value( X, Y );
	return(k);
}


int pre( vector<int> X, vector<int> Y, int v )          /* 计算 X 中集合的并集 减去 Y 中集合的并集 所得集,→ 合中 v 的前驱 */
{
	int k = ask( X, Y, v );                         /* 必须保证 Y 的并集 是 X 的并集 的子集 */
	if ( k == 0 )
		return(-inf);                           /* 找不到返回-inf */
	return(kth( X, Y, k ) );
}


int next( vector<int> X, vector<int> Y, int v )         /* 计算 X 中集合的并集 减去 Y 中集合的并集 所得,→ 集合中 v 的后继 */
{
	int k = ask( X, Y, v + 1 ) + 1;                 /* 必须保证 Y 的并集 是 X 的并集 的子集 */
	return(kth( X, Y, k ) );                        /* 找不到返回 inf */
}


int count( vector<int> X, vector<int> Y, int v )        /* 计算 X 中集合的并集 减去 Y 中集合的并集 所得,→ 集合中 v 的个数 */
{
	int	l	= 1, r = len;                   /* 必须保证 Y 的并集 是 X 的并集 的子集 */
	int	p	= ranking( v );
	if ( p > len || info[p] != v )
		return(0);
	while ( l < r )
	{
		int m = (l + r) / 2;
		if ( p <= m )
			left( X, Y ), r = m;
		else
			right( X, Y ), l = m + 1;
	}
	return(value( X, Y ) );
}


int maximum( vector<int> X, vector<int> Y )     /* 计算 X 中集合的并集 减去 Y 中集合的并集 所得集合中,→ 的最大值 */
{
	int l = 1, r = len;                     /* 必须保证 Y 的并集 是 X 的并集 的子集 */
	if ( value( X, Y ) == 0 )               /* 如果集合为空则返回-inf */
		return(-inf);
	while ( l < r )
	{
		int m = (l + r) / 2, t = right_value( X, Y );
		if ( t == 0 )
			left( X, Y ), r = m;
		else
			right( X, Y ), l = m + 1;
	}
	return(info[r]);
}


int minimum( vector<int> X, vector<int> Y )     /* 计算 X 中集合的并集 减去 Y 中集合的并集 所得集合中,→ 的最小值 */
{
	int l = 1, r = len;                     /* 必须保证 Y 的并集 是 X 的并集 的子集 */
	if ( value( X, Y ) == 0 )               /* 如果集合为空则返回 inf */
		return(inf);
	while ( l < r )
	{
		int m = (l + r) / 2, t = left_value( X, Y );
		if ( t > 0 )
			left( X, Y ), r = m;
		else
			right( X, Y ), l = m + 1;
	}
	return(info[r]);
}