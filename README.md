
run app:
<br/>`go run main.go`

run stress test client:
<br/>`go run client.go`

CTRL-C to exit both of them

---------------------

to improve:
1. orgnize http hanlders code
2. `setStatus()` of `JobRepo` doesnt atomically check if job exists in repo before changing status in sync map

tradeoffs:
1. in `JobCreate(newJob JobCreateDto)`, job can be created but not scheduled if the channel is full, it will be still deleted 2s later but its a bad practice when storing into databases.
2. if there are more requests than workers can handle, they will be stored in buffered channel, if channel buffer exceeds `chanBufferSize` it will return response 503