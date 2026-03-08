#include <bits/stdc++.h>
using namespace std;
#define int long long

int n, maxweight;
const int MM = 1e5+5;
int dp[MM];
int items[MM];
int value[MM], weight[MM];
signed main(){
    cin.sync_with_stdio(0);
    cin.tie(0);
    
    cin >> n >> maxweight;


    for(int i=0; i<n; i++){
        cin >> weight[i] >> value[i];
        for(int curweight=maxweight; curweight>=weight[i]; curweight--){
            dp[curweight] = max(
                dp[curweight],
                dp[curweight-weight[i]]+value[i]
            );
        }
    }

    cout << dp[maxweight] << endl;

}