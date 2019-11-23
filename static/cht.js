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
          refresh: 1000,
          delay: 1000,
          frameRate: 30,
          pause: false,

          onRefresh: function(chart) {
            let dataNum = chart.data.datasets.length;
            // for (let i = 0; i < dataNum; i++) {
            //   chart.data.datasets[i].data.push({
            //     x: Date.now(),
            //     y: Math.random() * 100
            //   });
            // }
            console.log(num);
            chart.data.datasets[0].data.push({
              x: Date.now(),
              y: num
            });
          }
        }
      }
    }
  });
};

var socket = new WebSocket("ws://localhost:8080/send");
var num = 0;
var data = "";
// var insertData = document.getElementById("insertData");

socket.addEventListener("open", e => {
  console.log("websocket connected");
});

socket.addEventListener("message", e => {
  num = parseInt(e.data, 10);
  // console.log(typeof parseInt(e.data, 10)); // number
  // console.log(e);
  data = e.data.replace(/\\"/g, '"');
  console.log("data is " + data);
  // data = JSON.parse(e);
  // console.log(data);
  // var p = document.createElement("p");
  // p.innerHTML = e.data;
  // insertData.appendChild(p);
});

socket.addEventListener("close", () => {
  console.log("websocket closed");
  num = 0;
});

socket.addEventListener("error", e => {
  console.log("websocket error : ", e);
});
