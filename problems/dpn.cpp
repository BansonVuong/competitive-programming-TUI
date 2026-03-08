#include <bits/stdc++.h>
using namespace std;
#define int long long

int n;
const int MM = 405;
int dp[MM][MM];
int val[MM];
int psa[MM];
signed main(){
    cin.sync_with_stdio(0);
    cin.tie(0);
    
    cin >> n ;
    for(int i=0; i<MM; i++){
        for(int j=0; j<MM; j++){
            dp[i][j] = 1e18;
        }
    }
    for(int i=1; i<=n; i++){
        cin >> val[i];
        psa[i] = val[i]+psa[i-1];
        dp[i][i] = 0;
    }
    
    // dp[i][i] = 0;
    for(int len=2; len<=n; len++){
        for(int left=1; left+len-1<=n; left++){
            int right = left+len-1;
            for(int split=left; split<=right; split++){
                dp[left][right] = min(
                    dp[left][right],
                    psa[right]-psa[left-1] + dp[left][split] + dp[split+1][right]
                );
            }
        }
    }
    cout << dp[1][n] << endl;
    
}