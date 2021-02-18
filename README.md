# Red1s 
An extremely simplified Redis server implementation that supports the following commands:

* GET
* DEL 
* SET

---

## Tasks
- [x] Implement TCP listener to handle incoming client requests.
- [x] Create parser / decoder to process incoming commands.
- [x] Create thread-safe key-value store to hold data.
- [x] Return Redis Wire protocol encoded responses.

### Nice to haves
- [ ] List command.
- [x] Data store abstraction. 
- [ ] Automatic key expiration.
- [ ] Persistence - save data to disk, read in contents during startup.
- [ ] Partitioning - would require a load balancer / router.

## Developer Notes
* Redis can store up to 512 MB per key. Obviously allocating this much data per request would be terribly inefficient, especially 
for requests that are only a few bytes in size. Reading incoming request data can be done using chunking where we read in say 1024 KB
at a time, writing the data to a byte buffer.
  * Note: I ended up opting for a basic solution for now due to issues reading CRLF requests correctly. Would like to re-address this down the road.

* Manually parsing the command bytes, but will probably refactor to use simple regex checks.

## Resources
* https://redis.io/topics/protocol

---

## Notes for the reviewer
There's a tremendous amount of stuff I would like to, and would have liked to do for this project, such as:
  * Passing a context to cancel requests.
  * Redo how I did the parsing entirely (Regex would have been simpler / cleaner).
  * Adding a graceful shutdown by monitoring sigint and having the network later update a channel whenever a request as in-flight.
  * TCP code is not great. I got hung up trying to get bufio/bytes package, etc. working with conn.Read. I suspect the issue with the CRLF, 
    but I decided to move after spending a bit too long on it. However, this left me having to use a pre-allocated byte slice of 4096 which is
    of course, much less than ideal :(
  * Package layout isn't great. Naming is hard. This would definitely be part of a refactor.
  * The storage system uses a standadrd `RWMutex` vs a `sync.Map` since `sync.Map` is more optimized for high reads and fewer rights.
  * I have not written as many tests and or bench marked everything as well as I'd have liked. I also haven't ran the race detector.

---

## Author 
Jason Welch (jasonw.developer@gmail.com)
