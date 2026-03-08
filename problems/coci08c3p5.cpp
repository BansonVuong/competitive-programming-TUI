#include <bits/stdc++.h>
using namespace std;
#define int long long

const int MM = 1e5*3+5;

int n;
int root;

pair<int, int> nodes[MM];

int ans;

int depth[MM];

set<int> inserted;

void insert(int x, int n){
    ans++;
    if(x < n){
        if(nodes[n].first == 0){
            nodes[n].first = x;
        }
        else{
            insert(x, nodes[n].first);
        }
    }
    else{
        if(nodes[n].second == 0){
            nodes[n].second = x;
        }
        else{
            insert(x, nodes[n].second);
        }
    }
}

signed main(){
    cin.sync_with_stdio(0);
    cin.tie(0);

    cin >> n;

    cin >> root;
    inserted.insert(root);
    cout << 0 << endl;
    for(int i=1; i<n; i++){
        int k;
        cin >> k;
        auto nxt = inserted.lower_bound(k);
        bool isnxt = (nxt != inserted.end());
        bool isprv = (nxt != inserted.begin());
        int succ = *nxt;
        int pred = *prev(nxt);

        depth[k] = 1 + max((isnxt?depth[succ]:-1), (isprv?depth[pred]:-1));
        ans += depth[k];
        inserted.insert(k);
        cout <<ans << endl;
    }
}