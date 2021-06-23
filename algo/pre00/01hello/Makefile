GCC_FALG=-g -Wall -std=c++11
I=helloworld.cpp
O=hello

build1:
	g++ $(GCC_FALG) $(I) -o $(O)
build2:
	gcc $(GCC_FALG) -lstdc++ $(I) -o $(O)2

run1: build1
	./$(O)
run2: build2
	./$(O)2

clean_all: clean_1 clean_2
clean_1:
	rm -rf $(O) $(O).dSYM
clean_2:
	rm -rf $(O)2 $(O)2.dSYM