#include <iostream>
#include <vector>
#include <string>
#include <algorithm>

using namespace std;

int main()
{
    int a, b, c;
    cin >> a >> b >> c;

    int m;
    cin >> m;

    vector<long long> usb_costs;
    vector<long long> ps2_costs;

    for (int i = 0; i < m; ++i)
    {
        long long cost;
        string type;
        cin >> cost >> type;
        if (type == "USB")
        {
            usb_costs.push_back(cost);
        }
        else
        {
            ps2_costs.push_back(cost);
        }
    }

    sort(usb_costs.begin(), usb_costs.end());
    sort(ps2_costs.begin(), ps2_costs.end());

    long long total_equipped = 0;
    long long total_cost = 0;

    int usb_ptr = 0;
    int ps2_ptr = 0;

    for (int i = 0; i < a; ++i)
    {
        if (usb_ptr < usb_costs.size())
        {
            total_cost += usb_costs[usb_ptr];
            total_equipped++;
            usb_ptr++;
        }
    }

    for (int i = 0; i < b; ++i)
    {
        if (ps2_ptr < ps2_costs.size())
        {
            total_cost += ps2_costs[ps2_ptr];
            total_equipped++;
            ps2_ptr++;
        }
    }

    vector<long long> remaining_mouses;
    while (usb_ptr < usb_costs.size())
    {
        remaining_mouses.push_back(usb_costs[usb_ptr]);
        usb_ptr++;
    }
    while (ps2_ptr < ps2_costs.size())
    {
        remaining_mouses.push_back(ps2_costs[ps2_ptr]);
        ps2_ptr++;
    }

    sort(remaining_mouses.begin(), remaining_mouses.end());

    for (int i = 0; i < c; ++i)
    {
        if (i < remaining_mouses.size())
        {
            total_cost += remaining_mouses[i];
            total_equipped++;
        }
    }

    cout << total_equipped << " " << total_cost << endl;

    return 0;
}