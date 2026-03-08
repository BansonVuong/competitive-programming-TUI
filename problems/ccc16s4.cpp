#include <bits/stdc++.h>
using namespace std;
#define int long long

int n;
const int MM = 405;
int dp[MM][MM];
int val[MM], psa[MM];
int ans;
signed main(){
    cin.sync_with_stdio(0);
    cin.tie(0);
    
    cin >> n;
    
    for(int i=1; i<=n; i++){
        cin >> val[i];
        psa[i] = val[i] + psa[i-1];
        dp[i][i] = val[i];
        ans = max(ans, dp[i][i]);
    }
    
    // so we need a left, a right
    for(int len=2; len<=n; len++){
        for(int leftstart=1; leftstart+len-1<=n; leftstart++){
            int rightend = leftstart+len-1; // 1, length 2: 1+2 -> 3, 3-1
            int sectionlength = rightend-leftstart+1; // 10-1 : 9 -> 9+1

            // case 1: there is no ball in between

            for(int leftEnd=leftstart; leftEnd<rightend; leftEnd++){
                int rightStart = leftEnd+1;

                bool isleftball = psa[leftEnd] - psa[leftstart-1] == dp[leftstart][leftEnd];
                bool isrightball = psa[rightend] - psa[rightStart-1] == dp[rightStart][rightend];
                int leftball = dp[leftstart][leftEnd];
                int rightball = dp[rightStart][rightend];
                bool ballsEqual = leftball == rightball;

                if(isleftball && isrightball && ballsEqual){
                    dp[leftstart][rightend] = max(dp[leftstart][rightend], psa[rightend] - psa[leftstart-1]);
                }
            }
            // case 2: there is a ball in between, and that ball needs a length

            if(len >= 3){
                for(int middlelength = 1; middlelength <= sectionlength-2; middlelength++){
                    // 1, 2, 3
                    // pos = 2
                    // 2+1-1 < 3
                    for(int middleStart=leftstart+1; middleStart+middlelength-1 < rightend; middleStart++){
                        int middleEnd = middleStart+middlelength-1;
                        int leftEnd = middleStart-1;
                        int rightStart = middleEnd+1;
                        
                        bool isleftball = psa[leftEnd] - psa[leftstart-1] == dp[leftstart][leftEnd];
                        bool isrightball = psa[rightend] - psa[rightStart-1] == dp[rightStart][rightend];
                        bool ismiddleball = psa[middleEnd] - psa[middleStart-1] == dp[middleStart][middleEnd];
                        int leftball = dp[leftstart][leftEnd];
                        int rightball = dp[rightStart][rightend];
                        bool ballsEqual = leftball == rightball;

                        if(isleftball && isrightball && ismiddleball && ballsEqual){
                            dp[leftstart][rightend] = max(dp[leftstart][rightend], psa[rightend] - psa[leftstart-1]);
                        }
                        
                    }
                }
            }
            ans = max(ans, dp[leftstart][rightend]);
        }
    }

    cout << ans << endl;
}