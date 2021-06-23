#include<iostream>
using namespace std;

//值传递
void change1(int n){
    //cout<<"change1() by val addr="<<&n<<endl;            //显示的是拷贝的地址而不是源地址
    n++;
}

//引用传递
void change2(int& n){                                       //	pushq	%rbp
    //cout<<"change2() by ref addr="<<&n<<endl;             //	movq	%rsp, %rbp
    n++;                                                    //	movq	%rdi, -0x8(%rbp)
}                                                           //	movq	-0x8(%rbp), %rax
                                                            //	movl	(%rax), %ecx
                                                            //	addl	$0x1, %ecx
                                                            //	movl	%ecx, (%rax)
                                                            // 	popq	%rbp
                                                            // 	retq
                                                            // 	nopw	%cs:(%rax,%rax)
                                                            // 	nop

//指针传递                                                    //	pushq	%rbp
void change3(int* n){                                       //	movq	%rsp, %rbp
    //cout<<"chagne3() by ptr addr="<<n<<endl;              //	movq	%rdi, -0x8(%rbp)
    *n=*n+1;                                                //	movq	-0x8(%rbp), %rax
}                                                           //	movl	(%rax), %ecx
                                                            //	addl	$0x1, %ecx
                                                            //	movq	-0x8(%rbp), %rax  //多取一次地址
                                                            //	movl	%ecx, (%rax)
                                                            //	popq	%rbp
                                                            //	retq
                                                            //	nopl	(%rax)

int main(){                                                 //  pushq	%rbp
    int n=10;                                               //  movq	%rsp, %rbp
    //cout<<"      main() val addr="<<&n<<endl;             //  subq	$0x10, %rsp
    //cout<<"         main() n="<<n<<endl;                  //  movl	$0x0, -0x4(%rbp)
    change1(n);                                             //  movl	$0xa, -0x8(%rbp)
    //cout<<"after change1() n="<<n<<endl;                  //  movl	-0x8(%rbp), %edi
    change2(n);                                          //  callq	0x100003f20
    //cout<<"after change2() n="<<n<<endl;                  //  leaq	-0x8(%rbp), %rdi   取地址入rdi
    change3(&n);                                            //  callq	0x100003f40
    // after chang3 n = 12                                  //  leaq	-0x8(%rbp), %rdi   取地址入rdi
    //cout<<"after change3() n="<<n<<endl;                  //  callq	0x100003f60        调用是完全一样的
    return 0;                                               //  xorl	%eax, %eax
}                                                           //  addq	$0x10, %rsp
                                                            //  popq	%rbp
                                                            //  retq


// 在汇编的层次上，除了细微的差别外，传地址和传ref的实现没有实质的区别。
// 具体的区别在更高级的语言抽象层面。