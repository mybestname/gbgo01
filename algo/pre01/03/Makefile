GCC_FALG=-g -Wall -std=c++11
I=two_sum.cpp
O=twosum

build:
	g++ $(GCC_FALG) $(I) -o $(O)

run: build
	./$(O)

clean_all: clean
clean:
	rm -rf $(O) $(O).dSYM