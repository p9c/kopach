# kopach

Miner for Parallelcoin

Kopach is a miner for Parallelcoin implementing CPU and OpenCL miners

Algorithms are detected via the block template provided by the server and if more than one algorithm is detected, it automatically selects the algorithm with the lowest difficulty, based on an automatic benchmark on first run that writes the benchmark figures for each algorithm into the configuration file.
