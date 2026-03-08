#include <bits/stdc++.h>
using namespace std;
#define int long long

int r, c;
const int MM = 105;
char grid[MM][MM];


typedef pair<int, int> pi;
vector<pi> cameras;
pi start;
vector<pi> conveyors;

int vis[MM][MM];
bool dfs(int curx, int cury){
    if(vis[curx][cury] == 1) return grid[curx][cury] != 'W';
    char cur = grid[curx][cury];
    if(cur == 'W'){
        return false;
    }
    if(cur == '.' || cur == 'S'){
        return true;
    }
    vis[curx][cury] = 2;
    if(cur == 'L'){
        int nx=curx, ny= cury-1;
        if(vis[nx][ny] == 2){
            return false;
        }
        bool good = dfs(nx, ny);
        vis[curx][cury] = 1;
        if(!good) grid[curx][cury] = 'W';
    }
    if(cur == 'R'){
        int nx=curx, ny= cury+1;
        if(vis[nx][ny] == 2){
            return false;
        }
        bool good = dfs(nx, ny);
        vis[curx][cury] = 1;
        if(!good) grid[curx][cury] = 'W';
    }
    if(cur == 'U'){
        int nx=curx-1, ny=cury;
        if(vis[nx][ny] == 2){
            return false;
        }
        bool good = dfs(nx, ny);
        vis[curx][cury] = 1;
        if(!good) grid[curx][cury] = 'W';
    }
    if(cur == 'D'){
        int nx=curx+1, ny=cury;
        if(vis[nx][ny] == 2){
            return false;
        }
        bool good = dfs(nx, ny);
        vis[curx][cury] = 1;
        if(!good) grid[curx][cury] = 'W';
    }
    vis[curx][cury] = 1;
    return true;
}
void output(int curx, int cury){
    for(int i=1; i<=r; i++){
        for(int j=1; j<=c; j++){
            cout << grid[i][j] << ' ';
        }
        cout << endl;
    }
    cout << "---------" << endl;
}
int dis[MM][MM];
vector<pi> ans;
int directions[4][2] ={
    {0, 1},
    {0,-1},
    {-1,0},
    { 1,0},
};
signed main(){
    cin.sync_with_stdio(0);
    cin.tie(0);
    
    cin >> r >> c;
    for(int i=1; i<=r; i++){
        for(int j=1; j<=c; j++){
            cin >> grid[i][j];
            char cur = grid[i][j];
            if(cur == 'C'){
                cameras.push_back({i, j});
            }
            if(cur == 'S'){
                start = {i, j};
            }
            if(cur == 'L' || cur == 'R' || cur == 'D' || cur == 'U'){
                conveyors.push_back({i, j});
            }
            if(cur == '.'){
                ans.push_back({i, j});
            }
        }
    }

    for(auto[sx, sy]:cameras){
        bool move = true;
        for(int i=sx; i<=r && move; i++){
            if(grid[i][sy] == 'S'){
                for(int j=0; j<ans.size(); j++) cout << -1 << endl;
                return 0;
            }
            if(grid[i][sy] == 'W'){
                move = false;
            }
            else if(grid[i][sy] == '.'){
                grid[i][sy] = 'B';
            }
            
        }
        move= true;
        for(int i=sx; i>0 && move; i--){
            if(grid[i][sy] == 'S'){
                for(int j=0; j<ans.size(); j++) cout << -1 << endl;
                return 0;
            }
            if(grid[i][sy] == 'W'){
                move = false;
            }
            else if(grid[i][sy] == '.'){
                grid[i][sy] = 'B';
            }
            
        }
        move= true;
        for(int i=sy; i<=c && move; i++){
            if(grid[sx][i] == 'S'){
                for(int j=0; j<ans.size(); j++) cout << -1 << endl;
                return 0;
            }
            if(grid[sx][i] == 'W'){
                move = false;
            }
            else if(grid[sx][i] == '.'){
                grid[sx][i] = 'B';
            }
            
        }
        move= true;
        for(int i=sy; i>0 && move; i--){
            if(grid[sx][i] == 'S'){
                for(int j=0; j<ans.size(); j++) cout << -1 << endl;
                return 0;
            }
            if(grid[sx][i] == 'W'){
                move = false;
            }
            else if(grid[sx][i] == '.'){
                grid[sx][i] = 'B';
            }
            
        }
    }
    for(auto[sx, sy]: cameras){
        grid[sx][sy] = 'W';
    }
    for(int i=1; i<=r; i++){
        for(int j=1; j<=c; j++){
            if(grid[i][j] == 'B') grid[i][j] = 'W';
        }
    }

    for(auto[sx, sy]:conveyors){
        dfs(sx, sy);
    }

    queue<pi> q;
    memset(dis, -1, sizeof(dis));
    q.push({start});
    dis[start.first][start.second] = 0;
    while(!q.empty()){
        auto[curx, cury] = q.front(); q.pop();
        for(int i=0; i<4; i++){
            int nx = directions[i][0]+curx, ny = directions[i][1]+cury;

            if(grid[nx][ny] != 'W' && dis[nx][ny] == -1){
                while(grid[nx][ny] != '.' && dis[nx][ny] == -1){
                    dis[nx][ny] = dis[curx][cury]+1;
                    if(grid[nx][ny] == 'L'){
                        ny--;
                    }
                    else if(grid[nx][ny] == 'R'){
                        ny++;
                    }
                    else if(grid[nx][ny] == 'U'){
                        nx--;
                    }
                    else if(grid[nx][ny] == 'D'){
                        nx++;
                    }
                }
                if(dis[nx][ny] == -1){
                    dis[nx][ny] = dis[curx][cury]+1;
                    q.push({nx, ny});
                }

            }
        }
    }
    for(auto[x, y]: ans){
        cout << dis[x][y] << endl;
    }
}