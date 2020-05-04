package main

http.Handle("/users", usersHandler)
log.Fatal(http.ListenAndServe(":8080", nil))