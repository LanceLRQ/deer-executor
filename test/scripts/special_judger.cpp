/*
@resources: BNUZCPC HOT BODY A
@date: 2017-5-10
@author: QuanQqqqq
@algorithm: check code
*/
#include <iostream>
//#include <bits/stdc++.h>

using namespace std;

string ans[] = {"100","1923","1902","5000","850","40","2003","2007"};

int main(int argc, char *argv[]){
	int total = 0;
	int p = 0;
	string str;
	freopen(argv[3],"r",stdin);
	while(cin >> str){
		p++;
		for(int i = 0;i < 8;i++){
			if(ans[i] == str){
				total++;
			}
		}
	}
	if(total == 1 && p == 1){
		return 0;
	} else {
		return 4;
	}
}