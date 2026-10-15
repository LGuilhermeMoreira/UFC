#include<bits/stdc++.h>
using namespace std;
typedef long long ll;


int ps[1001][1001];

int main(){
    ios::sync_with_stdio(0);
    cin.tie(0);

    int n, q;
    cin >> n >> q;

    
    for (int y = 1; y <= n; y++) {
        for (int x = 1; x <= n; x++) {
            char c;
            cin >> c;
            
            int tree = (c == '*'); 

            ps[y][x] = tree + ps[y-1][x] + ps[y][x-1] - ps[y-1][x-1];
        }
    }

    for (int i = 0; i < q; i++) {
        int y1, x1, y2, x2;
        cin >> y1 >> x1 >> y2 >> x2;

        int ans = ps[y2][x2] - ps[y1-1][x2] - ps[y2][x1-1] + ps[y1-1][x1-1];
        
        cout << ans << "\n";
    }

    return 0;
}