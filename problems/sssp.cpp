#include <bits/stdc++.h>
using namespace std;
#define int long long

int n, m;
const int MM = 5005;
typedef pair<int, int> pi;
vector<pi> adj[MM];
int dis[MM];
bool vis[MM];
signed main(){
    cin.sync_with_stdio(0);
    cin.tie(0);
    
    memset(dis, 0x3f3f3f3f, sizeof(dis));

    cin >> n >> m;
    for(int i=0; i<m; i++){
        int u, v, w;
        cin >> u >> v >> w;
        adj[u].push_back({v, w});
        adj[v].push_back({u, w});
    }

    priority_queue<pi, vector<pi>, greater<pi>> q;
    q.push({0, 1});
    vis[1] = 1;
    dis[1] = 0;
    while(!q.empty()){
        auto[d, u] = q.top(); q.pop();

        for(auto[v, w]:adj[u]){
            if(dis[u]+w < dis[v]){
                dis[v] = dis[u]+w;
                vis[v] = 1;
                q.push({dis[v], v});
            }
        }
    }

    for(int i=1; i<=n; i++){
        cout << (vis[i]?dis[i]:-1) << endl;
    }
}