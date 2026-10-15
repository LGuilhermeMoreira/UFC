#include <bits/stdc++.h>
using namespace std;
typedef long long ll;

int main()
{
    ios::sync_with_stdio(0);
    cin.tie(0);

    int n;
    cin >> n;
    vector<ll> a(n);
    for (int i = 0; i < n; i++)
        cin >> a[i];

    vector<ll> freq(n, 0);
    ll prefix = 0;
    freq[0] = 1; 
    for (int i = 0; i < n; i++)
    {
        prefix += a[i];
        ll mod = ((prefix % n) + n) % n;
        freq[mod]++;
    }

    ll ans = 0;
    for (int i = 0; i < n; i++)
    {
        ans += freq[i] * (freq[i] - 1) / 2;
    }

    cout << ans << "\n";
    return 0;
}
