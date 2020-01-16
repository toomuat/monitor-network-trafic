window.onload = function() {
  var ctx = document.getElementById("myChart").getContext("2d");
  var chart = new Chart(ctx, {
    type: "line",
    data: {
      // label: ["Android", "iOS", "Windows"],
      datasets: [
        // {
        //   label: "Android",
        //   type: "line",
        //   data: [],
        //   borderColor: "rgba(245, 146, 146, 0.9)",
        //   backgroundColor: "rgba(245, 146, 146, 0.9)",
        //   // backgroundColor: "rgba(218, 216, 216, 0)",
        //   fill: false
        // },
        {
          label: "iOS",
          type: "line",
          data: [],
          borderColor: "rgba(146, 201, 245, 0.9)",
          backgroundColor: "rgba(146, 201, 245, 0.9)",
          // backgroundColor: "rgba(218, 216, 216, 0)",
          fill: false
        }
        // {
        //   label: "Windows",
        //   type: "line",
        //   data: [],
        //   borderColor: "rgba(99, 238, 106, 0.9)",
        //   backgroundColor: "rgba(99, 238, 106, 0.9)",
        //   // backgroundColor: "rgba(218, 216, 216, 0)",
        //   fill: false
        // }
      ]
    },
    options: {
      scales: {
        xAxes: [
          {
            type: "realtime",
            scaleLabel: {
              display: true,
              labelString: "Time"
              // fontColor: "white",
              // fontSize: 20
            },
            gridLines: {
              // color: "rgba(255, 255, 255, .3)"
            },
            ticks: {
              // fontColor: "white",
              // fontSize: 14
            }
          }
        ],
        yAxes: [
          {
            scaleLabel: {
              display: true,
              labelString: "Traffic"
              // fontColor: "white",
              // fontSize: 20
            },
            gridLines: {
              // color: "rgba(255, 255, 255, .3)",
              // zeroLineColor: "rgba(255, 255, 255, .3)"
            },
            ticks: {
              // fontColor: "white",
              // fontSize: 14
            }
          }
        ]
      },
      legend: {
        display: true,
        position: "right",
        // position: "top",
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
            // let i = 0;
            // for (let os in counters) {
            //   // console.log("os: " + os + ", counter: " + counters[os]);
            //   chart.data.datasets[i].data.push({
            //     x: Date.now(),
            //     y: counters[os]
            //   });
            //   i++;
            // }

            chart.data.datasets[0].data.push({
              x: Date.now(),
              y: counters["iOS"]
            });
          }
        }
      }
    }
  });
};

// change ip address of the server depending on the environment
// var socket = new WebSocket("ws://localhost:8080/send");
// var socket = new WebSocket("ws://localhost/send");
var socket = new WebSocket("ws://192.168.11.195:8080/send");

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
  data = e.data;
  jsonData = eval(JSON.parse(data)); // string to object

  for (let i in jsonData) {
    os = jsonData[i].os;
    counters[os] = jsonData[i].counter;
  }
  // delete counters.Windows;
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
