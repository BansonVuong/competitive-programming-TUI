#include <bits/stdc++.h>
using namespace std;
#define int long long

string k;
string s;
const int MM = 2*1e6+5;

int psa[MM][3];
int ch[3];

int ans = 1e18;

signed main(){
    cin.sync_with_stdio(0);
    cin.tie(0);

    cin >> k;
    s = k+k;
    int first = s[0]-'A';
    ch[first]++;
    psa[0][first] = 1;
    int size = s.size();
    for(int i=1; i<s.size(); i++){
        int cur = s[i]-'A';      
        for(int j=0; j<3; j++){
            psa[i][j] = psa[i-1][j];
        }
        psa[i][cur]++;
        if(i < k.size()){
            ch[cur]++;
        }
    }

    int a = ch[0], b = ch[1], c = ch[2];

    for(int i=1; i<=k.size(); i++){
        int bInA, cInA;
        
        bInA = psa[i+a-1][1] - psa[i-1][1];
        cInA = psa[i+a-1][2] - psa[i-1][2];

        // case 1: B follows first
        int aInB1, cInB1;
        int aInC1, bInC1;

        aInB1 = psa[i+a+b-1][0]-psa[i+a-1][0];
        cInB1 = psa[i+a+b-1][2]-psa[i+a-1][2];

        aInC1 = psa[i+a+b+c-1][0]-psa[i+a+b-1][0];
        bInC1 = psa[i+a+b+c-1][1]-psa[i+a+b-1][1];

        int abs1, bcs1, acs1;
        
        abs1 = min(aInB1, bInA);
        bcs1 = min(bInC1, cInB1);
        acs1 = min(aInC1, cInA);

        ans = min(ans, (
            abs1+bcs1+acs1+(2*(max(aInB1, bInA)-abs1))
        ));
        // case 2: C follows first

        int aInC2, bInC2;
        int aInB2, cInB2;

        aInC2 = psa[i+a+c-1][0] - psa[i+a-1][0];
        bInC2 = psa[i+a+c-1][1] - psa[i+a-1][1];

        aInB2 = psa[i+a+b+c-1][0] - psa[i+a+c-1][0];
        cInB2 = psa[i+a+b+c-1][2] - psa[i+a+c-1][2];

        int abs2, bcs2, acs2;

        abs2 = min(aInB2, bInA);
        bcs2 = min(bInC2, cInB2);
        acs2 = min(aInC2, cInA);

        ans = min(ans, (
            abs2+bcs2+acs2+(2*(max(aInB2, bInA)-abs2))
        ));

    }

    cout << ans <<endl;
}