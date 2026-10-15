#include <bits/stdc++.h>
using namespace std;
typedef long long ll;

int main()
{
    ios::sync_with_stdio(0);
    cin.tie(0);

    int n, q;
    cin >> n >> q;
    vector<ll> x(n + 1, 0), pxor(n + 1, 0);

    for (int i = 1; i <= n; i++)
    {
        cin >> x[i];
        pxor[i] = pxor[i - 1] ^ x[i];
    }

    while (q--)
    {
        int a, b;
        cin >> a >> b;
        cout << (pxor[b] ^ pxor[a - 1]) << "\n";
    }

    return 0;
}
