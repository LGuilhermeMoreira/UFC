#include<bits/stdc++.h>
using namespace std;
typedef long long ll;

const int INF = 2e9;

struct SegTree {
    vector<int> tree;
    int n;
    bool isMinTree;

    SegTree(const vector<int>& arr, bool isMin) {
        n = arr.size();
        isMinTree = isMin;
        tree.resize(4 * n);
        build(arr, 1, 0, n - 1);
    }

    int combine(int a, int b) {
        if (isMinTree) return min(a, b);
        return max(a, b);
    }

    void build(const vector<int>& arr, int v, int tl, int tr) {
        if (tl == tr) {
            tree[v] = arr[tl];
        } else {
            int tm = (tl + tr) / 2;
            build(arr, v * 2, tl, tm);
            build(arr, v * 2 + 1, tm + 1, tr);
            tree[v] = combine(tree[v * 2], tree[v * 2 + 1]);
        }
    }

    int query(int v, int tl, int tr, int l, int r) {
        if (l > r) {
            return isMinTree ? INF : -INF;
        }
        if (l == tl && r == tr) {
            return tree[v];
        }
        int tm = (tl + tr) / 2;
        return combine(
            query(v * 2, tl, tm, l, min(r, tm)),
            query(v * 2 + 1, tm + 1, tr, max(l, tm + 1), r)
        );
    }

    int query(int l, int r) {
        return query(1, 0, n - 1, l, r);
    }
};

void solve() {
    int n;
    cin >> n;
    vector<int> a(n);
    for (int i = 0; i < n; ++i) {
        cin >> a[i];
    }

    SegTree st_max(a, false);
    SegTree st_min(a, true);

    for (int x = 1; x <= n - 2; ++x) {
        int k = st_max.query(0, x - 1);

        int j_lo = -1, j_hi = -1;
        int l, r, mid;

        l = x + 1, r = n - 1;
        while (l <= r) {
            mid = l + (r - l) / 2;
            int val = st_max.query(mid, n - 1);
            if (val == k) {
                j_lo = mid;
                r = mid - 1;
            } else if (val < k) {
                l = mid + 1;
            } else {
                r = mid - 1;
            }
        }

        if (j_lo == -1) continue;

        l = j_lo, r = n - 1;
        while (l <= r) {
            mid = l + (r - l) / 2;
            int val = st_max.query(mid, n - 1);
            if (val == k) {
                j_hi = mid;
                l = mid + 1;
            } else if (val < k) {
                l = mid + 1;
            } else {
                r = mid - 1;
            }
        }
        
        if (j_hi == -1) continue;

        int J_start = j_lo;
        int J_end = j_hi;

        if (st_min.query(x, J_start - 1) < k) continue;
        if (st_min.query(x, J_end - 1) > k) continue;

        int j_target = -1;
        l = J_start, r = J_end;
        while (l <= r) {
            mid = l + (r - l) / 2;
            int val = st_min.query(x, mid - 1);
            if (val <= k) {
                j_target = mid;
                r = mid - 1;
            } else {
                l = mid + 1;
            }
        }

        if (j_target != -1 && st_min.query(x, j_target - 1) == k) {
            int y = j_target - x;
            int z = n - j_target;
            cout << "YES\n";
            cout << x << " " << y << " " << z << "\n";
            return;
        }
    }

    cout << "NO\n";
}

int main(){
    ios::sync_with_stdio(0);
    cin.tie(0);

    int t;
    cin >> t;
    while (t--) {
        solve();
    }

    return 0;
}