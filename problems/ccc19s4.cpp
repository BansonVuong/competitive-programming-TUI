#include <bits/stdc++.h>
using namespace std;
#define int long long

int n, k;
const int MM = 1e6+5;
const int INF = 1e18;
int attr[MM];
int dp[MM];
int prefixmax[MM], suffixmax[MM];
int bestdp[MM], bestsum[MM];

signed main(){
    cin.sync_with_stdio(0);
    cin.tie(0);
    cin >> n>> k;
    for(int i=0; i<n; i++){
        cin >>attr[i];
        dp[i] = -INF;
    }

    int slack = (k-n%k)%k;
    int days = (n+k-1)/k;
    int mday = k-slack;

    prefixmax[0] = attr[0];
    for(int i=1; i<k; i++){
       prefixmax[i] = max(prefixmax[i-1], attr[i]);
    }

    for(int i=mday-1; i<k; i++){
        dp[i] = prefixmax[i];
    }

    for(int d=1; d<days; d++){
        int l = d*k;
        int r = min(n, (d+1) * k);

        for(int i=l; i<r; i++){
            if(i==l) prefixmax[i] = attr[i];
            else prefixmax[i] = max(prefixmax[i-1], attr[i]);
        }

        int pl = (d-1)*k;
        int pr = d*k-1;
        for(int i=pr; i>=pl; i--){
            if(i == pr) suffixmax[i] = attr[i];
            else suffixmax[i] = max(suffixmax[i+1], attr[i]);
        }

        bestdp[pr+1] = -INF;
        bestsum[pr+1] = -INF;

        for(int i=pr; i>=pl; i--){
            bestdp[i] = max(dp[i], bestdp[i+1]);

            int left = (i==0?0:dp[i-1]);
            bestsum[i] = max(left + suffixmax[i], bestsum[i+1]);
        }

        int start = d*k + mday-1;
        int finish = r-1;
        if(start > finish) continue;

        int it = max(pl, start-k);

        for(int i=start; i<=finish; i++){
            while(i-it > k) it++;

            int v1 = prefixmax[i] + bestdp[it];
            int v2 = bestsum[it+1];

            dp[i] = max(v1, v2);
        }
    }

    cout <<dp[n-1] <<endl;

}