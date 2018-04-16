# Description

**filegen** is random data generator tool that supports:
  * Generate tree of files with random data
  * Modify files with random data with controlling of modifications ranges

## Data generators

There are 3 data generators are supported now:
  * Crypto data generator
  * Pseudo random generator with seed support
  * Null blocks generator

# Usage

This sections describes how to use **filegen** tool to:
  * Generate new files
  * Modify existing files

## Generate new files

To generate new files filegen use **generate** (or **gen**) command. For example, the next command generates 5 folders with 10 file of 4096 bytes size  in each folder:
```
filegen gen -p /tmp/files -d 5 -f 10 -s 4K
```

There are options for **generate** command:
```
Common options:
  -p, --path                 Path to processing folder
  -q, --quiet                Quiet mode

Generate command options:
  -d, --dirs                 Directories count to generate
  -f, --files                Files count to generate
  -s, --size                 File size to generate. Size format: [\d{k,K,m,M,g,G}]

Generator options:
  -g, --generator            Type of generator to use
     crypto                  Crypto random data generator. Used by default.
     pseudo                  Pseudo random data generator
     null                    Null contains data generator
  --seed                     Initial seed for generated data. Can be used only with 'pseudo' generator
  --multiple-thread          Generate files with multiple threads. This mean than generated files can be fragmented
                             on disk. Can't be used with 'pseudo' generator
```

*-s, --size* option supports the following endings:
  * *k* - 10^3 bytes
  * *K* - 2^10 bytes
  * *m* - 10^6 bytes
  * *M* - 2^20 bytes
  * *g* - 10^9 bytes
  * *G* - 2^30 bytes

### Data generators

There are 3 supported data generators to create or modify files:

#### **crypto** generator

This generator is a cryptographically strong pseudo-random generator. On Linux, it uses getrandom(2) if available, /dev/urandom otherwise. On OpenBSD, generator uses getentropy(2). On other Unix-like systems, generator reads from /dev/urandom. On Windows systems, it uses the CryptGenRandom API. 

#### **pseudo** crypto generator
This generator is not cryptographically strong pseudo-random generator but it supports **seed** to regenerate the same random data. It used AES encryption. 

#### **null** blocks generator

This generator creates static blocks with nulls  

#### Generation data with multiple threads

This mode allow write files to disk in multiple threads. This means that writen files can be fragmented on disk. Each block is writen on disk in the thread that goroutines scheduler has selected. Option **--multiple-threads** is used for this. By default all files is writen in the single thread.

## Modify existing files

To modify existing files filegen use **change** (or **chg**) command. For example, the next command modify last 10% of file for 50% of existing files:
```
filegen chg -p /tmp/files --scale .5 -i 0,10% --once --reverse
```

There are options for **change** command:
```
Common options:
  -p, --path                 Path to processing folder
  -q, --quiet                Quiet mode

Change command options:
  --scale                    Percent of files cont to change. Range: [0;1]. By default is equal to 1
  -i, --interval             Interval to change file with. Format: ['data not to change', 'data to change',{'data not to change'}].
                             Data format: [\d{%, k,K,m,M,g,G}]. By default is [0,100%] and used until file ending
  --once                     Using of interval only once. Used only with -i, --interval option.
  --reverse                  Using interval from the file ending to begining. Used only with -i, --interval option.

Generator options:
  -g, --generator            Type of generator to use
     crypto                  Crypto random data generator. Used by default.
     pseudo                  Pseudo random data generator
     null                    Null contains data generator
  --seed                     Initial seed for generated data. Can be used only with 'pseudo' generator
  --multiple-thread          Generate files with multiple threads. This mean than generated files can be fragmented
                             on disk. Can't be used with 'pseudo' generator
```

### Intervals

Intervals is powerful tools to modify files. Interval is a triplet with optional last item: **(not to modify; modify; not to modify)**.
Each value in interval could be:
  * Absolute value in bytes. As *-s* parameter for *generate* command it supports endings *k*, *K*, *m*, *M*, *g*, *G*
  * Relative value in percents of file size. In this case **%** ending is used

To use intervals option **-i**, **--interval** is used. By default interval is (0, 100%).

There are special options for intervals:
  * **--once** used to apply interval once. In other case interval will be used until the end of file
  * **--reverse** used to apply interval from the end of file.

*Examples:*

Modify first 20% of file's data:
```
-i 0,20% --once
```

Modify each 20% of file's data with 1M gap from the end to the begining:
```
-i 0,20%,1M --reverse
```