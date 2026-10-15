#include<bits/stdc++.h>
using namespace std;
typedef long long ll;

int main(){
    ios::sync_with_stdio(0);
    cin.tie(0);

    int n;
    ll x;
    cin >> n >> x;

    map<ll, int> seen_values;

    for(int i = 0; i < n; ++i){
        ll current_val;
        cin >> current_val;
        
        ll complement = x - current_val;
        
       
        if(seen_values.count(complement)){     
            cout << seen_values[complement] << " " << i + 1 << "\n";
            return 0;
        }
        
        
        seen_values[current_val] = i + 1;
    }
    
    
    cout << "IMPOSSIBLE\n";

    return 0;
}