#include <bits/stdc++.h>
using namespace std;

vector<int> solve(vector<vector<int>>& adj) {
    int n = adj.size();
    vector<int> indegree(n, 0);
    queue<int> q;
    vector<int> list;

    for (int i = 0; i < n; i++) {
        for (int next : adj[i])
            indegree[next]++;
    }

   
    for (int i = 0; i < n; i++)
        if (indegree[i] == 0)
            q.push(i);


    while (!q.empty()) {
        int top = q.front();
        q.pop();
        list.push_back(top);
        for (int next : adj[top]) {
            indegree[next]--;
            if (indegree[next] == 0)
                q.push(next);
        }
    }

    return list;
}

void addEdge(vector<vector<int>>& adj, int u, int v) {
    adj[u].push_back(v);
}

int main() {
    int n,m;
    cin >> n >> m;
    vector<vector<int>> adj(n);
    
    for (int i = 0; i < m; i++)
    {
        int u,v;
        cin >> u >> v;
        u--;
        v--;
        addEdge(adj, u, v);
    }
    
    vector<int> res = solve(adj);
    
    if(res.size() != n) {
        cout << "IMPOSSIBLE" << endl;
        return 0;
    }
    
    for (int vertex : res)
        cout << vertex + 1 << " ";
    cout << endl;

    return 0;
}
