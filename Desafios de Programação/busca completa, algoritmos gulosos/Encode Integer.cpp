#include <bits/stdc++.h>
using namespace std;
typedef long long ll;

int main()
{
    ios::sync_with_stdio(0);
    cin.tie(0);

    int n;
    cin >> n;

    for (size_t i = 0; i < n; i++)
    {
        int N;
        cin >> N;

        if (N == 0)
        {
            cout << 10 << "\n";
            continue;
        }

        if (N == 1)
        {
            cout << 1 << "\n";
            continue;
        }

        string m_digits = "";
        for (int d = 9; d >= 2; d--)
        {
            while (N % d == 0)
            {
                N = N / d;
                cout << "d: " << d << endl;
                m_digits += to_string(d);
            }
        }

        if (N > 1)
        {
            cout << -1 << "\n";
        }
        else
        {
            sort(m_digits.begin(), m_digits.end());
            cout << m_digits << "\n";
        }
    }

    return 0;
}