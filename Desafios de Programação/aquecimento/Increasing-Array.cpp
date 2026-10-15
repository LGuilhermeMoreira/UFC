#include <bits/stdc++.h>
using namespace std;

int main()
{
    ios::sync_with_stdio(0);
    cin.tie(0);

    long long n;
    cin >> n;
    vector<long long> arr(n);

    // lendo o array
    for (int i = 0; i < n; i++)
    {
        int v = 0;
        cin >> v;
        arr[i] = v;
    }
    long max = 0;
    long count = 0;
    for (int i = 0; i < n; i++)
    {
        if (arr[i] >= max)
        {
            max = arr[i];
        }
        else
        {
            count += max - arr[i];
        }
    }

    cout << count << endl;
}