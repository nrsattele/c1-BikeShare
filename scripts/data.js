// Loads Data
function getData(flag) {
    fetch('data/metro-bike-share-trip-data.csv').then(function (response) {
        if (response.status !== 200) {
            throw response.status;
        }
        console.log(response.text());
        var myData = document.getElementById('myData');
        myData.innerText = response.text();
        return response.text();
    }).then(function(data) {
        console.log(data);
    }).catch(function (status) {
        console.log('Error ' + status);
    });
}

function httpGetAsync(theUrl, callback){
    var xmlHttp = new XMLHttpRequest();
    xmlHttp.onreadystatechange = function() { 
        if (xmlHttp.readyState == 4 && xmlHttp.status == 200)
            callback(xmlHttp.responseText);
    }
    xmlHttp.open("GET", theUrl, true); // true for asynchronous 
    xmlHttp.send(null);
}

function callback(respText) {
    console.log(respText);
}