#include <bits/stdc++.h>
using namespace std;
typedef long long ll;

/* 
Função para verificar se é possível dividir o array em no máximo k subarrays, 
tal que a soma de cada subarray seja <= max_allowed_sum
*/
bool check(const vector<ll>& x, ll k, ll max_allowed_sum) {
    ll subarrays_needed = 1;
    ll current_sum = 0;

    for (ll val : x) {
        if (val > max_allowed_sum) return false;

        if (current_sum + val <= max_allowed_sum) {
            current_sum += val;
        } else {
            subarrays_needed++;
            current_sum = val;
        }
    }

    return subarrays_needed <= k;
}

int main() {
    ios::sync_with_stdio(0);
    cin.tie(0);

    int n;
    ll k;
    cin >> n >> k;

    vector<ll> x(n);
    ll total_sum = 0;
    ll max_element = 0;

    for (int i = 0; i < n; ++i) {
        cin >> x[i];
        total_sum += x[i];
        max_element = max(max_element, x[i]);
    }

    ll low = max_element;
    ll high = total_sum; 
    ll ans = high;

    while (low <= high) {
        ll mid = low + (high - low) / 2;

        if (check(x, k, mid)) {
            ans = mid;
            high = mid - 1;
        } else {
            low = mid + 1;
        }
    }

    cout << ans << endl;
    return 0;
}
