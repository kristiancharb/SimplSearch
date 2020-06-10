import requests
import csv
import sys
import sqlite3

csv.field_size_limit(sys.maxsize)

project_path = f'{sys.argv[1]}/go/src/github.com/kristiancharb/SimplSearch'
conn = sqlite3.connect(f'{project_path}/docs.db')

create_table = 'CREATE TABLE IF NOT EXISTS docs (id INTEGER PRIMARY KEY, index_name VARCHAR(255), title TEXT, contents TEXT);'
cursor = conn.cursor()
cursor.execute(create_table)

docs = []
curr_title = ''
curr_contents = ''
insert_docs = "INSERT INTO docs (index_name, title, contents) VALUES (?, ?, ?)"
j = 0

with open(f'{project_path}/sample-data/simple-wiki.txt', newline='') as f:
    for line in f:
        if line == '\n':
            docs.append(('test', curr_title.strip(), curr_contents.strip()))
            curr_title = ''
            curr_contents = ''
        elif curr_title == '':
            curr_title = line
        else:
            curr_contents += line

        if j % 10 == 0:
            print(j, 'documents processed')
        if j % 1000 == 0:
            cursor.executemany(insert_docs, docs)
            docs = []
        j += 1

conn.commit()

conn.close()
