# synchronization-problems
### Classical (and not-so-classical) Synchronization Problems

All code tested on Ubuntu 16.04 LTS, using psutil (https://pypi.org/project/psutil/) for Python & gopsutil (https://github.com/shirou/gopsutil) for Go.

All code credits reproduced below, and included in the headers of their respective files.

- counter.c, counter.py, threading-cleanup.py: https://greenteapress.com/wp/

- producer-consumer.py: https://www.agiliq.com/blog/2013/10/producer-consumer-problem-in-python/ (modified)
- producer-consumer.go: http://www.golangpatterns.info/concurrency/producer-consumer (modified)

- dining-philosophers.py: https://rosettacode.org/wiki/Dining_philosophers#Python (modified)
- dining-philosophers.go: https://rosettacode.org/wiki/Dining_philosophers#Go (modified)

- dining-savages.py: algorithm from Allen B. Downey, *The Little Book of Semaphores*
- dining-savages.go: https://github.com/fsouza/lbos/blob/master/009-dining-savages.go (modified)
- dining-savages.c: https://github.com/pythonGuy/classic-threads/blob/master/DiningSavages.c (modified)

- modus-hall.py: algorithm from Allen B. Downey, *The Little Book of Semaphores*
- modus-hall.go: https://blog.ksub.org/bytes/post/modus-hall/modus-hall.go (modified)

- sushi-bar.py: algorithm from Allen B. Downey, *The Little Book of Semaphores*
- sushi-bar.go: https://blog.ksub.org/bytes/post/sushi-bar/sushi-bar.go (modified)
