#include <bits/stdc++.h>
using namespace std;
#define int long long

/*
3:03 start
subtasks:
1:
each road is connected to the one before and after it
ex. 1-2, 2-3, ... N-1-N
therefore the solution is to alternate red-blue edges on the edges going to/from node x-1 and x+1, and all other edges are grey
2. there is only one extra edge, forming a cycle somewhere. that extra edge is grey.
subtask 3: each cycle doesn't have to worry about having to handle multiple paths for the same grey edge. each cycle can be treated independently - one grey path within the cycle, coloured roads connecting separate cycles

full soln: must worry about each edge interacting with other cycles. i havent figured that part out yet.

insight: essentially, to minimize paintings, in any cycle, one edge will be grey, and the others alternating

insight 2: essentially, cut an edge w/in every cycle and you're left with a tree

*/

/*
40 minutes spent:
edge from odd parent to even child: red
edge from even parent to odd child: blue
*/

const int MM = 1e5*2+5;

int n, m;

vector<pair<int, int>> adj[MM];
unordered_set<int> q;
char road[MM];
bool vis[MM];

void dfs(int cur, bool par){
    q.erase(cur);
    for(auto[nxt, idx] : adj[cur]){
        if(!vis[nxt]){
            vis[nxt] = 1;
            if(par){
                road[idx] = 'B';
            }
            else{
                road[idx] = 'R';
            }
            dfs(nxt, !par);
        }
    }
}
signed main(){
    cin.tie(0);
    cin.sync_with_stdio(0);

    cin >> n >>m;

    for(int i=0; i<m; i++){
        int u, v;
        cin >> u >>v;
        
        q.insert(u);
        q.insert(v);

        adj[u].push_back({v, i});
        adj[v].push_back({u, i});
        road[i] = 'G';
    }

    while(!q.empty()){
        int cur = *q.begin();
        vis[cur] = 1;
        dfs(cur, 0);
    }
    
    for(int i=0; i<m; i++){
        cout <<road[i];
    }
    cout << endl;
}