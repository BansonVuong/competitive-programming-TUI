#include <bits/stdc++.h>
using namespace std;
#define int long long
int n;

const int MM= 1e6+5;


unordered_map<int, int> dip;


int dp(int n){
    if(n == 1 || n == 2) return 1;
    auto it = dip.find(n);
    if(it != dip.end()) return it->second;

    int ans=0;

    for(int j=2; j<=n; j++){

        int diff = (n/(n/j))-j;


        ans += (1+diff)*dp(n/j);
        j = n/(n/j);
    }

    dip[n] = ans;
    return dip[n];
    

}
signed main(){
    cin.sync_with_stdio(0);
    cin.tie(0);

    // dip.max_load_factor(0.5);
    // dip.reserve(200000);

    cin >>n;



    cout <<dp(n) << endl;
}