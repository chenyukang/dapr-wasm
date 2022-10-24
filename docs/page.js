var image;

function fileSelected(e) {
    const file = e.files[0];
    if (!file) {
        return;
    }

    if (!file.type.startsWith('image/')) {
        alert('Please select a image.');
        return;
    }

    const img = document.createElement('img-tag');
    img.file = file
    image = img;

    const reader = new FileReader();
    reader.onload = function(e) {
        var elem = document.getElementById("upload-pic");
        elem.src = e.target.result;
        elem.hidden = false;
        var origin_img = document.getElementById("origin-pic");
        origin_img.src = elem.src;
        var button = document.getElementById("run");
        button.removeAttribute("disabled");
        setRes("");
    }
    reader.readAsDataURL(file);
}

function setButton() {
    var button = document.getElementById("run");
    button.innerText = "Submit";
    button.disabled = false;
}

function setLoading(loading) {
    var button = document.getElementById("run");
    if (loading) {
        button.disabled = true;
        button.innerText = "Sending ...";
    } else {
        setButton();
    }
}

function setRes(res) {
    var elem = document.getElementById("result");
    elem.innerHTML = res;
    elem.hidden = false;
    var row = document.getElementById("grayscale-rows");
    row.hidden = true;
    var elem = document.getElementById("infer-rows");
    elem.hidden = false;
}

function setImageRes(data) {
    if (data == "ImageTooLarge") {
        alert("Image Too Large");
    } else {
        var row = document.getElementById("grayscale-rows");
        row.hidden = false;
        var elem = document.getElementById("infer-rows");
        elem.hidden = true;
        var img = document.getElementById("processed-pic");
        img.src = "data:image/png;base64, " + data;
        img.hidden = false;
        var origin_img = document.getElementById("origin-pic");
        origin_img.src = document.getElementById("upload-pic").src;
        origin_img.hidden = false;
    }
}

function getApi() {
    var select = document.getElementById('run-api');
    return select.options[select.selectedIndex].value;
}

function updateStat(api) {
    var xhr = new XMLHttpRequest();
    xhr.open('GET', 'http://23.100.38.125:9000/api/invokecount?api=' + api, true);
    xhr.responseType = 'json';
    xhr.onload = function() {
        var status = xhr.status;
        if (status === 200) {
            var res = xhr.response;
            document.getElementById('stat').hidden = false;
            console.log(res);
            let html = "";
            for (const elem of res) {
                items = elem.split("##");
                key = items[0];
                value = items[1];
                html += "<tr><td>" + key + "&nbsp</td><td>" + value + "</td></tr><br>";
            }
            document.getElementById('stat-text').innerHTML = html;
        }
    };
    xhr.send();
}

function runWasm(e) {
    const reader = new FileReader();
    reader.onload = function(e) {
        setLoading(true);
        var req = new XMLHttpRequest();
        var api = getApi();
        var uri = `http://23.100.38.125:3504/v1.0/invoke/image-api-${api}/method/${api}`;
        req.open("POST", uri, true);
        req.setRequestHeader('api', getApi());
        req.onload = function() {
            setLoading(false);
            if (req.status == 200) {
                var header = req.getResponseHeader("Content-Type");
                console.log(header);
                if (header == "image/png") {
                    setImageRes(req.response);
                } else {
                    setRes(req.response);
                }
                //updateStat(getApi());
            } else {
                setRes("API error with status: " + req.status);
            }
        };
        const blob = new Blob([e.target.result], {
            type: 'application/octet-stream'
        });
        req.send(blob);
    };
    console.log(image.file)
    reader.readAsArrayBuffer(image.file);
}