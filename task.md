Add logic to save a list of all processed files into ./data/files.txt
The file should contain one relative file path per line, with LF line ending, relative from the data folder.
The purpose of the file is to reconstruct the full dump.sql file after it gets split into parts.
Since we do not save order of parts yet, we need this files.txt to restore the original order.
