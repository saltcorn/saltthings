import uasyncio
import random
import time
#import inspect

random.seed(time.ticks_ms())

def randChar():
  return 'abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789-_'[random.getrandbits(6)]

def randStr(N=10):
	return ''.join([randChar() for i in range(N)])

print(randStr(20))