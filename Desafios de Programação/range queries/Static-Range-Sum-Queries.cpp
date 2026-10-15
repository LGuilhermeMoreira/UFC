#include<bits/stdc++.h>
using namespace std;
typedef long long ll;



int main(){
    ios::sync_with_stdio(0);
    cin.tie(0);
    
    int n;
    cin >> n;
    int q;
    cin >> q;

    vector<int> arr(n);
    vector<ll> sum(n);
    
    for(int i = 0; i < n; i++){
        cin >> arr[i];
    }

    sum[0] = arr[0];
    for(int i =1; i < n; i++){
        sum[i] = sum[i-1] + arr[i];
    }

    while (q--)
    {
        int l,r;
        cin >> l >> r;
        int rr,rl;
        rr = r-1;
        if(l == 1){
            cout << sum[rr] << endl;
        }else{
            rl = l-2;
            cout << sum[rr] - sum[rl] << endl;
        }
    }
    

    return 0;
}