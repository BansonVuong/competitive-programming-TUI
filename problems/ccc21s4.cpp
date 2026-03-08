#include <bits/stdc++.h>
using namespace std;
#define int long long

int n, w, d;
const int MM = 1e5*2+5;
vector<int> adj[MM];
multiset<int> ans;

int station[MM];
int dis[MM];

signed main(){
    cin.sync_with_stdio(0);
    cin.tie(0);
    
    cin >> n >> w >> d;

    for(int i=0; i<w; i++){
        int u, v;
        cin >> u >> v;
        adj[v].push_back(u);
    }

    memset(dis, -1, sizeof(dis));

    dis[n] = 0;
    queue<int> q;
    q.push(n);
    while(!q.empty()){
        int cur = q.front(); q.pop();
        for(int nxt:adj[cur]){
            if(dis[nxt] == -1){
                q.push(nxt);
                dis[nxt] = dis[cur]+1;
            }
        }
    }

    for(int i=1; i<=n; i++){
        cin >> station[i];
        ans.insert(dis[station[i]]==-1?1e18:dis[station[i]]+i);
    }

    for(int i=0; i<d; i++){
        int a, b;
        cin >> a >> b;

        int ba = dis[station[a]] == -1?1e18:dis[station[a]]+a;
        int bb = dis[station[b]] == -1? 1e18:dis[station[b]]+b;

        ans.erase(ans.find(ba));
        ans.erase(ans.find(bb));

        swap(station[a], station[b]);

        int aa = dis[station[a]] == -1? 1e18 : dis[station[a]]+a;
        int ab = dis[station[b]] == -1?1e18:dis[station[b]]+b;

        ans.insert(aa);
        ans.insert(ab);
        cout << *ans.begin()-1 << endl;
        
    }
}