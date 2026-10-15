#include<bits/stdc++.h>
using namespace std;
typedef long long ll;

int main(){
    ios::sync_with_stdio(0);
    cin.tie(0);

    int n, m;
    cin >> n >> m;

    vector<vector<pair<int, int>>> adj(n + 1);
    for (int i = 0; i < m; i++) {
        int a, b, c;
        cin >> a >> b >> c;
        adj[a].push_back({b, c});
    }

    priority_queue<pair<ll, pair<int, int>>, vector<pair<ll, pair<int, int>>>, greater<pair<ll, pair<int, int>>>> pq;
    vector<vector<ll>> dist(n + 1, vector<ll>(2, 1e18));

    dist[1][0] = 0;
    pq.push({0, {1, 0}});

    while (!pq.empty()) {
        ll d = pq.top().first;
        int u = pq.top().second.first;
        int f = pq.top().second.second;
        pq.pop();

        if (d > dist[u][f]) continue;

        for (auto e : adj[u]) {
            int v = e.first;
            int w = e.second;

            if (dist[u][f] + w < dist[v][f]) {
                dist[v][f] = dist[u][f] + w;
                pq.push({dist[v][f], {v, f}});
            }

            if (f == 0) {
                if (dist[u][0] + w / 2 < dist[v][1]) {
                    dist[v][1] = dist[u][0] + w / 2;
                    pq.push({dist[v][1], {v, 1}});
                }
            }
        }
    }

    cout << dist[n][1] << "\n";

    return 0;
}