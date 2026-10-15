#include <iostream>
using namespace std;

int main()
{
    ios::sync_with_stdio(0);
    cin.tie(0);
    long long int num = 0;
    cin >> num;
    while (true)
    {
        cout << num << " ";
        if (num == 1)
        {
            break;
        }
        if (num % 2 == 0)
        {
            num = num / 2;
        }
        else
        {
            num = num * 3 + 1;
        }
        // cout << num << " ";
    }

    return 0;
}