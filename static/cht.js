window.onload = function() {
  var ctx = document.getElementById("myChart").getContext("2d");
  var chart = new Chart(ctx, {
    type: "line",
    data: {
      // label: ["Android", "iOS", "Windows"],
      datasets: [
        {
          label: "Android",
          type: "line",
          data: [],
          borderColor: "rgba(245, 146, 146, 0.9)",
          backgroundColor: "rgba(245, 146, 146, 0.9)",
          // backgroundColor: "rgba(218, 216, 216, 0)",
          fill: false
        },
        {
          label: "iOS",
          type: "line",
          data: [],
          borderColor: "rgba(146, 201, 245, 0.9)",
          backgroundColor: "rgba(146, 201, 245, 0.9)",
          // backgroundColor: "rgba(218, 216, 216, 0)",
          fill: false
        },
        {
          label: "Windows",
          type: "line",
          data: [],
          borderColor: "rgba(99, 238, 106, 0.9)",
          backgroundColor: "rgba(99, 238, 106, 0.9)",
          // backgroundColor: "rgba(218, 216, 216, 0)",
          fill: false
        }
      ]
    },
    options: {
      scales: {
        xAxes: [
          {
            type: "realtime",
            scaleLabel: {
              display: true,
              labelString: "Time",
              // fontColor: "white",
              fontSize: 20
            },
            gridLines: {
              // color: "rgba(255, 255, 255, .3)"
            },
            ticks: {
              // fontColor: "white",
              fontSize: 14
            }
          }
        ],
        yAxes: [
          {
            scaleLabel: {
              display: true,
              labelString: "Traffic",
              // fontColor: "white",
              fontSize: 20
            },
            gridLines: {
              // color: "rgba(255, 255, 255, .3)",
              // zeroLineColor: "rgba(255, 255, 255, .3)"
            },
            ticks: {
              // fontColor: "white",
              fontSize: 14
            }
          }
        ]
      },
      legend: {
        display: true,
        position: "right",
        lables: {
          boxWidth: 40,
          padding: 10,
          fontSize: 25
        }
      },
      plugins: {
        streaming: {
          duration: 50000,
          refresh: 1500,
          delay: 1000,
          frameRate: 30,
          pause: false,

          onRefresh: function(chart) {
            // let dataNum = chart.data.datasets.length;
            // for (let i = 0; i < dataNum; i++) {
            //   chart.data.datasets[i].data.push({
            //     x: Date.now(),
            //     y: Math.random() * 100
            //   });
            // }

            var i = 0;
            for (let os in counters) {
              console.log("os: " + os + ", counter: " + counters[os]);
              chart.data.datasets[i].data.push({
                x: Date.now(),
                y: counters[os]
              });
              i++;
            }
          }
        }
      }
    }
  });
};

var socket = new WebSocket("ws://localhost:8080/send");
var num = 0;
var data;
var jsonData;
var os;
var counters = {
  Android: 0,
  iOS: 0,
  Windows: 0
};

socket.addEventListener("open", e => {
  console.log("websocket connected");
});

socket.addEventListener("message", e => {
  // num = parseInt(e.data, 10);
  // console.log(typeof parseInt(e.data, 10)); // number
  // console.log(e);

  // data = e.data.replace(/\\"/g, '"');
  // console.log(e.data);

  // jsonData = JSON.parse(e.data);
  data = e.data;
  jsonData = eval(JSON.parse(data)); // string to object
  console.log(jsonData);

  console.log(jsonData.length);
  console.log(typeof jsonData);

  for (let i in jsonData) {
    // console.log(os);
    // console.log(os + ": " + data[os].Counter);
    // console.log("i: " + i);
    os = jsonData[i].os;
    counters[os] = jsonData[i].counter;
    // console.log("os: " + os + ", counter: " + counters[os]);
  }
});

socket.addEventListener("close", () => {
  console.log("websocket closed");
  for (var key in counters) {
    counters[key] = 0;
  }
});

socket.addEventListener("error", e => {
  console.log("websocket error : ", e);
});
