# Load testing

**k6 load testing** for a webrtc signaling server with web sockets.

## Why

This load testing came to be for two reasons:

- Out of curiosity I wanted to test out just how many concurrent connections could the server take.
- Since this is just a small pet project, I was going to deploy it in a small VM (ignoring CPU, 500MB of RAM)

So the main question was, **how many connections with these constrains**.

## How

I chose to go with [k6](https://k6.io/open-source/). I had some prior experience using them, and they offered a ws lib.

In the load.js I have simulated a scenario in which a participant joins a room, and exchanges webrtc details with other participants, and just stays connected just further updates.

I had a small hiccup trying to coordinate the callId between the participants. Could not, for the life of me, get it working via k6. Whenever I tried setting it as a global variable, it would update for each new participant. Meaning each participant would be alone in their room.

In the end I just created a simple api whose only job is to serve a uuid to X amount of participants. Code located at: ./get_id.

So to wrap up, if you want to replicate this test:

- Run the main file in the directory **get_id**, this is the api to serve the roomId
- In the file /cmd/main.go, add the import _ "net/http/pprof" and uncomment the lines 23-50 and 67
- Build and run the signaling server
- Run the load.js file via k6

---

## Performance Metrics

To save up some time and also not saturate the get_id api (it's a rather simple service and did not want to optimize it too much), I decided for the steps to be as below:

```js
stages: [
    { duration: "10s", target: 10 },
    { duration: "30s", target: 50 },
    { duration: "1m", target: 250 },
    { duration: "2m", target: 500 },
    { duration: "4m", target: 1000 }, // just slowly ramp-up to a HUGE load
  ],
```

1000 concurrent users is a nice and round number to do some napkin math in the end.

### K6 Metrics

| **Category**            | **Metric**        | **Value**     | **Avg**  | **Min** | **Median** | **Max** | **P(90)** | **P(95)** |
| ----------------------- | ----------------- | ------------- | -------- | ------- | ---------- | ------- | --------- | --------- |
| **Data Throughput**     | Data Received     | 8.2 MB        | 17 kB/s  | -       | -          | -       | -         | -         |
|                         | Data Sent         | 6.3 MB        | 13 kB/s  | -       | -          | -       | -         | -         |
| **HTTP Requests**       | Blocked           | -             | 1.05ms   | 263µs   | 945µs      | 6.65ms  | 1.64ms    | 2.46ms    |
|                         | Connecting        | -             | 849.19µs | 204µs   | 721µs      | 6.41ms  | 1.36ms    | 1.99ms    |
|                         | Duration          | -             | 703.43µs | 154µs   | 652µs      | 4.1ms   | 985µs     | 1.49ms    |
|                         | Receiving         | -             | 124.53µs | 22µs    | 124µs      | 2.06ms  | 171µs     | 185µs     |
|                         | Sending           | -             | 131.54µs | 17µs    | 121µs      | 2.48ms  | 181.2µs   | 217.19µs  |
|                         | Waiting           | -             | 447.34µs | 101µs   | 381µs      | 3.72ms  | 641.2µs   | 935.69µs  |
|                         | Failed Requests   | 0.00% (0/999) | -        | -       | -          | -       | -         | -         |
|                         | Total Requests    | 999           | 2.04/s   | -       | -          | -       | -         | -         |
| **WebSocket Metrics**   | Connecting        | -             | 7.61ms   | 2.94ms  | 8.19ms     | 24.93ms | 8.93ms    | 9.18ms    |
|                         | Messages Received | 13,492        | 27.53/s  | -       | -          | -       | -         | -         |
|                         | Messages Sent     | 12,975        | 26.48/s  | -       | -          | -       | -         | -         |
|                         | Sessions          | 999           | 2.04/s   | -       | -          | -       | -         | -         |
| **Virtual Users (VUs)** | Active VUs        | 999           | -        | 1       | -          | 999     | -         | -         |
|                         | Max VUs           | 1000          | -        | 1000    | -          | 1000    | -         | -         |

## System Metrics

| **Metric**              | **Value** |
| ----------------------- | --------- |
| **Allocated Memory**    | 15.82 MB  |
| **Total Allocated**     | 101.09 MB |
| **System Memory**       | 42.64 MB  |
| **Garbage Collections** | 27        |

---

## Conclusion

From how I understand these data and some napkin math:

If for 1000 concurrent connections, the system memory got up to 42.64 MB. With my constraints of 500MB the server should be able to handle

$\frac{500 \times 1000}{43} =$ ~11.6K concurrent connections.
