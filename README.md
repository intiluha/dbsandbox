# dbsandbox

Allows you to create multiple readers and writers that manipulate one database concurrently. Readers delete rows as they read them and it is guaranteed that every row will be read only once.

Set MYSQL_USER and MYSQL_PASS before use.
```sh
$ export MYSQL_USER=username
$ export MYSQL_PASS=password
```
