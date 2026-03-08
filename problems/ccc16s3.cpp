#include <bits/stdc++.h>
using namespace std;
#define int long long

int n, m;
const int MM = 1e5+5;
vector<int> adj[MM];
vector<int> nadj[MM];

bool vis[MM];
bool pho[MM];

bool dfs(int cur){
    bool iscurpho = false;
    for(int nxt: adj[cur]){
        if(!vis[nxt]){
            vis[nxt] = 1;
            bool ispho = dfs(nxt);
            if(ispho){
                iscurpho = true;
            nadj[cur].push_back(nxt);
            nadj[nxt].push_back(cur);
            }
        }
    }    
    if(iscurpho || pho[cur]){
        return true;
    }
    return false;
}

int start;

int ans;

int dis[MM];
int par[MM];
int far, bfar;
int edges=0;
void bfs(int start){
    queue<int> q;
    memset(dis, -1, sizeof(dis));
    dis[start] = 0;
    q.push(start);
    par[start] = -1;
    while(!q.empty()){
        int cur =q.front(); q.pop();
        if(dis[cur] > far){
            far = dis[cur];
            bfar = cur;
        }
        for(int nxt:nadj[cur]){
            if(dis[nxt] == -1){
                dis[nxt] = dis[cur]+1;
                par[nxt] = cur;
                q.push(nxt);
                edges++;
            }
        }
    }
}
bool trunk[MM];
signed main(){
    cin.sync_with_stdio(0);
    cin.tie(0);
    
    cin >> n >> m;

    for(int i=0; i<m; i++){
        int u;
        cin >> u;
        pho[u] = 1;
        start = u;
    }

    for(int i=0; i<n-1; i++){
        int u, v;
        cin >> u >> v;
        adj[u].push_back(v);
        adj[v].push_back(u);
    }
    memset(vis, 0, sizeof(vis));
    vis[start] = 1;
    dfs(start);

    bfs(start);
    int d1 = bfar;
    far=0;
    edges = 0;
    bfs(d1);
    int d2 = bfar;

    cout << 2*edges-dis[d2] << endl;

    
    
}