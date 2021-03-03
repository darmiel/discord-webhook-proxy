const api = "/cms/history";

function htmlDecode(input) {
    const elem = document.createElement('div');
    elem.innerHTML = input;
    return elem.childNodes.length === 0 ? "" : elem.childNodes[0].nodeValue;
}

function showDiffWithCurrent(url, index) {

    const b64 = btoa(url);
    console.log("b64:", b64);

    const req = api + "/" + b64 + "/" + index;
    console.log("req:", req);

    Swal.fire({
        position: 'center',
        title: "Loading Diff...",
        showConfirmButton: false,
        didOpen: async () => {
            Swal.showLoading();

            // make request
            console.log("fetching", req, "...");
            let res = await fetch(req, {
                method: "GET"
            });
            console.log("fetched:");
            console.log(res);

            // get response
            let text = await res.text();
            if (res.status === 401 && text === "") {
                text = "Unauthorized, you are not allowed to view the history of a cms page";
            }
            console.log("text:", text);

            if (res.status !== 200) {
                return Swal.insertQueueStep({
                    position: 'top-end',
                    icon: "error",
                    title: "Error:",
                    text: text
                });
            }

            console.log("nice :)");

            Swal.fire({
                position: 'center',
                title: "Diff #" + index + " for " + url,
                html: `<div style="text-align: left !important; color: #99AAB5 !important;">${text}</div>`
            });
        }
    });
}