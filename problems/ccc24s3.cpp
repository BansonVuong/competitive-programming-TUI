#include <bits/stdc++.h>
using namespace std;
#define int long long

int n;
const int MM = 1e5*3+5;
int first[MM], second[MM];
typedef pair<int, int> pi;
vector<pi> lswipe, rswipe;

bool cmp(pi &a, pi&b){
    return a.first>b.first;
}

signed main(){
    cin.sync_with_stdio(0);
    cin.tie(0);
    
    cin >> n;
    for(int i=0; i<n; i++){
        cin >> first[i];
    }

    for(int i=0; i<n; i++){
        cin >> second[i];
    }

    int cur = second[0], start=0;
    second[n] = -1;

    for(int f=0, s=0; f<=n; f++){
        // the previous indice was the end of the last contiguous cur
        if(second[f] != cur){
            // search for the corresponding indice. looking for the first one is fine.

            int end = f-1;

            while(first[s] != cur){
                if(s == n){ // we've reached the end.
                    cout << "NO" << endl;
                    return 0;
                }
                s++;
            }

            // s is now the first one that matches. now we have to find where s ends.
            int beginning = s;

            while(first[s] == cur && s<n){
                s++;
            }
            int ed = s-1;

            // we only need to swipe if beginning > start or end < ed
            if(beginning > start || ed < end){
                // 3 cases
                // only right swipe (s is before start)
                // only left swipe(s is after end)
                // both swipe (s is in between start and end)
                
                // why would we need swipe right? to get to end
                if(ed < end){
                    // swipe right from s to end
                    rswipe.push_back({ed, end});
                }
                if(beginning > start){
                    lswipe.push_back({start, beginning});
                }
            }

            //reset
            start = f;
            cur = second[f];
        }
    }

    cout << "YES" << endl;
    cout << rswipe.size()+lswipe.size() << endl;
    
    // order right swipes from right to left 
    sort(rswipe.begin(), rswipe.end(), cmp);
    sort(lswipe.begin(), lswipe.end());

    for(auto[s, e]:rswipe){
        printf("R %d %d\n", s, e);
    }

    for(auto[s, e]:lswipe){
        printf("L %d %d\n", s, e);
    }

}