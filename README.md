# go-git-stats

Utility to generate statistics from git repositories.

## Disclaimer

This is a work in progress. It is not ready for production use and currently only availbale as standalone binary.

## Installation

Download Repository and install binary:

```bash
$ make install
```

## Usage

Currently `go-git-stats` only prints out the number of lines of code (LOC) for each author in the repository.

```bash
$ go-git-stats
-----------------------+-------+------------+
| AUTHOR                |   LOC | PERCENTAGE |
+-----------------------+-------+------------+
| <pbertram@aerosys.io> |   223 | 100.00%    |
+-----------------------+-------+------------+
|                       | TOTAL | 223        |
+-----------------------+-------+------------+
```
