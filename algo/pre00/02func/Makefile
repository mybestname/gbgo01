GCC_FALG=-g -Wall -std=c++11
I=func.cpp
O=func

build:
	g++ $(GCC_FALG) $(I) -o $(O)

run: build
	./$(O)

dump_1:build
	objdump -S $(O)

dump_2:build
	otool -tv $(O)

dump_3:build
	go tool objdump -gnu -s change1 $(O)
	go tool objdump -gnu -s change2 $(O)
	go tool objdump -gnu -s change3 $(O)

clean_all: clean
clean:
	rm -rf $(O) $(O).dSYM
