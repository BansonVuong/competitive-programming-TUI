#include <bits/stdc++.h>
using namespace std;
#define int long long

string n, h;
const int MM = 1e5*2+5;
int freqn[26];
int freqh[26];
int ans;
const int MOD = 1e18+9;
unordered_set<unsigned long long> track; // always use ull for hashes
int matches;

signed main(){
    cin.sync_with_stdio(0);
    cin.tie(0);
    
    cin >> n >> h;

    for(int i=0; i<n.size(); i++){
        char cur=n[i];
        int val=cur-'a';
        freqn[val]++;
    }

    for(int i=0; i<n.size(); i++){
        int cur = h[i]-'a';
        freqh[cur]++;
    }
    for(int i=0; i<26; i++){
        if(freqh[i] == freqn[i]) matches++;
    }
    unsigned long long hsh = 0, pw=1;
    for(int i=0; i<n.size(); i++){
        hsh = (((__int128)hsh * (__int128) 131) % MOD + h[i])%MOD;
        if(i>0) pw = (__int128) pw*131 % MOD;
    }
    for(int l=0; l+n.size()<=h.size(); l++){
        int r = l+n.size();
        // 0 length 3
        // right (next letter incoming) would be 3

        if(matches==26){
            if(track.find(hsh) == track.end()){
                track.insert(hsh);
                ans++;
            }
        }
        // leaving frequency 
        /*
            How to know if matches changed?
            
        */
        int leave = h[l]-'a';
        int leavebefore = freqh[leave];
        freqh[leave]--;
        int leaveafter = freqh[leave];

        // we lost a match
        if( (leavebefore == freqn[leave] && leaveafter != freqn[leave]) ) matches--;
        // we gained a match
        if( (leavebefore != freqn[leave] && leaveafter == freqn[leave]) ) matches++;

        if(r<h.size()) {
            int incoming = h[r]-'a';
            int incbefore = freqh[incoming];
            freqh[h[r]-'a']++;
            int incafter = freqh[incoming];

            // we lost a match
            if( (incbefore == freqn[incoming] && incafter != freqn[incoming]) ) matches--;
            // we gained a match
            if( (incbefore != freqn[incoming] && incafter == freqn[incoming]) ) matches++;
            hsh = ((__int128)hsh - (__int128)pw * h[l] % MOD + MOD*2) % MOD;// make sure to remove the leaving letter!!
            hsh=((__int128)hsh*131+h[r]) % MOD;
            
        }

        
    }

    cout <<ans << endl;
    

}