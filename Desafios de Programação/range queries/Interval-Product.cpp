#include <iostream>
#include <vector>
#include <string>

const int MAXN = 1e5 + 10;
int n, k;
int arr[MAXN];
int seg_tree[MAXN * 4];
const int NEUTRAL = 1;

int get_sign(int val) {
    if (val > 0) return 1;
    if (val < 0) return -1;
    return 0;
}

int calculate(int a, int b) {
    return a * b;
}


void build(int node, int left, int right) {
    if (left == right) {
        seg_tree[node] = arr[left];
    } else {
        int middle = left + (right - left) / 2;
        int left_node = 2 * node + 1;
        int right_node = 2 * node + 2;

        build(left_node, left, middle);
        build(right_node, middle + 1, right);

        seg_tree[node] = calculate(seg_tree[left_node], seg_tree[right_node]);
    }
}

void update(int node, int left, int right, int index, int value) {
    if (left == right) {
        seg_tree[node] = value;
    } else {
        int middle = left + (right - left) / 2;
        int left_node = 2 * node + 1;
        int right_node = 2 * node + 2;

        if (index <= middle) {
            update(left_node, left, middle, index, value);
        } else {
            update(right_node, middle + 1, right, index, value);
        }

        seg_tree[node] = calculate(seg_tree[left_node], seg_tree[right_node]);
    }
}

int query(int node, int left, int right, int x, int y) {
    if (x <= left && right <= y) {
        return seg_tree[node];
    }
    if (right < x || left > y) {
        return NEUTRAL;
    }

    int middle = left + (right - left) / 2;
    int left_node = 2 * node + 1;
    int right_node = 2 * node + 2;

    int answer_left = query(left_node, left, middle, x, y);
    int answer_right = query(right_node, middle + 1, right, x, y);

    return calculate(answer_left, answer_right);
}


int main() {
    std::ios_base::sync_with_stdio(false);
    std::cin.tie(NULL);

    while (std::cin >> n >> k) {
        for (int i = 0; i < n; i++) {
            int val;
            std::cin >> val;
            arr[i] = get_sign(val);
        }

        build(0, 0, n - 1);

        std::string result_string = "";
        for (int i = 0; i < k; i++) {
            char command;
            int arg1, arg2;
            std::cin >> command >> arg1 >> arg2;

            if (command == 'C') {
                int index = arg1 - 1; 
                int new_value_sign = get_sign(arg2);
                update(0, 0, n - 1, index, new_value_sign);
            } else if (command == 'P') {
                int start = arg1 - 1; 
                int end = arg2 - 1;   
                int product_sign = query(0, 0, n - 1, start, end);

                if (product_sign > 0) {
                    result_string += '+';
                } else if (product_sign < 0) {
                    result_string += '-';
                } else {
                    result_string += '0';
                }
            }
        }
        std::cout << result_string << std::endl;
    }

    return 0;
}
