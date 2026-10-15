#include <iostream>
#include <vector>
#include <map>
#include <algorithm> 
using namespace std;

void solve() {
    long long n;
    cin >> n;

    if (n == 0) {
        cout << 0 << "\n";
        return;
    }

    map<long long, long long> mp;
    for (long long i = 0; i < n; i++) {
        long long num;
        cin >> num;
        mp[num]++; 
    }

    long long max_freq = 0;
    for (auto const& x : mp) {
        max_freq = max(max_freq, x.second);
    }

    cout << n - max_freq << "\n";
}

int main() {
    ios_base::sync_with_stdio(false);
    cin.tie(NULL);

    long long t;
    cin >> t;
    while (t--) {
        solve();
    }

    return 0;
}