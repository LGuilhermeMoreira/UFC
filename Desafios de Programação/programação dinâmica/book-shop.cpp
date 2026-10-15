#include <bits/stdc++.h>
using namespace std;
typedef long long ll;
const int MAXN = 1010;
const int MAXX = 100010;
const int INF = MAXN * MAXX;
int dp[MAXN][MAXX];
int n, x;
vector<int> h(MAXN);
vector<int> s(MAXN);

int solve(int i, int cap)
{
    if (cap < 0)
    {
        return -INF;
    }
    if (i == n)
    {
        return 0;
    }
    if (dp[i][cap] != -1)
    {
        return dp[i][cap];
    }

    int colocar = s[i] + solve(i + 1, cap - h[i]);
    int n_colocar = solve(i + 1, cap);

    return dp[i][cap] = max(colocar, n_colocar);
}

int main()
{
    ios::sync_with_stdio(0);
    cin.tie(0);

    cin >> n >> x;

    for (int i = 0; i <= n; i++)
    {
        for (int j = 0; j <= x; j++)
        {
            dp[i][j] = -1;
        }
    }

    for (int i = 0; i < n; i++)
    {
        cin >> h[i];
    }

    for (int i = 0; i < n; i++)
    {
        cin >> s[i];
    }

    cout << solve(0, x) << "\n";

    return 0;
}
