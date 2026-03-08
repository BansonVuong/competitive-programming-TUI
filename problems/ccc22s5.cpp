#include <bits/stdc++.h>
using namespace std;
#define int long long

int n;
const int MM = 1e5*2+5;
bool yes[MM];
int cost[MM];
int paidneedshelp[MM], paidiscovered[MM], notpaidneedshelp[MM], notpaidcovered[MM];
vector<int> adj[MM];

int vis[MM];
void dfs(int cur){
    vis[cur] = 1;
    bool children = false;
    vector<int> child;
    for(int nxt:adj[cur]){
        if(!vis[nxt]){
            dfs(nxt);
            children = true;
            child.push_back(nxt);
        }
    }

    if(!children){
        // leaf base cases:
        // paidneedshelp - impossible
        paidneedshelp[cur] = 1e10;
        // paidiscovered - possible if Y, impossible otherwise
        paidiscovered[cur] = (yes[cur]?cost[cur]:1e10);
        // notpaidneedshelp
        notpaidneedshelp[cur] = (!yes[cur]?0:1e10);
        // notpaidcovered
        notpaidcovered[cur] = (yes[cur]?0:1e10);
    }
    else{
        // paidneedshelp
        // ONLY possible if cur is not Y
        // CANT PICK:   paidiscovered - cause then parent would just be covered
        // CAN PICK:    notpaidcovered - child doesn't propagate a Y up to the parent
        //              paidneedshelp - sending help down
        //              notpaidneedshelp - sending help down
        if(!yes[cur]){
            int sum = 0;
            for(int nxt:child){
                sum += min(notpaidcovered[nxt], min(paidneedshelp[nxt], notpaidneedshelp[nxt]));
            }
            paidneedshelp[cur] = cost[cur] + sum;
        }
        else{
            paidneedshelp[cur] = 1e10;
        }
        // paidiscovered
        // if Y
        //      just grab min of all 4 states of children
        // if not, we need a child to contribute a Y
        //      for each child, choose it to contribute a Y, and then just take the min of all 4 states of all other children
        int allstates = 0;
        for(int nxt:child){
            allstates += min(min(notpaidcovered[nxt], notpaidneedshelp[nxt]), min(paidiscovered[nxt], paidneedshelp[nxt]));
        }

        if(yes[cur]){
            paidiscovered[cur] = allstates+cost[cur];
        }
        else{
            paidiscovered[cur] = 1e10;
            for(int nxt:child){
                int childcontrib = min(
                    min(notpaidcovered[nxt], notpaidneedshelp[nxt]),
                    min(paidiscovered[nxt], paidneedshelp[nxt])
                );
                paidiscovered[cur] = min(
                    paidiscovered[cur],
                    (
                        cost[cur]
                        + paidiscovered[nxt]
                        + allstates
                        - childcontrib
                    )
                );
            }
        }

        // notpaidiscovered
        // if yes? just child sum
        // otherwise
        // again have to pick one child to contrib 
        int twostates = 0;
        for(int nxt:child){
            twostates += min(
                paidiscovered[nxt],
                notpaidcovered[nxt]
            );
        }
        if(yes[cur]){
            notpaidcovered[cur] = twostates;
        }
        else{
            notpaidcovered[cur] = 1e10;
            for(int nxt:child){
                notpaidcovered[cur] = min(
                    notpaidcovered[cur],
                    (
                        twostates
                        + paidiscovered[nxt]
                        - min(
                            paidiscovered[nxt],
                            notpaidcovered[nxt]
                        )
                    )
                );
            }
        }

        // notpaidneedshelp
        // children can be whatever 
        notpaidneedshelp[cur] = twostates;
    }
}
signed main(){
    cin.sync_with_stdio(0);
    cin.tie(0);
    
    cin >> n;
    for(int i=1; i<n;  i++){
        int u, v;
        cin >> u >> v;
        adj[u].push_back(v);
        adj[v].push_back(u);
    }

    for(int i=1; i<=n; i++){
        char k;
        cin >> k;
        yes[i] = k=='Y';
    }

    for(int i=1; i<=n; i++){
        cin >> cost[i];
    }

    dfs(1);

    cout << min(
        paidiscovered[1],
        notpaidcovered[1]
    ) << endl;
}