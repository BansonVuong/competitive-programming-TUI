#include <bits/stdc++.h>
using namespace std;
#define int long long

typedef pair<int, int> pi;

int n, k;
const int MM = 3*1e6+5;
multiset<int> bags;
vector<pair<int, int>> gems;
int ans;

bool cmp(const pi &a, const pi &b){
    if(a.first != b.first) return a.first > b.first;
    return a.second < b.second;
}
signed main(){
    cin.sync_with_stdio(0);
    cin.tie(0);
    
    cin >> n >> k;
    for(int i=0; i<n; i++){
        int u , v;
        cin >> u >> v;
        gems.push_back({v, u});
    }

    for(int i=0; i<k; i++){
        int u;
        cin >> u;
        bags.insert(u);
    }

    sort(gems.begin(), gems.end(), cmp);

    for(auto[v, m] : gems){
        auto it = bags.lower_bound(m);
        if(it == bags.end()){
            continue;
        }
        ans += v;
        bags.erase(it);
    }

    cout << ans << endl;

    


}