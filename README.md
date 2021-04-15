# dbsandbox

Allows you to create multiple readers and writers that manipulate one database concurrently. Readers delete rows as they read them and it is guaranteed that every row will be read only once.
