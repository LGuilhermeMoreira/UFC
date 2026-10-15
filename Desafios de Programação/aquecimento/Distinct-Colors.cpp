#include<bits/stdc++.h>
using namespace std;


void solve() {
    int n;
    cin >> n;
    
    int max_balls = 0;
    for (int i = 0; i < n; ++i) {
        int current_balls;
        cin >> current_balls;
        max_balls = max(max_balls, current_balls);
    }
    
    cout << max_balls << "\n";
}

int main() {
    ios_base::sync_with_stdio(false);
    cin.tie(NULL);
    int t;
    cin >> t;
    while (t--) {
        solve();
    }
    return 0;
}