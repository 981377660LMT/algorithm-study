#include <cstdio>
#include <cstring>
#include <iostream>
using namespace std;

const int MAXN = 500000 , MAXK = 2;
int n;
char str[ MAXN + 5 ];
struct Palindrome_Automaton{
	int Size , Last , Root0 , Root1 , Trans[ MAXN + 5 ][ MAXK + 5 ] , Link[ MAXN + 5 ];
	int Len[ MAXN + 5 ] , Cnt[ MAXN + 5 ];

	Palindrome_Automaton( ) {
		Root0 = Size ++ , Root1 = Size ++; Last = Root0;
		Len[ Root0 ] = 0  , Link[ Root0 ] = Root1;
		Len[ Root1 ] = -1 , Link[ Root1 ] = Root1; 
	}
	void Extend( int ch , int dex ) {
		int u = Last;
		for( ; u != 1 && ( str[ dex - Len[ u ] - 1 ] == str[ dex ] || ( dex - Len[ u ] - 1 == -1 ) ) ; u = Link[ u ] );
		if( str[ dex - Len[ u ] - 1 ] == str[ dex ] || ( dex - Len[ u ] - 1 == -1 ) ) { Last = 0; return; }
        if( !Trans[ u ][ ch ] ) {
			int Newnode = ++ Size , v = Link[ u ];
			Len[ Newnode ] = Len[ u ] + 2;
			
			for( ; v != 1 && ( str[ dex - Len[ v ] - 1 ] == str[ dex ] || ( dex - Len[ u ] - 1 == -1 ) ) ; v = Link[ v ] );
            Link[ Newnode ] = str[ dex - Len[ u ] - 1 ] == str[ dex ] ? 0 : Trans[ v ][ ch ];
            Trans[ u ][ ch ] = Newnode;
			Cnt[ Newnode ] = Cnt[ Link[ Newnode ] ] + 1;
		}
		Last = Trans[ u ][ ch ];
	}
	
	long long Build( char *str ) {
		int len = strlen( str );
        long long Ans = 0;
		for( int i = 0 ; i < len ; i ++ ) {
			Extend( str[ i ] - '0' + 1 , i );
			Ans += Cnt[ Last ];
		}	
        return Ans;
	}
}PAM;

int main() {
    scanf("%d",&n);
    scanf("%s", str );
    printf("%lld", PAM.Build( str ) );
    return 0;
}