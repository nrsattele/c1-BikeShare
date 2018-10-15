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
    }).catch(function (status) {
        console.log('Error ' + status);
    });
}