document.getElementById('searchBtn').addEventListener('click', function () {
  const track = document.getElementById('track').value;

  // const proxyUrl = 'http://192.168.0.28:8081/proxy';
  // using localhost to avoid ip changes
  const proxyUrl = 'http://localhost:8081/proxy';
  const apiUrl = `https://api.deezer.com/search?q=${track}`;
  const url = `${proxyUrl}?url=${encodeURIComponent(apiUrl)}`;

  fetch(url)
    .then(response => response.json())
    .then(data => { displayResults(data); })
    .catch(error => console.error('Error: ', error));
});


function displayResults(data) {
  const resultsDiv = document.getElementById('results');
  resultsDiv.innerHTML = '';

  data.data.forEach(track => {
    const trackInfo = document.createElement('div');
    const downloadBtnId = "downloadBtn" + track.id;
    trackInfo.innerHTML = `
      <h3>${track.title}</h3>
      <p>Link: ${track.link}</p>
      <p>Artist: ${track.artist.name}</p>
      <p>Album: ${track.album.title}</p>
      <img src="${track.album.cover_small}" alt="${track.album.title}">
      <audio controls>
          <source src="${track.preview}" type="audio/mpeg">
          Your browser does not support the audio element.
      </audio>
      <button id="${downloadBtnId}">Download</button>
      <hr>
    `;
    resultsDiv.appendChild(trackInfo);
    loadDownloadFuncionality(track);
  });
}

function loadDetailsFuncionality(track) {
  const detailsButton = "detailsBtn" + track.id;

}

function loadDownloadFuncionality(track) {
  const downloadBtnId = "downloadBtn" + track.id;
  let userSelected = document.getElementById('userSelected');
  console.log(userSelected.value);
  document.getElementById(downloadBtnId).addEventListener('click', function () {
    // fetch('http://192.168.0.28:8081/download', {
    fetch('http://localhost:8081/download', {
      method: "POST",
      headers: {
        'Content-Type': 'application/json'
      },
      body: JSON.stringify({
        url: track.link,
        title: track.title,
        artist: track.artist.name,
        album: track.album.title,
        user: userSelected.value,
      })
    })
      .then(response => response.text())
      .then(downloadId => {
        // Redirect to the status page with the download ID
        window.location.href = `/download-status.html?downloadId=${encodeURIComponent(downloadId)}`;
      })
      .catch(error => {
        console.error("Error calling download in the backend", error);
        alert("Error al iniciar la descarga");
      });
    console.log("Sending data:", {
      url: track.link,
      title: track.title,
      artist: track.artist.name,
      album: track.album.title,
      user: userSelected.value,
    });
  });
}

