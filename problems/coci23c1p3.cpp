#include <bits/stdc++.h>
using namespace std;
#define int long long

int n, m;
const int MM = 4005;
int grid[MM][MM];
int r, s;
int dq[MM];
int mx[MM][MM];
int ans[MM][MM];

signed main(){
    cin.sync_with_stdio(0);
    cin.tie(0);
    
    cin >> n >> m;

    for(int i=1; i<=n; i++){
        for(int j=1; j<=m; j++){
            cin >> grid[i][j];
        }
    }

    cin >> r >> s;

    int width = m-s+1;
    int height = n-r+1;

    for(int i=1; i<=n; i++){
        memset(dq, 0, sizeof(dq));

        int head = 0, tail = -1;

        for(int j=1; j<=m; j++){
            while(head <= tail && grid[i][dq[tail]] <= grid[i][j]) tail--;

            dq[++tail] = j;

            while(head <= tail && dq[head] <= j-s) head++;

            if(j >= s){
                int left = j-s+1;
                mx[i][left] = grid[i][dq[head]];
            }
        }
    }

    for(int j=1; j<=width; j++){
        memset(dq, 0, sizeof(dq));

        int head=0, tail=-1;
        for(int i=1; i<=n; i++){
            while(head <= tail && mx[dq[tail]][j] <= mx[i][j]) tail--;
            dq[++tail] = i;

            while(head <= tail && dq[head] <= i-r) head++;

            if(i>=r){
                int top = i-r+1;
                ans[top][j] = mx[dq[head]][j];
            }
        }
    }

    for(int i=1; i<=height; i++){
        for(int j=1; j<=width; j++){
            cout << ans[i][j] << " ";
        }
        cout << endl;
    }
}