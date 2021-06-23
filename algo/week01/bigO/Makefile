CC=gcc
CFLAGS=-g -Wall -std=c99

all:calc fib

calc: calc.c
	$(CC) $(CFLAGS) -o $@ $<
fib: fib.c
	$(CC) $(CFLAGS) -o $@ $<

clean_all: clean
clean:
	rm -rf calc fib *.dSYM