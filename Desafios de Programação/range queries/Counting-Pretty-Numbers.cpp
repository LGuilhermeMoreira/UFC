#include <bits/stdc++.h>
using namespace std;
typedef long long ll;

int main()
{
    ios::sync_with_stdio(0);
    cin.tie(0);

    int T;
    cin >> T;
    while (T--)
    {
        int L, R;
        cin >> L >> R;
        int count = 0;
        for (int i = L; i <= R; i++)
        {
            int last = i % 10;
            if (last == 2 || last == 3 || last == 9)
                count++;
        }
        cout << count << "\n";
    }

    return 0;
}
