// Load CSV File as a String
function runCSV(flag) {
    fetch('data/metro-bike-share-trip-data.csv').then(function (response) {
        if (response.status !== 200) {
            throw response.status;
        }
        console.log(response.text());
        return response.text();
    }).catch(function (status) {
        console.log('Error ' + status);
    });
}