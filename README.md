# Red1s 
An extremely simplified Redis server implementation that supports the following commands:

* GET
* DEL 
* SET

## Tasks
- [ ] Implement TCP listener to handle incoming client requests.
- [ ] Create parser / decoder to process incoming commands.
- [ ] Create thread-safe key-value store to hold data.
- [ ] Create response encoder to return Redis Wire protocol encoded responses.

### Nice to haves
- [ ] Data store abstraction. 
- [ ] Automatic key expiration.
- [ ] Persistence - save data to disk, read in contents during startup.
- [ ] Partitioning - would require a load balancer / router.

## Developer Notes
* Redis can store up to 512 MB per key. Obviously allocating this much data per request would be terribly inefficient, especially 
for requests that are only a few bytes in size. Reading incoming request data can be done using chunking where we read in say 1024 KB
at a time, writing the data to a byte buffer. The max data size should be configurable via env var.

* Automatic key expiration could be done by keeping track of which keys are flagged to automatically expire.
  * Should the expirey data be part of the k/v data structure or part of an expirey data structure?
  * What happens if we have 10k auto expirey keys? What would be the best way to efficiently determine what should be removed? A priority queue perhaps?
  * How frequently would we run the check?