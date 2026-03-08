#include <bits/stdc++.h>
using namespace std;
#define int long long

const int MM = 1e3*2+5;

int par[MM], dis[MM];

int n, m, ans;

typedef pair<int, int> pi;

vector<pi> adj[MM];

int find_par(int u){
    if(par[u] == u) return u;
    par[u] = find_par(par[u]);
    return par[u];
}

void merge(int u, int v){
    int paru = find_par(u), parv = find_par(v);
    par[paru] = parv;
}

vector<pair<pi, pi>> edges;

bool cmp(pair<pi, pi> a, pair<pi, pi> b){
    if (a.second.first == b.second.first)
        return a.second.second < b.second.second;
    return a.second.first < b.second.first;
}

int dijk(int st, int ed, int ln){
    memset(dis, 0x3f3f3f3f, sizeof(dis));
    priority_queue<pi, vector<pi>, greater<pi>> q;
    q.push({0, st});
    dis[st] = 0;
    while(!q.empty()){
        auto [d, u] = q.top(); q.pop();
        if(d > ln) return -1;
        if(u == ed) return d;
        for(auto[v, w] : adj[u]){
            if(dis[u] + w < dis[v]){
                dis[v] = dis[u]+w;
                q.push({dis[v], v});
            }
        }
    }
    return -1;
}
signed main(){
    cin.sync_with_stdio(0);
    cin.tie(0);

    cin >> n >> m;

    for(int i=0; i<=n; i++){
        par[i] = i;
    }
    for(int i=0; i<m; i++){
        int u, v, l, c;
        cin >> u >> v >> l >> c;

        edges.push_back({{u, v}, {l, c}});
    }

    sort(edges.begin(), edges.end(), cmp);

    for(auto[n, w] : edges){
        auto [u, v] = n;
        auto [l, c] = w;

        if(find_par(u) != find_par(v)){
            ans += c;
            merge(u, v);
            adj[u].push_back({v, l});
            adj[v].push_back({u, l});
        }
        else{
            int ds = dijk(u, v, l);
            if(ds == -1 || ds > l){
                ans += c;
                merge(u, v);
                adj[u].push_back({v, l});
                adj[v].push_back({u, l});
            }
        }

    }

    cout << ans << endl;
}