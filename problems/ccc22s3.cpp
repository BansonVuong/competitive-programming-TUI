#include <bits/stdc++.h>
using namespace std;
#define int long long

int n, m, k;
const int MM = 1e6+5;
int ans[MM];
int len[MM];
int sum;
signed main(){
    cin.sync_with_stdio(0);
    cin.tie(0);
    
    cin >> n >> m >> k;

    if(k < n){
        cout << -1 << endl;
        return 0;
    }
    if(k > n*(n+1)/2 - (n-min(n,m))*(n-min(n,m)+1)/2){
        cout << -1 << endl;
        return 0; 
    }
    
    for(int i=1; i<=n; i++){
        int prebudget = k-n+i;
        int budget = prebudget-sum;
        budget = min(budget, i);
        budget = min(budget, m);
        budget = min(budget, len[i-1]+1);
        len[i] = budget;
        if(budget == len[i-1]+1){
            cout << budget << " ";
            ans[i] = budget;
            sum += budget;
        }
        else{
            cout << ans[i-budget] << " ";
            ans[i] = ans[i-budget];
            sum += budget;
        }
        
    }

    cout << endl;

}

/*
1 2 3 4 1
12, 23, 34, 42
123, 234, 342
1234, 2341
*/