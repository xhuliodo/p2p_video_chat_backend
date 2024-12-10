import ws from "k6/ws";
import http from "k6/http";

import { sleep, check } from "k6";
import { uuidv4 } from "https://jslib.k6.io/k6-utils/1.4.0/index.js";

export const options = {
  stages: [
    { duration: "10s", target: 10 },
    { duration: "30s", target: 50 },
    { duration: "1m", target: 250 },
    { duration: "2m", target: 500 },
    { duration: "4m", target: 1000 }, // just slowly ramp-up to a HUGE load
  ],
};

// The function that defines VU logic.
export default function () {
  const res = http.get("http://localhost:3030/call/id");
  if (!res.body) {
    return;
  }
  const callId = res.body;
  // Simulate a participant joining the current room
  const userId = uuidv4(); // Generate a unique user ID for each participant

  // if enables, this flag tests out the messages logic, if left disabled, the test
  // leaves the connections on throughout the whole load testing. so it depends on what you
  // want to test
  const messageLogic = false;

  participant(callId, userId, messageLogic);
  sleep(1);
}

const sampleOffer = `v=0
o=- 46117316 2 IN IP4 127.0.0.1
s=-
t=0 0
a=group:BUNDLE 0
a=extmap-allow-mixed
a=msid-semantic: WMS
m=audio 9 UDP/TLS/RTP/SAVPF 111 103 104 9 0 8 106 105 13 110 112 113 126
c=IN IP4 0.0.0.0
a=rtcp:9 IN IP4 0.0.0.0
a=ice-ufrag:EYbn
a=ice-pwd:1F9tWQ0tBngv1YQ7JtxK2jKe
a=ice-options:trickle
a=fingerprint:sha-256 67:88:74:16:95:0F:A3:18:FA:56:D9:44:BD:62:62:6F:11:64:26:85:3A:2C:BC:5D:A9:BD:4D:EB:87:C1:15:73
a=setup:actpass
a=mid:0
a=extmap:1 urn:ietf:params:rtp-hdrext:ssrc-audio-level
a=extmap:2 http://www.webrtc.org/experiments/rtp-hdrext/abs-send-time
a=extmap:3 urn:ietf:params:rtp-hdrext:sdes:mid
a=sendrecv
a=rtcp-mux
a=rtpmap:111 opus/48000/2
a=rtcp-fb:111 transport-cc
a=fmtp:111 minptime=10;useinbandfec=1
a=rtpmap:103 ISAC/16000
a=rtpmap:104 ISAC/32000
a=rtpmap:9 G722/8000
a=rtpmap:0 PCMU/8000
a=rtpmap:8 PCMA/8000
a=rtpmap:106 CN/32000
a=rtpmap:105 CN/16000
a=rtpmap:13 CN/8000
a=rtpmap:110 telephone-event/48000
a=rtpmap:112 telephone-event/32000
a=rtpmap:113 telephone-event/16000
a=rtpmap:126 telephone-event/8000
a=ssrc:292201749 cname:N9bwn0mI+IReSRy8
a=ssrc:292201749 msid:stream1 audio1
a=ssrc:292201749 mslabel:stream1
a=ssrc:292201749 label:audio1
`;
const sampleAnswer = `v=0
o=- 46117316 2 IN IP4 127.0.0.1
s=-
t=0 0
a=group:BUNDLE 0
a=extmap-allow-mixed
a=msid-semantic: WMS
m=audio 9 UDP/TLS/RTP/SAVPF 111
c=IN IP4 0.0.0.0
a=rtcp:9 IN IP4 0.0.0.0
a=ice-ufrag:7D7b
a=ice-pwd:XcNvi0bD69N3yL0UNM8y2qVQ
a=fingerprint:sha-256 8A:98:1A:D9:44:9D:72:F1:AE:12:88:21:44:9A:17:AE:24:C3:DE:E9:FE:8A:FA:1E:43:96:3F:65:E8:3D:46
a=setup:active
a=mid:0
a=recvonly
a=rtcp-mux
a=rtpmap:111 opus/48000/2
a=fmtp:111 minptime=10;useinbandfec=1
`;
const sampleIceCandidate = `candidate:842163049 1 udp 1677729535 192.168.1.2 63068 typ srflx raddr 0.0.0.0 rport 0 generation 0 ufrag EYbn network-id 1`;

async function participant(callId, userId, messageLogic) {
  const steps = {
    new_participant: false,
    offer: false,
    answer: false,
    ice_candidate: false,
    participant_left: false,
  };
  // establish connection with the ws server
  const url = "wss://localhost:8080/calls/" + callId;
  const params = { userId, callId };
  const res = ws.connect(url, params, function (socket) {
    socket.on("open", () => {
      // console.log(`[Participant ${userId} Joining call ${callId}`);
      // immediately send an event that a new participant has been added to the call
      const newParticipantEvent = newNewParticipantEvent(userId);
      socket.send(newParticipantEvent);
    });

    socket.on("message", (data) => {
      const event = JSON.parse(data);
      switch (event.type) {
        case "new_participant": {
          // console.log("got newParticipant");
          const eventData = event.payload;
          // send a mock offer
          const offerEvent = newOfferEvent(
            sampleOffer,
            eventData.participantId
          );
          socket.send(offerEvent);

          // send mock x amount of ice candidates an a random intervals
          simulateIceCandidatesGeneration(socket, 3, eventData.participantId);

          // set step as done
          steps["new_participant"] = true;
          break;
        }
        case "offer": {
          // console.log("got offer");
          const eventData = event.payload;
          // send mock answer back
          const answerEvent = newAnswerEvent(sampleAnswer, eventData.from);
          socket.send(answerEvent);

          // send mock x amount of ice candidates an a random intervals
          simulateIceCandidatesGeneration(socket, 3, eventData.from);

          // set step as done
          steps["offer"] = true;
          break;
        }
        case "answer": {
          // console.log("got answer");
          // set step as done
          steps["answer"] = true;
          break;
        }
        case "ice_candidate": {
          // do nothing
          // console.log("got ice_candidate");
          // set step as done
          steps["ice_candidate"] = true;

          // consider this step as the last one, and close connection afterward
          if (messageLogic) {
            socket.close();
          }
          break;
        }
        case "participant_left": {
          // do nothing
          // console.log("got participantLeft");
          // set step as done
          steps["participant_left"] = true;
          break;
        }
      }
    });
    socket.on("close", () => console.log("disconnected"));
    socket.on("error", function (e) {
      if (e.error() != "websocket: close sent") {
        console.log("An unexpected error occured: ", e.error());
      }
    });
  });

  check(res, { "status is 101": (r) => r && r.status === 101 });

  if (messageLogic) {
    check(steps, {
      "message logic": (s) => {
        // these are the events from the offerer side
        if (s.new_participant && s.answer && s.ice_candidate) {
          return true;
        }
        // these are from the reciever
        if (s.offer && s.ice_candidate) {
          return true;
        }

        return false;
      },
    });
  }
}

async function simulateIceCandidatesGeneration(socket, howMany, to) {
  for (let i = 0; i < howMany; i++) {
    const iceCandidateEvent = newIceCandidateEvent(sampleIceCandidate, to);
    socket.send(iceCandidateEvent);
    // console.log("sent ice candidate");
    sleep(getRandomInt(3));
  }
}

function getRandomInt(max) {
  return Math.floor(Math.random() * max);
}

function newNewParticipantEvent(userId) {
  const e = {
    type: "new_participant",
    payload: {
      userId,
    },
  };

  return JSON.stringify(e);
}

// const participantLeftEvent = () => {
//   const e = {
//     type: "participant_left",
//     payload: {},
//   };

//   return JSON.stringify(e);
// };

function newOfferEvent(offer, to) {
  const e = {
    type: "offer",
    payload: {
      offer,
      to,
    },
  };
  return JSON.stringify(e);
}

function newIceCandidateEvent(iceCandidate, to) {
  const e = {
    type: "ice_candidate",
    payload: {
      iceCandidate,
      to,
    },
  };

  return JSON.stringify(e);
}

function newAnswerEvent(answer, to) {
  const e = {
    type: "answer",
    payload: {
      answer,
      to,
    },
  };

  return JSON.stringify(e);
}
