#include <bits/stdc++.h>
using namespace std;
#define int long long

int n, m;
const int MM = 1e5+5;
int par[MM];
vector<int> ans;
int find_par(int u){
    if(par[u] == u) return u;
    return par[u] = find_par(par[u]);
}

void merge(int u, int v){
    int fu = find_par(u), fv=find_par(v);
    par[fu] = fv;
}
signed main(){
    cin.sync_with_stdio(0);
    cin.tie(0);
    
    cin >> n >> m;

    for(int i=1; i<=n; i++){
        par[i] = i;
    }

    for(int i=1; i<=m; i++){
        int u, v;
        cin >> u >> v;
        if(find_par(u) != find_par(v)){
            merge(u, v);
            ans.push_back(i);
        }
    }

    for(int i=2; i<=n; i++){
        if(find_par(i) != find_par(i-1)){
            cout << "Disconnected Graph" << endl;
            return 0;
        }
    }

    for(int u:ans){
        cout << u << endl;
    }
}