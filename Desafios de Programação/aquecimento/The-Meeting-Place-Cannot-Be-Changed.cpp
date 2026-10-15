#include<bits/stdc++.h>
using namespace std;
typedef long long ll;

int n;
vector<int> x, v;

bool can_meet(double t) {
    double max_left_reach = -1e18; 
    double min_right_reach = 1e18;  

    for (int i = 0; i < n; ++i) {
        double left_bound = (double)x[i] - (double)v[i] * t;
        double right_bound = (double)x[i] + (double)v[i] * t;

        max_left_reach = max(max_left_reach, left_bound);
        min_right_reach = min(min_right_reach, right_bound);
    }

    return max_left_reach <= min_right_reach;
}

int main(){
    ios::sync_with_stdio(0);
    cin.tie(0);

    cin >> n;
    x.resize(n);
    v.resize(n);

    for(int i = 0; i < n; ++i) cin >> x[i];
    for(int i = 0; i < n; ++i) cin >> v[i];

    double low = 0.0;
    double high = 1e9 + 7;

    for(int iter = 0; iter < 100; ++iter) {
        double mid = low + (high - low) / 2.0;
        if (can_meet(mid)) {
            high = mid;
        } else {
            low = mid;
        }
    }

    cout << fixed << setprecision(10) << high << endl;

    return 0;
}


