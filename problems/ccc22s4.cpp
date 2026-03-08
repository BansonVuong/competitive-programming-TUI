#include <bits/stdc++.h>
using namespace std;
#define int long long

int n, c;

int good;

vector<int> circ;

signed main(){
    cin.sync_with_stdio(0);
    cin.tie(0);

    cin >> n >> c;
    good = (n*(n-1)*(n-2))/6;

    for(int i=0; i<n; i++){
        int u;
        cin >> u;
        circ.push_back(u);
        circ.push_back(u+c);
    }
    sort(circ.begin(), circ.end());
    for(int i=0; i<n; i++){
        int cur = circ[i];
        int up = c/2+cur-(!(c%2));
        auto pos = circ.begin()+i+1;
        auto end = upper_bound(circ.begin(), circ.end(), up);
        //1a:
        int cnt = end-pos;
        int choose = (cnt*(cnt-1))/2;
        good -= choose;
    }
    if(c%2==0){
        for(int i=0; i<c/2; i++){

            int up = c/2+i;

            auto lstart = lower_bound(circ.begin(), circ.end(), i);
            auto rstart = upper_bound(circ.begin(), circ.end(), i);
            auto lend = lower_bound(circ.begin(), circ.end(), up);
            auto rend = upper_bound(circ.begin(), circ.end(), up);

            //2a:
            int bg = rstart-lstart;
            int ed = rend-lend;
            int md = lend-rstart;
            good-= bg*ed*(n-bg-ed);
            //2b:
            good -= (bg*(bg-1))/2*ed;
            //2c:
            good -= (ed*(ed-1))/2*bg;


        }
    }

    cout <<good << endl;
}