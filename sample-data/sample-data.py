import requests
import csv
import sys
import sqlite3

csv.field_size_limit(sys.maxsize)

conn = sqlite3.connect('/Users/kristiancharbonneau/go/src/github.com/kristiancharb/SimplSearch/docs.db')

create_table = 'CREATE TABLE IF NOT EXISTS docs (id INTEGER PRIMARY KEY, index_name VARCHAR(255), title TEXT, contents TEXT);'
cursor = conn.cursor()
cursor.execute(create_table)

docs = []
j = 0
for i in range(1, 4):
    with open(f'articles{i}.csv', newline='') as csvfile:
        reader = csv.reader(csvfile, delimiter=',')
        for row in reader:
            if (len(row) < 10): 
                continue
            title = row[2].replace('"', '')
            contents = row[9].replace('"', '')
            docs.append(('test', title, contents))
            if j % 10 == 0:
                print(j, 'documents processed')
            j += 1

insert_docs = "INSERT INTO docs (index_name, title, contents) VALUES (?, ?, ?)"

cursor.executemany(insert_docs, docs)
print(cursor.rowcount, 'rows inserted')
conn.commit()

conn.close()






# host = 'http://localhost:8000'
# index_name = 'test'

# r = requests.post(f'{host}/index', json={ 'name': 'test' })
# if r.status_code != requests.codes.ok:
#     print('Error occured with request:')
#     print(r)
#     exit(1)

# j = 0
# for i in range(1, 4):
#     with open(f'articles{i}.csv', newline='') as csvfile:
#         reader = csv.reader(csvfile, delimiter=',')
#         for row in reader:
#             if (len(row) < 10): 
#                 continue
#             title = row[2].replace('"', '')
#             contents = row[9].replace('"', '')
#             r = requests.post(f'{host}/index/{index_name}', json={
#                 'title': title,
#                 'contents': contents
#             })
#             if r.status_code != requests.codes.ok:
#                 print('Error occured with request:')
#                 print(r)
#             if j % 20 == 0:
#                 print(j, 'requests sent')
#             j += 1
        


