// Loads Data
function getData(url) {
    var xmlHttp = new XMLHttpRequest();
    xmlHttp.onreadystatechange = function () {
        if (xmlHttp.readyState == 4 && xmlHttp.status == 200)
            csvToJSON(xmlHttp.responseText);
    }
    xmlHttp.open("GET", url, true);
    xmlHttp.send(null);
}


// Converts data to JSON object
function csvToJSON(data) {

    // Splits data by row and extracts headers
    var elements = data.split("\n");
    var headers = elements[0].split(",");
    var answer = [];

    // Loops through each element of the line and makes it an object, then adds it to answer[]
    for (let i = 1; i < elements.length; i++) {
        var thisObj = {};
        var thisLine = elements[i].split(",");
        for (let j = 0; j < headers.length; j++) {
            thisObj[headers[j]] = thisLine[j];
        }
        answer.push(thisObj);
    }

    popularity(answer);
}

// Calculates which stations are most popular. Requires data = JSON
function popularity(data) {
    var d = data;
    answer = {};
    data.forEach(element => {
        start = element["Starting Station ID"];
        end = element["Ending Station ID"];
        if (start) {
            if (answer[start]) {
                answer[start]++;
            } else {
                answer[start] = 1;
            }
        }
        if (end) {
            if (answer[end]) {
                answer[end]++;
            } else {
                answer[end] = 1;
            }
        }
    });

    // Sorts the data
    sortedData = [];
    var sum = 0;
    Object.keys(answer).forEach(key => {
        var v = {
            id: key,
            freq: answer[key],
        }
        sortedData.push(v);
        sum += v.freq;
    });
    sortedData.sort(function (a, b) {
        return b.freq - a.freq;
    });

    var table = document.getElementById("popularity");

    for (let index = 0; index < sortedData.length; index++) {
        const element = sortedData[index];

        var rank = document.createElement("td");
        rank.innerText = index + 1;
        var id = document.createElement("td");
        id.innerText = element.id;
        var freq = document.createElement("td");
        freq.innerText = element.freq;

        var row = document.createElement("tr");
        row.appendChild(rank);
        row.appendChild(id);
        row.appendChild(freq);
        table.appendChild(row);
    }

    sortedData.forEach(element => {
        var id = document.createElement("td");
        id.innerText = element.id;
        var freq = document.createElement("td");
        freq.innerText = element.freq;

        var row = document.createElement("tr");
        row.appendChild(id);
        row.appendChild(freq);
        table.appendChild(row);
    });
    myData.innerText = sum;
    // myData.innerText = JSON.stringify(sortedData);
}
