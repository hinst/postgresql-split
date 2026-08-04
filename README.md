# postgresql-split

Split or rebuild PostgreSQL dump files.

## Usage

### Split (`default`)

```sh
postgresql-split <input.sql>
```

Splits a large PostgreSQL dump into smaller per-object SQL files under `data/`.  
The first positional argument is the input file path (required).

### Build (`-build`)

```sh
postgresql-split -build <source_dir> [output.sql]
```

Reconstructs a single dump from split files in `<source_dir>` by reading its `files.txt` manifest and concatenating all referenced SQL files.  
If no output filename is given, defaults to `dump.sql`.

## Arguments

| Argument     | Mode     | Description                                                |
| ------------ | -------- | ---------------------------------------------------------- |
| `-build`     | Build    | Path to directory containing split dump files              |
| `<input>`    | Split    | Path to the PostgreSQL dump file to split                  |
| `<output>`   | Build    | Optional output file (default: `dump.sql`)                 |
