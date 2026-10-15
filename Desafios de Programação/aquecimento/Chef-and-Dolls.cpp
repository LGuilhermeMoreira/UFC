#include <bits/stdc++.h>
using namespace std;

int main(int argc, char const *argv[])
{
    long long n = 0;
    cin >> n;

    for (long long i = 0; i < n; i++)
    {
        long long a = 0;
        cin >> a;
        map<long long, long long> mp;

        for (long long i = 0; i < a; i++)
        {
            long long c = 0;
            cin >> c;
            mp[i]++;
        }

        for (auto x : mp)
        {
            if (!(x.second % 2 == 0))
            {
                cout << x.first << endl;
            }
        }
    }

    return 0;
}
