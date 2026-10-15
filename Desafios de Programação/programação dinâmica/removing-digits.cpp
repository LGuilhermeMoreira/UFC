#include <iostream>
#include <algorithm>

using namespace std;

typedef long long ll;

int main()
{
    ios::sync_with_stdio(0);
    cin.tie(0);

    ll n;
    cin >> n;
    ll cont = 0;

    if(n == 0) {
        cout << "0\n";
        return 0;
    }
    if(n < 10) {
        cout << "1\n";
        return 0;
    }

    while (n > 0)
    {
        ll temp = n;
        int maior_digito = 0;

        while (temp > 0)
        {
            int digito = temp % 10;
            if (digito > maior_digito)
            {
                maior_digito = digito;
            }
            temp /= 10;
        }

        n -= maior_digito;
        cont++;
    }

    cout << cont << "\n";

    return 0;
}