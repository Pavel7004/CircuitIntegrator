let getDocumentValue = (id) => {
    return document.getElementById(id).valueAsNumber;
};

let hostname = document.location.host
let socket = new WebSocket(`ws://${hostname}/gate`);
let plotter = document.getElementById('plotter');
let statusLabel = document.getElementById('status');

var dataX = [];
var dataY = [];

function buttonClick() {
    dataX = [];
    dataY = [];

    let msg = JSON.stringify({
        supply_voltage:      getDocumentValue('supplVolts'),
        capacity:            getDocumentValue('capacity'),
        resistance:          getDocumentValue('resistance'),
        stages_count:        getDocumentValue('stagesCount'),
        gap_trigger_voltage: getDocumentValue('gapVolts'),
        holding_voltage:     getDocumentValue('holdingVolts'),
        load_resistance:     getDocumentValue('loadResistance'),
        step:                getDocumentValue('step'),
        int_num:             getDocumentValue('intNum'),
    });
    socket.send(msg);
}
socket.onclose = function (event) {
  if (!event.wasClean) {
    statusLabel.innerText = 'Ошибка! Соединение разорвано';
  } else {
    statusLabel.innerText = 'Соединение прервано';
  }
};
socket.onmessage = (event) => {
    let msg = JSON.parse(event.data);

    if (msg.type === "point") {
        dataX = dataX.concat(msg.x)
        dataY = dataY.concat(msg.y)
    }
    if (msg.type === "end") {
        let data = {
            x: dataX,
            y: dataY,
            mode: 'lines'
        };

        let layout = {
            title: 'X(t)',
            showlegend: true
        };

        Plotly.newPlot(plotter, data, layout, {displayModeBar: true});
    }
};
socket.onerror = function (error) {
  statusLabel.innerText = `[error] ${error.message}`;
};
