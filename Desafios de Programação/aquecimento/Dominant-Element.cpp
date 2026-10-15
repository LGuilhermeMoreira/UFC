#include <bits/stdc++.h> 
using namespace std;

void solve() {
    int n;
    cin >> n;


    unordered_map<int, int> freq;
    for (int i = 0; i < n; ++i) {
        int num;
        cin >> num;
        freq[num]++;
    }

    if (freq.size() == 1) {
        
        cout << "YES\n";
        return;
    }

    int max_freq = 0;
    int second_max_freq = 0;

   
    for (auto const& pair : freq) {
        if (pair.second > max_freq) {
            second_max_freq = max_freq;
            max_freq = pair.second;
        } else if (pair.second > second_max_freq) {
            second_max_freq = pair.second;
        }
    }

    if (max_freq > second_max_freq) {
        cout << "YES\n";
    } else {
        cout << "NO\n";
    }
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